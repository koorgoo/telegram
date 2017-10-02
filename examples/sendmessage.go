package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	tg "github.com/koorgoo/telegram"
)

func main() {
	ctx := context.Background()
	b, err := tg.NewBot(ctx, os.Getenv("TOKEN"), tg.WithoutUpdates())
	if err != nil {
		panic(err)
	}

	id, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 32)
	m, err := b.SendMessage(ctx, &tg.TextMessage{
		ChatID:              int64(id),
		Text:                "At your service.",
		DisableNotification: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("sent: %+v\n", m)
}
