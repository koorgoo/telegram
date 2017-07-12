package telegram

//go:generate python keyboards.py

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"
)

type Update struct {
	UpdateID          int      `json:"update_id"`
	Message           *Message `json:"message"`
	EditedMessage     *Message `json:"edited_message"`
	ChannelPost       *Message `json:"channel_post"`
	EditedChannelPost *Message `json:"edited_channel_post"`
	// InlineQuery
	// ChosenInlineResult
	CallbackQuery *CallbackQuery `json:"callback_query"`
	// ShippingQuery
	// PreCheckoutQuery
}

type WebhookInfo struct {
	URL                  string   `json:"url"`
	HasCustomCertificate bool     `json:"has_custom_certificate"`
	PendingUpdateCount   int      `json:"pending_update_count"`
	LastErrorDate        int      `json:"last_error_date"`
	LastErrorMessage     string   `json:"last_error_message"`
	MaxConnections       int      `json:"max_connections"`
	AllowedUpdates       []string `json:"allowed_updates"`
}

type User struct {
	ID           int     `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     *string `json:"last_name"`
	Username     *string `json:"username"`
	LanguageCode *string `json:"language_code"`
}

type Chat struct {
	ID                          int64      `json:"id"`
	Type                        string     `json:"type"`
	Title                       *string    `json:"title"`
	Username                    *string    `json:"username"`
	FirstName                   *string    `json:"first_name"`
	LastName                    *string    `json:"last_name"`
	AllMembersAreAdministrators *bool      `json:"all_members_are_administrators"`
	ChatPhoto                   *ChatPhoto `json:"chat_photo"`
	Description                 *string    `json:"description"`
	InviteLink                  *string    `json:"invite_link"`
}

func (c *Chat) IsPrivate() bool    { return c.Type == "private" }
func (c *Chat) IsGroup() bool      { return c.Type == "group" }
func (c *Chat) IsSupergroup() bool { return c.Type == "supergroup" }
func (c *Chat) IsChannel() bool    { return c.Type == "channel" }

type Message struct {
	MessageID       int              `json:"message_id"`
	From            *User            `json:"from"`
	Date            int              `json:"date"`
	Chat            Chat             `json:"chat"`
	ForwardFrom     *User            `json:"forward_from"`
	ForwardFromChat *Chat            `json:"forward_from_chat"`
	ForwardDate     *int             `json:"forward_date"`
	ReplyToMessage  *Message         `json:"reply_to_message"`
	EditDate        *int             `json:"edit_date"`
	Text            *string          `json:"text"`
	Entities        []*MessageEntity `json:"entities"`
	// Audio
	// Document
	// Game
	// Photo
	// Sticker
	// Video
	// Voice
	// VideoNote
	NewChatMembers []*User `json:"new_chat_members"`
	Caption        *string `json:"caption"`
	// Contact
	// Location
	// Venue
	NewChatMember  *User   `json:"new_chat_member"`
	LeftChatMember *User   `json:"left_chat_member"`
	NewChatTitle   *string `json:"new_chat_title"`
	// NewChatPhoto
	DeleteChatPhoto       *bool    `json:"delete_chat_photo"`
	GroupChatCreated      *bool    `json:"group_chat_created"`
	SupergroupChatCreated *bool    `json:"supergroup_chat_created"`
	ChannelChatCreated    *bool    `json:"channel_chat_created"`
	MigrateToChatID       *int64   `json:"migrate_to_chat_id"`
	MigrateFromChatID     *int64   `json:"migrate_from_chat_id"`
	PinnedMessage         *Message `json:"pinned_message"`
	// Invoice
	// SuccessfulPayment
}

type MessageEntity struct {
	Type   string  `json:"type"`
	Offset int     `json:"offset"`
	Length int     `json:"length"`
	URL    *string `json:"url"`
	User   *User   `json:"user"`
}

func (e *MessageEntity) IsMention() bool    { return e.Type == "mention" }
func (e *MessageEntity) IsHashtag() bool    { return e.Type == "hastag" }
func (e *MessageEntity) IsBotCommand() bool { return e.Type == "bot_command" }
func (e *MessageEntity) IsURL() bool        { return e.Type == "url" }
func (e *MessageEntity) IsEmail() bool      { return e.Type == "email" }

// PhotoSize
// Audio
// Document
// Sticker
// Video
// Voice
// VideoNote
// Contact
// Location
// Venue
// UserProfilePhotos
// File

type ReplyKeyboardMarkup struct {
	Keyboard        [][]*KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool                `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool                `json:"one_time_keyboard,omitempty"`
	Selective       bool                `json:"selective,omitempty"`
}

type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard,omitempty"`
	Selective      bool `json:"selective,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]*InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	URL                          string `json:"url,omitempty"`
	CallbackData                 string `json:"callback_data,omitempty"`
	SwitchInlineQuery            string `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"`
	// CallbackGame
	// Pay
}

// TODO: Ensure Data may not be sent to leave a pointer to string.
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            User     `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageID *string  `json:"inline_message_id"`
	ChatInstance    string   `json:"chat_instance"`
	Data            *string  `json:"data"`
	// GameShortName *string
}

type ForceReply struct {
	ForeceReply bool `json:"force_reply"`
	Selective   bool `json:"selective"`
}

type ChatPhoto struct {
	SmallFileID string `json:"small_file_id"`
	BigFileID   string `json:"big_file_id"`
}

type ChatMember struct {
	User      User   `json:"user"`
	Status    string `json:"status"`
	UntilDate *int   `json:"until_date"`
	// administrators only stuff
	CanBeEdited           *bool `json:"can_be_edited"`
	CanChangeInfo         *bool `json:"can_change_info"`
	CanPostMessages       *bool `json:"can_post_messages"`
	CanEditMessages       *bool `json:"can_edit_messages"`
	CanDeleteMessages     *bool `json:"can_delete_messages"`
	CanInviteUsers        *bool `json:"can_invite_users"`
	CanRestrictMembers    *bool `json:"can_restrict_members"`
	CanPinMessages        *bool `json:"can_pin_messages"`
	CanPromoteMembers     *bool `json:"can_promote_members"`
	CanSendMessages       *bool `json:"can_send_messages"`
	CanSendMediaMessages  *bool `json:"can_send_media_messages"`
	CanSendOtherMessages  *bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews *bool `json:"can_add_web_page_previews"`
}

type ResponseParameters struct {
	MigrateToChatID *int64 `json:"migrate_to_chat_id"`
	RetryAfter      *int   `json:"retry_after"`
}

type InputFile interface {
	io.Reader
	Name() string
}

// Parse modes.
const (
	ModeDefault  ParseMode = 0
	ModeMarkdown           = 1
	ModeHTML               = 2
)

type ParseMode int

// MarshalJSON implements json.Marshaler interface.
func (m ParseMode) MarshalJSON() (b []byte, err error) {
	switch m {
	case ModeDefault:
		// ModeDefault should be evaluated as empty by json package and skipped
		// in message marshalling. But empty string was a simplest solution to
		// leave ParseMode a simple type.
		b = []byte(`""`)
	case ModeMarkdown:
		b = []byte(`"Markdown"`)
	case ModeHTML:
		b = []byte(`"HTML"`)
	}
	return
}

type NewMessage struct {
	ChatID                int64     `json:"chat_id"`
	Text                  string    `json:"text"`
	ParseMode             ParseMode `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool      `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool      `json:"disable_notification,omitempty"`
	ReplyToMessageID      int       `json:"reply_to_message_id,omitempty"`
	ReplyMarkup           Markup    `json:"reply_markup,omitempty"`
}

type Markup interface {
	json.Marshaler
	json.Unmarshaler
}

var _ = Markup((*ReplyKeyboardMarkup)(nil))
var _ = Markup((*ReplyKeyboardRemove)(nil))
var _ = Markup((*InlineKeyboardMarkup)(nil))
var _ = Markup((*ForceReply)(nil))

type MessageText struct {
	ChatID                int64                 `json:"chat_id,omitempty"`
	MessageID             int                   `json:"message_id,omitempty"`
	InlineMessageID       int                   `json:"inline_message_id,omitempty"`
	Text                  string                `json:"text"`
	ParseMode             ParseMode             `json:"parse_mode"`
	DisableWebPagePreview bool                  `json:"disable_web_page_preview"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type ForwardMessage struct {
	ChatID              int64 `json:"chat_id"`
	FromChatID          int64 `json:"from_chat_id"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
	MessageID           int   `json:"message_id"`
}

type PhotoMessage struct {
	ChatID              int64     `json:"chat_id"`
	Photo               InputFile `json:"-"`
	PhotoID             string    `json:"photo,omitempty"`
	Caption             string    `json:"caption,omitempty"`
	DisableNotification bool      `json:"disable_notification,omitempty"`
	ReplyToMessageID    int       `json:"reply_to_message_id,omitempty"`
}

var _ = (Multiparter)((*PhotoMessage)(nil))

// Multipart implements Multiparter interface.
func (m *PhotoMessage) Multipart() *Multipart {
	if m.Photo == nil {
		return nil
	}
	return &Multipart{
		Files: map[string]InputFile{"photo": m.Photo},
		Form: url.Values{
			"chat_id":              {strconv.FormatInt(m.ChatID, 10)},
			"caption":              {m.Caption},
			"disable_notification": {strconv.FormatBool(m.DisableNotification)},
			"reply_to_message_id":  {strconv.FormatInt(int64(m.ReplyToMessageID), 10)},
		},
	}
}

// AudioMessage
// DocumentMessage
// StickerMessage
// VideoMessage
// VoiceMessage
// VideoNoteMessage
// LocationMessage
// VenueMessage
// ContactMessage
// ChatActionMessage
// UserProfilePhotosMessage
// KickChatMemberMessage
// UnbanChatMemberMessage
// RestrictChatMemberMessage
// PromoteChatMemberMessage
// ExportChatInviteLinkMessage
// SetChatPhotoMessage
// DeleteChatPhotoMessage
// SetChatTitleMessage
// SetChatDescriptionMessage
// PinChatMessageMessage
// UnpinChatMessageMessage
// LeaveChatMessage

// getChat
// getChatAdministrators
// getChatMembersCount
// getChatMember

type CallbackQueryAnswer struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
	URL             string `json:"url,omitempty"`
	CacheTime       int    `json:"cache_time,omitempty"`
}

type MessageCaption struct {
	ChatID          int64                 `json:"chat_id,omitempty"`
	MessageID       int                   `json:"message_id,omitempty"`
	InlineMessageID int                   `json:"inline_message_id,omitempty"`
	Caption         string                `json:"caption,omitempty"`
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type MessageReplyMarkup struct {
	ChatID          int64                 `json:"chat_id,omitempty"`
	MessageID       int                   `json:"message_id,omitempty"`
	InlineMessageID int                   `json:"inline_message_id,omitempty"`
	ReplyMarkup     *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type MessageDeletion struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int   `json:"message_id"`
}
