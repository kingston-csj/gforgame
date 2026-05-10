package state

type StateType string

const (
	// 眩晕
	StateType_Stun StateType = "stun"
	// 睡眠
	StateType_Sleep StateType = "sleep"
	// 沉默
	StateType_Silent StateType = "silent"
)
