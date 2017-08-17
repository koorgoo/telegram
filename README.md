### telegram

`telegram` is a simple Telegram Bot API client.


#### Principles

> Good design is honest. 
> Dieter Rams

Types must be clear about its fields' type and optionality. Because of this some
ids are int64. And some fields are pointers to bool, int, or string.

Types in the package explicitly mirror types of
[Telegram Bot API](https://core.telegram.org/bots/api#available-types).


#### Usage

First, install the package.

```
go get -u github.com/koorgoo/telegram
```

Then import it and write your own bot.

```
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
    id, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
    if err != nil {
        panic(err)
    }

    m, err := b.SendMessage(ctx, &tg.TextMessage{
        ChatID: id,
        Text:   "At your service.",
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("sent: %+v\n", m)
}
```

Handle updates as a channel.

```
package main

import (
    "context"
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

    go func(done func(), updatesc <-chan []*tg.Update) {
        defer done()
        for updates := range updatesc {
            // Handle a slice of updates.
            _ = updates
        }
    }(wg.Done, b.Updates())

    go func(done func(), errorc <-chan error) {
        defer done()
        for err := range errorc {
            // Handle an error.
            _ = err
        }
    }(wg.Done, b.Errors())

    wg.Wait()
}
```