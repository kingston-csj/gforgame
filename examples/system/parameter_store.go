package system

import (
	"io/github/gforgame/examples/context"
)

func saveSystemParameterValue(id string, value string, payload any) {
	cache, _ := context.CacheManager.GetCache(systemParameterCacheTable)
	cache.Set(id, payload)
	record := GetSystemParameterService().GetOrCreateSystemParameterRecord(id)
	record.Data = value
	context.DbService.SaveToDb(record)
}

func loadSystemParameterValue(id string) string {
	record := GetSystemParameterService().GetOrCreateSystemParameterRecord(id)
	if record == nil {
		return ""
	}
	return record.GetData()
}
