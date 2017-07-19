package telegram

import "testing"

var SafeSliceTests = []struct {
	S        string
	From, To int
	Slice    string
}{
	{"01234", 0, 4, "0123"},
	{"01234", -10, 10, "01234"},
	{"01234", 3, 10, "34"},
	{"01234", 5, 10, ""},
	{"01234", 10, 10, ""},
	{"01234", 20, 10, ""},
}

func TestSafeSlice(t *testing.T) {
	for _, tt := range SafeSliceTests {
		if s := safeSlice(tt.S, tt.From, tt.To); s != tt.Slice {
			t.Errorf("want %q, got %q", tt.Slice, s)
		}
	}
}

var SplitCommandTests = []struct {
	S, Command, Mention string
}{
	{"/command", "/command", ""},
	{"/command@bot", "/command", "bot"},
}

func TestSplitCommand(t *testing.T) {
	for _, tt := range SplitCommandTests {
		if c, m := splitCommand(tt.S); c != tt.Command || m != tt.Mention {
			t.Errorf("want (%q, %q), got (%q, %q)", tt.Command, tt.Mention, c, m)
		}
	}
}

var SplitArgsTests = []struct {
	S    string
	Args []string
}{
	{"", nil},
	{"", []string{}},
	{"    ", []string{}},
	{" a  b   ", []string{"a", "b"}},
}

func TestSplitArgs(t *testing.T) {
	for _, tt := range SplitArgsTests {
		if args := splitArgs(tt.S); !stringsEqual(args, tt.Args) {
			t.Errorf("want %s, got %s", tt.Args, args)
		}
	}
}

func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
