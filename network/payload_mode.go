package network

// PayloadMode controls whether a session decodes message body into typed structs
// or forwards raw payload bytes with only protocol headers parsed.
type PayloadMode int

const (
	// PayloadModeDecode decodes body by cmd -> message type mapping.
	PayloadModeDecode PayloadMode = iota
	// PayloadModeRawBody keeps body as []byte and only parses protocol headers.
	PayloadModeRawBody
)
