package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sealsee/web-base/public/cst/httpStatus"
	"github.com/sealsee/web-base/public/ds/page"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/errs"
)

const (
	OK   = "操作成功"
	ERR  = "系统繁忙,请稍后再试!"
	PAGE = "page"
	LIST = "list"
	OBJ  = "obj"
)

type JsonResult struct {
	Success bool                   `json:"success"`
	Code    string                 `json:"code"`
	Msg     string                 `json:"msg"`
	Data    map[string]interface{} `json:"data"`
	c       *gin.Context
}

func NewJsonResult(c *gin.Context) *JsonResult {
	return &JsonResult{
		Success: true,
		Code:    strconv.Itoa(int(httpStatus.Success)),
		Msg:     OK,
		Data:    make(map[string]interface{}, 1),
		c:       c}
}

func (json *JsonResult) SetSysErr() *JsonResult {
	json.Code = strconv.Itoa(int(httpStatus.Error))
	json.Msg = ERR
	json.Success = false
	return json
}

func (json *JsonResult) SetErr(err string) *JsonResult {
	json.Code = strconv.Itoa(int(httpStatus.Error))
	json.Msg = err
	json.Success = false
	return json
}

func (json *JsonResult) SetErrs(err errs.ERROR) *JsonResult {
	json.Code = err[0]
	json.Msg = err[1]
	json.Success = false
	return json
}

func (json *JsonResult) SetErrsWithArgs(err errs.ERROR, args ...string) *JsonResult {
	json.Code = err[0]
	json.Msg = fmt.Sprintf(err[1], args)
	json.Success = false
	return json
}

func (json *JsonResult) checkDataInit() {
	if json.Data == nil {
		json.Data = make(map[string]interface{}, 1)
	}
}

func (json *JsonResult) AddData(key string, ary interface{}) *JsonResult {
	json.checkDataInit()
	if key != "" && ary != nil {
		json.Data[key] = ary
	}
	return json
}

func (json *JsonResult) SetPage(ary interface{}) *JsonResult {
	json.checkDataInit()
	json.Data[PAGE] = ary
	return json
}

func (json *JsonResult) SetList(ary interface{}) *JsonResult {
	json.checkDataInit()
	json.Data[LIST] = ary
	return json
}

func (json *JsonResult) SetPageList(ary interface{}, page *page.Page) *JsonResult {
	json.checkDataInit()
	json.Data[LIST] = ary
	json.Data[PAGE] = page
	return json
}

func (json *JsonResult) SetObj(ary interface{}) *JsonResult {
	json.checkDataInit()
	json.Data[OBJ] = ary
	return json
}

func (json *JsonResult) Render() {
	json.c.JSON(http.StatusOK, json)
}

func (json *JsonResult) RenderWithCode(code int) {
	json.c.JSON(code, json)
}
