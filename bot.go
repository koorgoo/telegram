package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

const (
	jsonContentType = "application/json;chartset=utf-8"
)

const (
	defaultURL         = "https://api.telegram.org/bot"
	defaultErrTimeout  = 5 * time.Second
	defaultPollTimeout = time.Minute
)

var (
	ErrNotDeleted  = errors.New("telegram: message not deleted")
	ErrNotEdited   = errors.New("telegram: message not edited")
	ErrNotAnswered = errors.New("telegram: query not answered")
)

type Bot struct {
	name        string
	url         string
	ctx         context.Context
	errTimeout  time.Duration
	pollTimeout time.Duration
	noUpdates   bool
	updatec     chan []*Update
	errorc      chan error
}

func NewBot(ctx context.Context, token string, opts ...BotOption) (*Bot, error) {
	bot := newBot(ctx, token, opts...)
	// ensure bot works
	_, err := bot.GetMe(context.TODO())
	if err != nil {
		return nil, err
	}
	if !bot.noUpdates {
		go bot.listenToUpdates()
	}
	return bot, nil
}

type botOptions struct {
	Name        string
	URL         string
	ErrTimeout  time.Duration
	PollTimeout time.Duration
	NoUpdates   bool
}

type BotOption func(*botOptions)

func withURL(url string) BotOption {
	return func(o *botOptions) {
		o.URL = url
	}
}

func WithErrTimeout(t time.Duration) BotOption {
	return func(o *botOptions) {
		o.ErrTimeout = t
	}
}

func WithPollTimeout(t time.Duration) BotOption {
	return func(o *botOptions) {
		o.PollTimeout = t
	}
}

func WithoutUpdates() BotOption {
	return func(o *botOptions) {
		o.NoUpdates = true
	}
}

func newBot(ctx context.Context, token string, opts ...BotOption) *Bot {
	o := &botOptions{URL: defaultURL, ErrTimeout: defaultErrTimeout, PollTimeout: defaultPollTimeout}
	for _, opt := range opts {
		opt(o)
	}
	bot := &Bot{
		name:        o.Name,
		url:         o.URL + token,
		ctx:         ctx,
		errTimeout:  o.ErrTimeout,
		pollTimeout: o.PollTimeout,
		noUpdates:   o.NoUpdates,
		updatec:     make(chan []*Update),
		errorc:      make(chan error),
	}
	if bot.noUpdates {
		close(bot.updatec)
		close(bot.errorc)
	}
	return bot
}

func (b *Bot) listenToUpdates() {
	var offset int
	donec := b.ctx.Done()
loop:
	for {
		u, err := b.GetUpdates(b.ctx, WithOffset(offset), WithTimeout(b.pollTimeout))
		// Handle context errors differently - shutdown gracefully.
		switch err {
		case context.Canceled, context.DeadlineExceeded:
			break loop
		}

		if err != nil {
			select {
			case b.errorc <- err:
				sleepctx(b.ctx, b.errTimeout)
				continue
			case <-donec:
				break
			}
		}
		// No updates this time - repeat the loop and wait for another pack.
		if len(u) == 0 {
			continue
		}
		// Increment offset according to the last update id. Next time updates
		// pack will not contain updates up to this last one.
		offset = u[len(u)-1].UpdateID + 1

		select {
		case b.updatec <- u:
			continue
		case <-donec:
			break
		}
	}

	// TODO: How to ensure updatesc and errorc to be drained?

	// Don't forget to close channels.
	close(b.updatec)
	close(b.errorc)
}

// sleepctx pauses for at lease t duration. It returns early if ctx is cancelled or
// its deadline is exceeded.
func sleepctx(ctx context.Context, t time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(t):
	}
}

func (b *Bot) Updates() <-chan []*Update { return b.updatec }
func (b *Bot) Errors() <-chan error      { return b.errorc }

// call issues HTTP request to API for the method with form values and decodes
// received data in v. It returns error otherwise.
func (b *Bot) do(ctx context.Context, method string, data interface{}, v interface{}) error {
	client := http.DefaultClient
	url := b.url + "/" + method

	body, contentType, err := b.encode(data)
	if err != nil {
		return err
	}
	resp, err := ctxhttp.Post(ctx, client, url, contentType, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r := new(apiResponse)
	if err := json.Unmarshal(bdata, r); err != nil {
		return err
	}

	if !r.OK {
		return &APIError{
			ErrorCode:   r.ErrorCode,
			Description: r.Description,
			Parameters:  r.Parameters,
		}
	}

	return json.Unmarshal([]byte(r.Result), v)
}

func (b *Bot) encode(data interface{}) (io.Reader, string, error) {
	if m, ok := data.(Multiparter); ok {
		if v := m.Multipart(); v != nil {
			return v.Encode()
		}
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return nil, "", err
	}
	return buf, jsonContentType, nil
}

// Multiparter is an interface for messages that may be converted to a multipart
// form (e.g. photo, document, video). *Multipart may be nil meaning unavailable
// conversion.
type Multiparter interface {
	Multipart() *Multipart
}

type Multipart struct {
	Form  url.Values
	Files map[string]InputFile
}

// Encode encodes Multipart to multipart/form-data. It returns io.Reader for
// content, content type with boundary and error. In case of failed encoding
// io.Reader is nil, content type is an empty string.
func (m *Multipart) Encode() (io.Reader, string, error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	for key := range m.Form {
		if err := w.WriteField(key, m.Form.Get(key)); err != nil {
			return nil, "", err
		}
	}

	for key := range m.Files {
		file := m.Files[key]
		if dest, err := w.CreateFormFile(key, file.Name()); err != nil {
			return nil, "", err
		} else {
			if _, err := io.Copy(dest, file); err != nil {
				return nil, "", err
			}
		}
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return buf, w.FormDataContentType(), nil
}

type updatesOptions struct {
	Offset  int
	Limit   int
	Timeout time.Duration
	// AllowedUpdates []string
}

// MarshalJSON implements json.Marshaler interface.
func (o *updatesOptions) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	if o.Offset > 0 {
		m["offset"] = o.Offset
	}
	if o.Limit > 0 {
		m["limit"] = o.Limit
	}
	if o.Timeout > 0 {
		m["timeout"] = int(o.Timeout.Seconds())
	}
	return json.Marshal(m)
}

type UpdatesOption func(*updatesOptions)

// WithOffset sets id of the first expected update in response. Usually offset
// should equal last update's id + 1.
func WithOffset(offset int) UpdatesOption {
	return func(o *updatesOptions) {
		o.Offset = offset
	}
}

// WithLimit modifies updates request to limit the number of updates in response.
func WithLimit(limit int) UpdatesOption {
	return func(o *updatesOptions) {
		if limit < 1 {
			limit = 1
		}
		if limit > 100 {
			limit = 100
		}
		o.Limit = limit
	}
}

// WithTimeout modifies timeout of updates request. 0 duration means short
// polling (for testing only).
func WithTimeout(t time.Duration) UpdatesOption {
	return func(o *updatesOptions) {
		o.Timeout = t
	}
}

// APIError represents an error returned by API. It satisfies error interface.
type APIError struct {
	ErrorCode   int
	Description string
	Parameters  *ResponseParameters
}

// Error returns an error string.
func (e *APIError) Error() string {
	return fmt.Sprintf("telegram: %d %s", e.ErrorCode, e.Description)
}

// apiResponse represents API response. When OK is false then ErrorCode and
// Description defines the error situation.
type apiResponse struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	// error part
	ErrorCode   int                 `json:"error_code,omitempty"`
	Description string              `json:"description,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}
