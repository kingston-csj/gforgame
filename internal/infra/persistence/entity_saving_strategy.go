package persistence

import (
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/internal/constants"
	playerpo "github.com/forfun/gforgame/internal/infra/persistence/po"
	"github.com/forfun/gforgame/persist"
)

type EntitySavingStrategy struct {
}

func (s *EntitySavingStrategy) DoSave(entity persist.Entity) error {
	if entity.IsDeleted() {
		return Db.Delete(entity).Error
	}

	// if err := entity.BeforePersist(); err != nil {
	// 	return err
	// }
	if p, ok := entity.(*playerpo.PlayerPO); ok {
		logger.LogWithActor(p.Id, p.Name, constants.LoggerDebug, "model", "savePlayer2")
	}
	return Db.Save(entity).Error
}
