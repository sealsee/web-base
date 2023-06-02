package dbf

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/sealsee/web-base/public/IOFile"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/utils/file/dbf/godbf"
	"go.uber.org/zap"
)

const (
	fileEncoding = "GBK"
	max_page     = 2000
)

type dbf struct {
}

func (d *dbf) GetHeaders(arg any) ([]string, error) {
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

	dbfTable, err := godbf.NewFromByteArray(bs, fileEncoding)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	fieldNames := dbfTable.FieldNames()
	if fieldNames == nil || len(fieldNames) <= 0 {
		return nil, errors.New("fieldNames invalid")
	}
	return nil, nil
}

func (d *dbf) Import(bytes []byte, handler ImpHandler) error {
	if bytes == nil || handler == nil {
		return errors.New("params is invalid")
	}

	dbfTable, err := godbf.NewFromByteArray(bytes, fileEncoding)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	fieldNames := dbfTable.FieldNames()
	if fieldNames == nil || len(fieldNames) <= 0 {
		return errors.New("fieldNames invalid")
	}

	handler.Headers(fieldNames)
	rowCount := dbfTable.NumberOfRecords()
	zap.L().Info("rowCount:", zap.Int("", rowCount))

	for i := 0; i < rowCount; i++ {
		row := map[string]string{}
		for _, fname := range fieldNames {
			val, _ := dbfTable.FieldValueByName(i, fname)
			row[fname] = val
		}
		handler.Row(&row)

	}

	handler.After()
	return nil
}

func (d *dbf) ImportWithUrl(url string, handler ImpHandler) error {
	if url == "" || handler == nil {
		return errors.New("params is invalid")
	}
	bytes, err := IOFile.GetConfig().Download(url)
	if err != nil {
		return errors.New("download with " + url + " err")
	}

	return d.Import(bytes, handler)
}

func (d *dbf) Export(handler ExpHandler) ([]byte, error) {
	if handler == nil || handler.Fields() == nil || len(handler.Fields()) < 1 {
		return nil, errors.New("invalid params")
	}

	dbfTable := godbf.New(fileEncoding)

	for _, v := range handler.Fields() {
		dbfname := v.DBFname
		if v.Type == S {
			len := v.Length
			if v.Length <= 0 {
				len = 32
			}
			dbfTable.AddTextField(dbfname, byte(len))
		} else if v.Type == N {
			dbfTable.AddNumberField(dbfname, 32, 0)
		} else if v.Type == D {
			dbfTable.AddFloatField(dbfname, 32, 3)
		} else if v.Type == DT {
			dbfTable.AddDateField(dbfname)
		} else {
			dbfTable.AddTextField(dbfname, 64)
		}
	}

	page := page.NewExportPage()
	for i := 0; i < max_page; i++ {
		list := handler.Rows(page)
		len := len(list)
		if list == nil || len < 1 {
			break
		}
		for _, v := range list {
			_createRow(dbfTable, v, handler.Fields())
		}

		if !page.NextPage() {
			break
		}
	}

	handler.Finish("")
	return godbf.GetDbfFileData(dbfTable), nil
}

func _createRow(table *godbf.DbfTable, row map[string]interface{}, fields []DBFField) {
	rowId, _ := table.AddNewRecord()
	for _, h := range fields {
		val := row[h.Dname]
		table.SetFieldValueByName(rowId, h.DBFname, _valCvt(val))
	}
}

func _valCvt(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	case time.Time:
		t, _ := value.(time.Time)
		key = t.String()
		key = strings.Replace(key, " +0800 CST", "", 1)
		key = strings.Replace(key, " +0000 UTC", "", 1)
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
