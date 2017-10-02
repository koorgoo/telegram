package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	tg "github.com/koorgoo/telegram"
)

func main() {
	ctx := context.Background()
	b, err := tg.NewBot(ctx, os.Getenv("TOKEN"), tg.WithErrTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func(done func(), updatec <-chan *tg.Update) {
		defer done()
		for update := range updatec {
			log.Printf("update: %+v\n", update)
		}
	}(wg.Done, b.Updates())

	go func(done func(), errorc <-chan error) {
		defer done()
		for err := range errorc {
			log.Println("error:", err)
		}
	}(wg.Done, b.Errors())

	wg.Wait()
}
