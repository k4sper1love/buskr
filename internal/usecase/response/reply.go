package response

type Kind int

const (
	KindSend Kind = iota
	KindEdit
	KindSendImage
)

type Text struct {
	Key          string
	Args         map[string]any
	Fallback     string
	SubKeyArgs   []string // keys in Args whose values are i18n keys and need translating
	NoEscapeArgs []string // keys in Args that should NOT be escaped
}

type Reply struct {
	Kind     Kind
	Text     Text // or caption for image
	Keyboard Keyboard
	Image    []byte // nil if no image
}

type Keyboard struct {
	Remove bool

	ReplyRows  [][]Text
	InlineRows [][]Button
}

type Button struct {
	Text      Text
	Data      CallbackData
	URL       string
	WebAppURL string
}

type CallbackData struct {
	Unique string
	Args   []string
}

func (r *Reply) IsEmpty() bool {
	if r.Text.Key != "" || r.Text.Fallback != "" {
		return false
	}
	if len(r.Image) > 0 {
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
