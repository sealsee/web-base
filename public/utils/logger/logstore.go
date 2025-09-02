package logger

var LogStore ILogStore

const CAPACITY = 1024

var log_chan = make(chan map[string]any, CAPACITY)
var log_err_chan = make(chan map[string]any, CAPACITY)

type ILogStore interface {
	SaveLog(map[string]any)
	SaveErrLog(map[string]any)
}

func ConfigStore(log ILogStore) {
	LogStore = log
	if LogStore == nil {
		return
	}
	go initLog()
	go initLogErr()
}

func initLog() {
	for {
		if log, ok := <-log_chan; ok {
			LogStore.SaveLog(log)
		}
	}
}

func initLogErr() {
	for {
		if log_err, ok1 := <-log_err_chan; ok1 {
			LogStore.SaveErrLog(log_err)
		}
	}
}

func Log(params map[string]any) {
	if LogStore == nil {
		return
	}
	log_chan <- params

}

func ErrLog(params map[string]any) {
	if LogStore == nil {
		return
	}
	log_err_chan <- params
}
