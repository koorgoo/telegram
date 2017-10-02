package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tg "github.com/koorgoo/telegram"
)

func main() {
	bot, err := tg.NewBot(context.Background(), os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	cmd := tg.NewCommands(bot.Username())
	cmd.Add("/hello", Hello(bot))

	callCommand := func(u *tg.Update) {
		if err, _ := cmd.Run(u); err != nil {
			log.Println("error:", err)
		}
	}

	for update := range bot.Updates() {
		go callCommand(update)
	}
}

func Hello(bot tg.Bot) tg.CommandFunc {
	return func(c *tg.Command, u *tg.Update) error {
		name := "user"
		if len(c.Args) > 0 {
			name = c.Args[0]
		}
		_, err := bot.SendMessage(context.Background(), &tg.TextMessage{
			ChatID: u.Message.Chat.ID,
			Text:   fmt.Sprintf("Hello, %s!", name),
		})
		return err
	}
}
