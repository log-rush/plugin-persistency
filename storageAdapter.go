package pluginPersistency

type StorageAdapter interface {
	AppendLogs(key string, logs string)
}
