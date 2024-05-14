package excel

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/sealsee/web-base/public/IOFile"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

var (
	// first        = 65
	// heads_symbol []string
	id          = "rid"
	max_page    = 2000
	isStart     atomic.Bool
	tasks       = make(chan string, 100)
	tasks_entry sync.Map
)

// func init() {
// 	for i := first; i < first+26; i++ {
// 		heads_symbol = append(heads_symbol, fmt.Sprintf("%c", i))
// 	}
// }

type excel struct {
}

func (e *excel) GetHeaders(arg any) ([]string, error) {
	if arg == nil {
		return nil, nil
	}

	var bs []byte
	switch x := arg.(type) {
	case string:
		b, err := IOFile.GetConfig().Download(x)
		if err != nil {
			return nil, errors.New("download with " + x + " err")
		}
		bs = b
	case []byte:
		bs = x
	default:
		return nil, errors.New("arg is invalid,need url or []byte")
	}

	file, err := excelize.OpenReader(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	defSheet := file.GetSheetList()[0]
	if defSheet == "" {
		zap.L().Error("no sheet default")
		return nil, errors.New("no sheet default")
	}

	rows, _ := file.Rows(defSheet)
	if rows.Next() {
		return rows.Columns()
	}

	return nil, nil
}

func (e *excel) Import(bs []byte, handler ImpHandler) error {
	if bs == nil || len(bs) <= 0 || handler == nil {
		return errors.New("params is invalid")
	}

	file, err := excelize.OpenReader(bytes.NewReader(bs))
	if err != nil {
		return err
	}
	defer file.Close()

	defSheet := file.GetSheetList()[0]
	if defSheet == "" {
		zap.L().Error("no sheet default")
		return errors.New("no sheet default")
	}

	var headers []string
	headerMap := map[int]string{}
	rows, _ := file.Rows(defSheet)
	if rows.Next() {
		i := 0
		cols, _ := rows.Columns()
		for _, v := range cols {
			i++
			headerMap[i] = v
			headers = append(headers, v)
		}
		handler.Headers(headers)
	}

	rowIdx := 1
	for rows.Next() {
		cols, _ := rows.Columns()
		if len(cols) < 1 {
			continue
		}
		j := 0
		row := map[string]string{}
		row[id] = strconv.Itoa(rowIdx)
		for _, v := range cols {
			j++
			header := headerMap[j]
			row[header] = v
		}
		status := handler.Row(&row)
		if status == Exit {
			break
		}
		rowIdx++
	}

	handler.After()
	return nil
}

func (e *excel) ImportWithUrl(url string, handler ImpHandler) error {
	if url == "" || handler == nil {
		return errors.New("params is invalid")
	}
	bytes, err := IOFile.GetConfig().Download(url)
	if err != nil {
		return errors.New("download with " + url + " err")
	}
	return e.Import(bytes, handler)
}

func (e *excel) ExportSync(handler ExpHandler) ([]byte, error) {
	return _export(&Task{Handler: handler})
}

func (e *excel) ExportAsync(handler ExpHandler) (string, error) {
	if handler == nil || handler.HeaderColumn() == nil || len(handler.HeaderColumn()) < 1 {
		return "", errors.New("invalid params")
	}

	tid, err := _addTask(handler)
	if err != nil {
		return "", err
	}

	if !isStart.Load() {
		isStart.Store(true)
		go func() {
			for {
				id := <-tasks
				v, ok := tasks_entry.Load(id)
				if !ok {
					continue
				}
				t, _ := v.(*Task)
				start := time.Now()
				_runAsyncExp(t)
				cost := time.Since(start)
				t.CostTime = cost
			}
		}()

		go func() {
			for {
				ticker := time.NewTicker(time.Second * 5)
				<-ticker.C
				tasks_entry.Range(func(key, value any) bool {
					t, _ := value.(*Task)
					// fmt.Println(key, "----process:", t.Process, t.Expcount, t.TotalSize, t.CostTime)
					if t.timerAndExpire() {
						tasks_entry.Delete(t.TaskId)
					}
					return true
				})
			}
		}()
	}

	return tid, nil
}

func _runAsyncExp(task *Task) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("", zap.Any("", err))
		}
	}()

	bs, err := _export(task)
	if err != nil {
		zap.L().Error(err.Error())
	}
	url, err := IOFile.GetConfig().Upload(bytes.NewReader(bs), "", "xlsx", true)
	if err != nil {
		zap.L().Error(err.Error())
	}
	//TODO
	task.Handler.Finish(url)
}

func _export(task *Task) ([]byte, error) {
	handler := task.Handler
	sheetName := "Sheet1"
	file := excelize.NewFile()
	defer file.Close()
	file.NewSheet(sheetName)

	style, err := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "gradient", Pattern: 1, Color: []string{"#FF0000"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	headers, cols := _headerColumn(handler.HeaderColumn())
	err = file.SetSheetRow(sheetName, "A1", &headers)
	file.SetRowStyle(sheetName, 1, 1, style)
	if err != nil {
		zap.L().Error(err.Error())
	}

	rid := 2
	rowCount := 0
	page := page.NewExportPage()
	for i := 0; i < max_page; i++ {
		list := handler.Rows(page)
		len := len(list)
		if list == nil || len < 1 {
			break
		}

		rowCount += len
		for _, v := range list {
			ary := _map2array(v, &cols)
			err := file.SetSheetRow(sheetName, "A"+strconv.Itoa(rid), &ary)
			if err != nil {
				zap.L().Error(err.Error())
			}
			rid++
		}

		task.Process = float32(rowCount / page.TotalSize)
		task.Expcount = rowCount
		task.TotalSize = page.TotalSize

		if !page.NextPage() {
			break
		}
	}

	task.Process = 100

	buf, err := file.WriteToBuffer()
	return buf.Bytes(), err
}

func (e *excel) GetProcess(taskid string) float32 {
	if taskid == "" {
		return 0
	}
	value, ok := tasks_entry.Load(taskid)
	if !ok {
		return 100
	}
	t, _ := value.(*Task)
	return t.Process
}

func _map2array(row map[string]interface{}, cols *[]string) []interface{} {
	vals := make([]interface{}, 0, len(*cols))
	for _, v := range *cols {
		value := row[v]
		vals = append(vals, value)
	}
	return vals
}

func _headerColumn(fs []string) ([]string, []string) {
	if fs == nil || len(fs) < 1 {
		return nil, nil
	}
	var headers []string
	var cols []string
	for _, v := range fs {
		if v == "" {
			continue
		}
		ary := strings.Split(v, ",")
		if len(ary) != 2 {
			continue
		}
		header := ary[0]
		column := ary[1]
		if column == "" {
			continue
		}
		if header == "" {
			header = column
		}

		headers = append(headers, header)
		cols = append(cols, column)
	}
	return headers, cols
}

func _addTask(handler ExpHandler) (string, error) {
	u := uuid.New()
	id := strconv.Itoa((int(u.ID())))
	select {
	case tasks <- id:
		tasks_entry.Store(id, &Task{TaskId: id, Title: handler.Title(), Handler: handler, AddTime: time.Now()})
		return fmt.Sprint(id), nil
	default:
		zap.L().Error("export task is full")
		return "", errors.New("export task is full")
	}
}

func GetExcelTask() sync.Map {
	return tasks_entry
}
