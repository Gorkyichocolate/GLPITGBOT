package telegram

type State int

const (
	StateIdle State = iota
	StateWaitTitle
	StateWaitDescription
)
