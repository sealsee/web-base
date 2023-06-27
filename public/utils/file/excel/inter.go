package excel

import (
	"time"

	"github.com/sealsee/web-base/public/ds/page"
)

type RunStatus int

const (
	Next RunStatus = 1 //下一个
	Exit RunStatus = 2 //退出
)

type ImpHandler interface {
	Headers([]string)
	Row(*map[string]string) RunStatus
	After()
}

type ExpHandler interface {
	Title() string
	HeaderColumn() []string // {"表头1,字段名1","表头2,字段名2","表头3,字段名3"...}
	Rows(*page.Page) []map[string]interface{}
	Finish(url string)
}

type ImpExp interface {
	GetHeaders(arg any) ([]string, error)
	Import(bytes []byte, handler ImpHandler) error
	ImportWithUrl(url string, handler ImpHandler) error
	ExportSync(handler ExpHandler) ([]byte, error)
	ExportAsync(handler ExpHandler) (string, error)
	GetProcess(taskid string) float32
}

type Task struct {
	TaskId     string
	Title      string
	AddTime    time.Time
	StartTime  time.Time
	FinishTime time.Time
	Handler    ExpHandler
	Process    float32
	Expcount   int
	TotalSize  int
	CostTime   time.Duration
	Timer      int //12c
}

func (t *Task) timerAndExpire() bool {
	if t.Timer >= 11 {
		return true
	}
	t.Timer++
	return false
}
