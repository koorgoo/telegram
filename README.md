### telegram

`telegram` is a simple Telegram Bot API client.


#### Principles

> Good design is honest.
> Dieter Rams

Types must be clear about its fields' type and optionality. Because of this some
ids are int64. And some fields are pointers to bool, int, or string.

Types in the package explicitly mirror types of
[Telegram Bot API](https://core.telegram.org/bots/api#available-types).


#### Installation

```
go get -u github.com/koorgoo/telegram
```


#### Examples

- [How to send a message](examples/sendmessage.go)
- [How to receive updates using channel](examples/updates.go)
- [How to handle commands from updates](examples/commands.go)
