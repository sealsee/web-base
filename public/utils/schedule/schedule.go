package schedule

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

var crontab *cron.Cron
var schedule_tables = []ScheduleTask{}

type ScheduleTask struct {
	EntryID int    `json:"entry_id,omitempty"`
	JobName string `json:"job_name,omitempty"`
	CronExp string `json:"cron_exp,omitempty"`
}

func init() {
	crontab = cron.New(cron.WithSeconds())
	crontab.Start()
}

func AddJob(jobName, cronExp string, task func()) {
	if cronExp == "" || task == nil {
		return
	}

	entryID, _ := crontab.AddFunc(cronExp, task)
	fmt.Println("entryID:", entryID)
	schTask := ScheduleTask{EntryID: int(entryID), JobName: jobName, CronExp: cronExp}
	schedule_tables = append(schedule_tables, schTask)
}

func ListTask() []ScheduleTask {
	return schedule_tables
}
