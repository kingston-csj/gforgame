package camp

// 根据阵营获取主公武将ID
func GetHeroIdByCamp(camp int32) int32 {
	switch camp {
	case Camp_Shu:
		return 1001
	case Camp_Wei:
		return 1002
	case Camp_Wu:
		return 1003
	case Camp_Hao:
		return 1004
	}
	return 0
}
