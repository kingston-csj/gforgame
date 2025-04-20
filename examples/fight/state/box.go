package state

type StateBox struct {
	states map[StateType]bool
}

func NewStateBox() *StateBox {
	return &StateBox{
		states: make(map[StateType]bool),
	}
}

func (s *StateBox) AddState(state StateType) {
	s.states[state] = true
}

func (s *StateBox) RemoveState(state StateType) {
	s.states[state] = false
}

func (s *StateBox) HasState(state StateType) bool {
	return s.states[state]
}

func (s *StateBox) CanAttack() bool {
	return !s.HasState(StateType_Stun) && !s.HasState(StateType_Sleep)
}
