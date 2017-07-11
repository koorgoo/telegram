package telegram

import (
	"encoding/json"
	"testing"
)

func TestReplyKeyboardMarkup_MarshalJSON(t *testing.T) {
	m := &ReplyKeyboardMarkup{
		Keyboard: [][]*KeyboardButton{{{Text: "test"}}},
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	want := `"{\"keyboard\":[[{\"text\":\"test\"}]]}"`
	if s := string(b); s != want {
		t.Fatalf("json: want %q, got %q", want, s)
	}
}

func TestReplyKeyboardMarkup_UnmarshalJSON(t *testing.T) {
	var m ReplyKeyboardMarkup
	s := `"{\"keyboard\":[[{\"text\":\"test\"}]]}"`
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatal(err)
	}
	if s := m.Keyboard[0][0].Text; s != "test" {
		t.Fatalf("button: want %q, got %q", "test", s)
	}
}
