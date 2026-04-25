package persist

type Entity interface {
	// 获取实体的唯一标识符
	GetId() string
	// 是否为逻辑已删除状态
	IsDeleted() bool
	// 设置为等删除状态（逻辑已删除，物理上待删除）
	SetDeleted()
	// 在持久化前调用（与 ORM 无关）
	BeforePersist() error
	// 在加载后调用（与 ORM 无关）
	AfterLoad() error
}

type BaseEntity struct {
	Id     string `json:"id"`
	Delete bool   `json:"delete"`
}

func (s *BaseEntity) GetId() string {
	return s.Id
}

func (s *BaseEntity) IsDeleted() bool {
	return s.Delete
}

func (s *BaseEntity) SetDeleted() {
	s.Delete = true
}
