package telegram

import (
	"context"
	"errors"
)

var (
	ErrNotDeleted  = errors.New("telegram: message not deleted")
	ErrNotEdited   = errors.New("telegram: message not edited")
	ErrNotAnswered = errors.New("telegram: query not answered")
)

// GetMe returns a basic information about the bot.
func (b *bot) GetMe(ctx context.Context) (*User, error) {
	var u *User
	if err := b.do(ctx, "getMe", nil, &u); err != nil {
		return nil, err
	}
	return u, nil
}

// GetUpdates returns a slice of updates received with given options.
func (b *bot) GetUpdates(ctx context.Context, opts ...UpdatesOption) ([]*Update, error) {
	uo := new(updatesOptions)
	for _, opt := range opts {
		opt(uo)
	}
	var v []*Update
	if err := b.do(ctx, "getUpdates", uo, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// SetWebhook
// DeleteWebhook
// GetWebhookInfo
// WebhookInfo

// SendMessages sends m message and returns a sent message instance on success.
func (b *bot) SendMessage(ctx context.Context, m *NewMessage) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "sendMessage", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (b *bot) ForwardMessage(ctx context.Context, m *ForwardMessage) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "forwardMessage", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (b *bot) SendPhoto(ctx context.Context, m *PhotoMessage) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "sendPhoto", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// SendAudio
// SendDocument
// SendSticker
// SendVideo
// SendVoice
// SendVoiceNote
// SendLocation
// SendVenue
// SendContact
// SendChatAction
// GetUserProfilePhotos
// GetFile
// KickChatMember
// UnbanChatMember
// RestrictChatMember
// PromoteChatMember
// ExportChatInviteLink
// SetChatPhoto
// SetChatDescription
// PinChatMessage
// UnpinChatMessage
// LeaveChat
// GetChat
// GetChatAdministrators
// GetChatMembersCount
// GetChatMembers

func (b *bot) AnswerCallbackQuery(ctx context.Context, a *CallbackQueryAnswer) error {
	var ok bool
	if err := b.do(ctx, "answerCallbackQuery", a, &ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotAnswered
	}
	return nil
}

// TODO: What does True mean for edit* methods?
// > On success, if edited message is sent by the bot, the edited Message is
// > returned, otherwise True is returned.
// https://core.telegram.org/bots/api#editmessagetext

func (b *bot) EditMessageText(ctx context.Context, t *MessageText) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "editMessageText", t, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (b *bot) EditMessageCaption(ctx context.Context, c *MessageCaption) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "editMessageCaption", c, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (b *bot) EditMessageReplyMarkup(ctx context.Context, m *MessageReplyMarkup) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "editMessageReplyMarkup", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (b *bot) DeleteMessage(ctx context.Context, d *MessageDeletion) error {
	var ok bool
	if err := b.do(ctx, "deleteMessage", d, &ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotDeleted
	}
	return nil
}
