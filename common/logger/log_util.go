package logger

func LogWithActor(actorID, actorName string, logType string, args ...interface{}) {
	baseArgs := []interface{}{
		"actorId", actorID,
		"actorName", actorName,
	}
	allArgs := append(baseArgs, args...)
	Log(logType, allArgs...)
}