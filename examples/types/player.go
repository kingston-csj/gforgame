package types

// Player represents a basic player type that can be used across different packages
type Player interface {
	GetId() string
	// Add other common player methods that are needed across packages
}
