### telegram

`telegram` is a simple Telegram Bot API client.


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

    id, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 32)
    m, err := b.SendMessage(ctx, &tg.NewMessage{
        ChatID: int(id),
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