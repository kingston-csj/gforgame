package reward

type Reward interface {
	GetAmount() int

	AddAmount(amount int)

	GetType() string

	Serial() string
}
