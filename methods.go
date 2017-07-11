package telegram

import "context"

// GetMe returns a basic information about the bot.
func (b *Bot) GetMe(ctx context.Context) (*User, error) {
	var u *User
	if err := b.do(ctx, "getMe", nil, &u); err != nil {
		return nil, err
	}
	return u, nil
}

// GetUpdates returns a slice of updates received with given options.
func (b *Bot) GetUpdates(ctx context.Context, opts ...UpdatesOption) ([]*Update, error) {
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
func (b *Bot) SendMessage(ctx context.Context, m *NewMessage) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "sendMessage", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (b *Bot) ForwardMessage(ctx context.Context, m *ForwardMessage) (*Message, error) {
	var v *Message
	if err := b.do(ctx, "forwardMessage", m, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (b *Bot) SendPhoto(ctx context.Context, m *PhotoMessage) (*Message, error) {
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

func (b *Bot) AnswerCallbackQuery(ctx context.Context, a *CallbackQueryAnswer) error {
	var ok bool
	if err := b.do(ctx, "answerCallbackQuery", a, &ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotAnswered
	}
	return nil
}

func (b *Bot) EditMessageText(ctx context.Context, t *MessageText) error {
	var ok bool
	if err := b.do(ctx, "editMessageText", t, &ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotEdited
	}
	return nil
}

func (b *Bot) DeleteMessage(ctx context.Context, d *MessageDeletion) error {
	var ok bool
	if err := b.do(ctx, "deleteMessage", d, &ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotDeleted
	}
	return nil
}
