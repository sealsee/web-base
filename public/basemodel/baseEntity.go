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

// 设置需要置空(set NULL)的字段，传数据库字段名
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

// AND column LIKE %?%
func (p *BaseEntityQuery) AddLikeAll(column, conditionVal string) *BaseEntityQuery {
	if len(strings.TrimSpace(column)) > 0 && len(strings.TrimSpace(conditionVal)) > 0 {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" LIKE ?")
		p.condVals = append(p.condVals, "%"+conditionVal+"%")
	}
	return p
}

// AND column LIKE %?
func (p *BaseEntityQuery) AddLikeLeft(column string, conditionVal string) *BaseEntityQuery {
	if len(strings.TrimSpace(column)) > 0 && len(strings.TrimSpace(conditionVal)) > 0 {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" LIKE ?")
		p.condVals = append(p.condVals, "%"+conditionVal)
	}
	return p
}

// AND column LIKE ?%
func (p *BaseEntityQuery) AddLikeRight(column string, conditionVal string) *BaseEntityQuery {
	if len(strings.TrimSpace(column)) > 0 && len(strings.TrimSpace(conditionVal)) > 0 {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" LIKE ?")
		p.condVals = append(p.condVals, conditionVal+"%")
	}
	return p
}

// AND column <> ?
func (p *BaseEntityQuery) AddNot(column string, value interface{}) *BaseEntityQuery {
	return p.buildCompare(column, value, "<>")
}

// AND column < ?
func (p *BaseEntityQuery) AddLt(column string, value interface{}) *BaseEntityQuery {
	return p.buildCompare(column, value, "<")
}

// AND column <= ?
func (p *BaseEntityQuery) AddLe(column string, value interface{}) *BaseEntityQuery {
	return p.buildCompare(column, value, "<=")
}

// AND column > ?
func (p *BaseEntityQuery) AddGt(column string, value interface{}) *BaseEntityQuery {
	return p.buildCompare(column, value, ">")
}

// AND column >= ?
func (p *BaseEntityQuery) AddGe(column string, value interface{}) *BaseEntityQuery {
	return p.buildCompare(column, value, ">=")
}

// 构建比较方法
func (p *BaseEntityQuery) buildCompare(column string, value interface{}, cond string) *BaseEntityQuery {
	if value != nil {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" "+cond+" ?")
		p.condVals = append(p.condVals, value)
	}
	return p
}

// AND column IN (?)
func (p *BaseEntityQuery) AddIn(column string, conditionVal ...interface{}) *BaseEntityQuery {
	if len(conditionVal) > 0 {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" IN ?")
		p.condVals = append(p.condVals, conditionVal)
	}
	return p
}

// AND column NOT IN (?)
func (p *BaseEntityQuery) AddNotIn(column string, conditionVal ...interface{}) *BaseEntityQuery {
	if len(conditionVal) > 0 {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" NOT IN ?")
		p.condVals = append(p.condVals, conditionVal)
	}
	return p
}

// AND column BETWEEN ? AND ?
func (p *BaseEntityQuery) AddBetween(column string, condStart, condEnd interface{}) *BaseEntityQuery {
	if condStart != nil && condEnd != nil {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" BETWEEN ? AND ?")
		p.condVals = append(p.condVals, condStart, condEnd)
	}
	return p
}

// AND column IS NULL
func (p *BaseEntityQuery) AddIsNull(column string) *BaseEntityQuery {
	p.whereCols = append(p.whereCols, column)
	p.whereCond = append(p.whereCond, column+" IS NULL")
	return p
}

// AND column IS NOT NULL
func (p *BaseEntityQuery) AddIsNotNull(column string) *BaseEntityQuery {
	p.whereCols = append(p.whereCols, column)
	p.whereCond = append(p.whereCond, column+" IS NOT NULL")
	return p
}

func (p *BaseEntityQuery) GetConditions() ([]string, string, []interface{}) {
	if p.whereCols == nil || len(p.whereCols) <= 0 {
		return nil, "", nil
	}
	return p.whereCols, strings.Join(p.whereCond, " AND "), p.condVals
}
