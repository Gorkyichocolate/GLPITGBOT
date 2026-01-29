package models

type State int

const (
	StateIdle State = iota
	StateWaitTitle
	StateWaitDescription
)
