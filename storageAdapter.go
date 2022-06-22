package pluginPersistency

type StorageAdapter interface {
	Initialize(path string) error
	AppendLogs(key string, logs string)
}
