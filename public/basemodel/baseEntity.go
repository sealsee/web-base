package basemodel

import (
	"strings"
	"time"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds/page"
)

type IQuery interface {
	GetOrders() string
	GetConditions() ([]string, string, []interface{})
}

type IEntidy interface {
	GetToNullCols() []string
}

type Entity struct {
	Deleted int `json:"-"`
}

type BaseEntity struct {
	CreateBy   int64    `json:"createBy,string,omitempty"` //创建人
	CreateTime BaseTime `json:"createTime,omitempty"`      //创建时间
	UpdateBy   int64    `json:"updateBy,string,omitempty"` //修改人
	UpdateTime BaseTime `json:"updateTime,omitempty"`      //修改时间
	Entity
	toNullCols []string `gorm:"-" json:"-"` //更新时需要置空的字段列表
}

type BaseEntityQuery struct {
	CurPage  int `gorm:"-" form:"curPage" json:"curPage,omitempty"`   //第几页
	PageSize int `gorm:"-" form:"pageSize" json:"pageSize,omitempty"` //数量
	Entity
	orders    []string      `gorm:"-" json:"-"`
	whereCols []string      `gorm:"-" json:"-"` // 扩展条件字段名
	whereCond []string      `gorm:"-" json:"-"` // 扩展条件内容
	condVals  []interface{} `gorm:"-" json:"-"` // 扩展条件值
}

func (p *BaseEntity) SetCreateBy(createBy int64) {
	p.CreateBy = createBy
	p.CreateTime = BaseTime(time.Now())
}

func (p *BaseEntity) SetUpdateBy(updateBy int64) {
	p.UpdateBy = updateBy
	p.UpdateTime = BaseTime(time.Now())
}

func (p *BaseEntity) SetDeleteBy(deleteBy int64) {
	p.Deleted = common.Deleted
	p.UpdateBy = deleteBy
	p.UpdateTime = BaseTime(time.Now())
}

func (p *BaseEntity) SetNullableCols(cols ...string) {
	p.toNullCols = append(p.toNullCols, cols...)
}

func (p *BaseEntity) GetToNullCols() []string {
	return p.toNullCols
}

func (p *BaseEntityQuery) GetPage() *page.Page {
	page := page.NewPage()
	if p.CurPage > 0 {
		page.CurPage = p.CurPage
	}
	if p.PageSize > 0 {
		page.PageSize = p.PageSize
	}
	return page
}

func (p *BaseEntityQuery) AddOrderAsc(column string) {
	if p.orders == nil {
		p.orders = make([]string, 0)
	}
	p.orders = append(p.orders, column+" ASC")
}

func (p *BaseEntityQuery) AddOrderDesc(column string) {
	if p.orders == nil {
		p.orders = make([]string, 0)
	}
	p.orders = append(p.orders, column+" DESC")
}

func (p *BaseEntityQuery) GetOrders() string {
	if p.orders == nil || len(p.orders) <= 0 {
		return ""
	}
	return strings.Join(p.orders, ",")
}

// 左LIKE %?
func (p *BaseEntityQuery) likeL(colum string) string {
	if len(strings.TrimSpace(colum)) == 0 {
		return ""
	}
	return "%" + colum
}

// 右LIKE ?%
func (p *BaseEntityQuery) likeR(colum string) string {
	if len(strings.TrimSpace(colum)) == 0 {
		return ""
	}
	return colum + "%"
}

// 全LIKE %?%
func (p *BaseEntityQuery) likeA(colum string) string {
	if len(strings.TrimSpace(colum)) == 0 {
		return ""
	}
	return "%" + colum + "%"
}

func (p *BaseEntityQuery) AddLikeAll(colum, conditionVal string) {
	if len(strings.TrimSpace(colum)) > 0 && len(strings.TrimSpace(conditionVal)) > 0 {
		p.whereCols = append(p.whereCols, colum)
		p.whereCond = append(p.whereCond, colum+" like ?")
		p.condVals = append(p.condVals, p.likeA(conditionVal))
	}
}

func (p *BaseEntityQuery) AddLikeLeft(colum string, conditionVal string) {
	if len(strings.TrimSpace(colum)) > 0 && len(strings.TrimSpace(conditionVal)) > 0 {
		p.whereCols = append(p.whereCols, colum)
		p.whereCond = append(p.whereCond, colum+" like ?")
		p.condVals = append(p.condVals, p.likeL(conditionVal))
	}
}

func (p *BaseEntityQuery) AddLikeRight(colum string, conditionVal string) {
	if len(strings.TrimSpace(colum)) > 0 && len(strings.TrimSpace(conditionVal)) > 0 {
		p.whereCols = append(p.whereCols, colum)
		p.whereCond = append(p.whereCond, colum+" like ?")
		p.condVals = append(p.condVals, p.likeR(conditionVal))
	}
}

func (p *BaseEntityQuery) GetConditions() ([]string, string, []interface{}) {
	if p.whereCols == nil || len(p.whereCols) <= 0 {
		return nil, "", nil
	}
	return p.whereCols, strings.Join(p.whereCond, " and "), p.condVals
}
