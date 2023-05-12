package logger

var LogStore ILogStore

type ILogStore interface {
	SaveLog(map[string]any)
	SaveErrLog(map[string]any)
}

func ConfigStore(log ILogStore) {
	LogStore = log
}

func Log(params map[string]any) {
	if LogStore == nil {
		return
	}
	LogStore.SaveLog(params)
}

func ErrLog(params map[string]any) {
	if LogStore == nil {
		return
	}
	LogStore.SaveErrLog(params)
}
