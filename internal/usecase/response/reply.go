package response

type Kind int

const (
	KindSend Kind = iota
	KindEdit
)

type Text struct {
	Key        string
	Args       map[string]any
	Fallback   string
	SubKeyArgs []string // keys in Args whose values are i18n keys and need translating
}

type Reply struct {
	Kind Kind

	Text Text

	Keyboard Keyboard
}

type Keyboard struct {
	Remove bool

	ReplyRows  [][]Text
	InlineRows [][]Button
}

type Button struct {
	Text Text
	Data CallbackData
}

type CallbackData struct {
	Unique string
	Args   []string
}

func (r *Reply) IsEmpty() bool {
	if r.Text.Key != "" || r.Text.Fallback != "" {
		return false
	}
	if r.Keyboard.Remove {
		return false
	}
	if len(r.Keyboard.ReplyRows) > 0 {
		return false
	}
	if len(r.Keyboard.InlineRows) > 0 {
		return false
	}
	return true
}
