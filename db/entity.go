package db

type Entity interface {
	GetId() string
	IsDeleted() bool
	SetDeleted()
}

//type BaseEntity struct {
//	Id     string `json:"id"`
//	Delete bool
//}
//
//func (s *BaseEntity) GetId() string {
//	return s.Id
//}
//
//func (s *BaseEntity) IsDeleted() bool {
//	return s.Delete
//}
//
//func (s *BaseEntity) SetDeleted() {
//	s.Delete = true
//}
