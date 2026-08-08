package system

import (
	"go.uber.org/dig"
	"gorm.io/gorm"

	mysqldb "github.com/forfun/gforgame/internal/infra/persistence"
	"github.com/forfun/gforgame/persist"
)

// SystemParameterEnt 仅表示系统参数的持久化记录。
type SystemParameterEnt struct {
	persist.BaseEntity
	Data string `gorm:"column:data"`
}

func (s *SystemParameterEnt) GetData() string {
	return s.Data
}

func (s *SystemParameterEnt) SetData(data string) {
	s.Data = data
}

func (s *SystemParameterEnt) GetId() string {
	return s.Id
}

func (s *SystemParameterEnt) GetKey() string {
	return "SystemParameter_" + s.Id
}

func (s *SystemParameterEnt) BeforePersist() error {
	return nil
}

func (s *SystemParameterEnt) AfterLoad() error {
	return nil
}

func (s *SystemParameterEnt) BeforeSave(tx *gorm.DB) error {
	return s.BeforePersist()
}

func (s *SystemParameterEnt) AfterFind(tx *gorm.DB) error {
	return s.AfterLoad()
}

func (s *SystemParameterEnt) SnapshotEntity() (persist.Entity, error) {
	return &SystemParameterEnt{
		BaseEntity: s.BaseEntity,
		Data:       s.Data,
	}, nil
}

// MySQLSystemRepository 纯 db 仓储：只管系统参数的持久化读写，不碰 cache。
type MySQLSystemRepository struct {
	db        *gorm.DB
	dbService *mysqldb.AsyncDBService
}

type MySQLSystemRepositoryParams struct {
	dig.In

	DB        *gorm.DB
	DbService *mysqldb.AsyncDBService
}

func NewMySQLSystemRepository(params MySQLSystemRepositoryParams) *MySQLSystemRepository {
	return &MySQLSystemRepository{
		db:        params.DB,
		dbService: params.DbService,
	}
}

// GetRecord 按主键从 db 查询系统参数记录，未命中返回 nil。
func (r *MySQLSystemRepository) GetRecord(id string) *SystemParameterEnt {
	var record SystemParameterEnt
	result := r.db.First(&record, "id=?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return nil
	}
	return &record
}

// SaveRecord 异步落库系统参数记录。
func (r *MySQLSystemRepository) SaveRecord(record *SystemParameterEnt) {
	r.dbService.SaveToDb(record)
}
