package persistence

import (
	"github.com/forfun/gforgame/persist"
)

type EntitySavingStrategy struct {
}

func (s *EntitySavingStrategy) DoSave(entity persist.Entity) error {
	if entity.IsDeleted() {
		return Db.Delete(entity).Error
	}
	// if _, ok := entity.(*playerpo.PlayerPO); ok {
	// 	logger.LogSystem(constants.LoggerDebug.Code, "model", "savePlayer2")
	// }
	return Db.Save(entity).Error
}
