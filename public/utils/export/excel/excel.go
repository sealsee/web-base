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
	"github.com/sealsee/web-base/public/utils/export/internal"
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

type Excel struct {
}

func ExportExcel(dataList [][]interface{}) (data []byte) {
	f := excelize.NewFile()
	defer f.Close()
	for i, row := range dataList {
		if i == 0 {
			f.SetSheetRow("Sheet1", "A1", &row)
		} else {
			f.SetSheetRow("Sheet1", "A"+strconv.Itoa(i+1), &row)
		}
	}
	buffer, _ := f.WriteToBuffer()
	return buffer.Bytes()
}

func (e *Excel) Import(bs []byte, handler internal.ImpHandler) error {
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
		j := 0
		row := map[string]string{}
		row[id] = strconv.Itoa(rowIdx)
		for _, v := range cols {
			j++
			header := headerMap[j]
			row[header] = v
		}
		handler.Row(&row)
		rowIdx++
	}

	handler.After()
	return nil
}

func (e *Excel) ImportWithUrl(url string, handler internal.ImpHandler) error {
	if url == "" || handler == nil {
		return errors.New("params is invalid")
	}
	bytes, err := IOFile.GetConfig().Download(url)
	if err != nil {
		return errors.New("download with url err")
	}
	return e.Import(bytes, handler)
}

func (e *Excel) ExportSync(handler internal.ExpHandler) ([]byte, error) {
	return _export(&internal.Task{Handler: handler})
}

func (e *Excel) ExportAsync(handler internal.ExpHandler) (string, error) {
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
				t, _ := v.(*internal.Task)
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
					t, _ := value.(*internal.Task)
					// fmt.Println(key, "----process:", t.Process, t.Expcount, t.TotalSize, t.CostTime)
					if t.TimerAndExpire() {
						tasks_entry.Delete(t.TaskId)
					}
					return true
				})
			}
		}()
	}

	return tid, nil
}

func _runAsyncExp(task *internal.Task) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("", zap.Any("", err))
		}
	}()

	_, err := _export(task)
	if err != nil {
		zap.L().Error(err.Error())
	}

	//TODO
	task.Handler.Finish("")
}

func _export(task *internal.Task) ([]byte, error) {
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

func (e *Excel) GetProcess(taskid string) float32 {
	if taskid == "" {
		return 0
	}
	value, ok := tasks_entry.Load(taskid)
	if !ok {
		return 100
	}
	t, _ := value.(*internal.Task)
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

func _addTask(handler internal.ExpHandler) (string, error) {
	u := uuid.New()
	id := strconv.Itoa((int(u.ID())))
	select {
	case tasks <- id:
		tasks_entry.Store(id, &internal.Task{TaskId: id, Title: handler.Title(), Handler: handler, AddTime: time.Now()})
		return fmt.Sprint(id), nil
	default:
		zap.L().Error("export task is full")
		return "", errors.New("export task is full")
	}
}
