package consume

type Consume interface {
	GetType() string

	Serial() string
}
