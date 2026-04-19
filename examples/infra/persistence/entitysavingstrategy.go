package persistence

import "io/github/gforgame/persist"

type EntitySavingStrategy struct {
}

func (s *EntitySavingStrategy) DoSave(entity persist.Entity) error {
	if entity.IsDeleted() {
		return Db.Delete(entity).Error
	}

	if err := entity.BeforeSave(nil); err != nil {
		return err
	}

	return Db.Save(entity).Error
}
