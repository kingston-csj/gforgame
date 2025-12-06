package domain

// Player represents a basic player type that can be used across different packages
type Player interface {
	GetId() string
	GetName() string
}


