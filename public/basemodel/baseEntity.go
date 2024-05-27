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
	GetToZeroCols() []string
	GetConditions() ([]string, string, []interface{})
}

type Entity struct {
	Deleted   int           `json:"-"`
	whereCols []string      `gorm:"-" json:"-"` // 扩展条件字段名
	whereCond []string      `gorm:"-" json:"-"` // 扩展条件内容
	condVals  []interface{} `gorm:"-" json:"-"` // 扩展条件值
}

type BaseEntity struct {
	CreateBy   int64    `json:"createBy,string,omitempty"` //创建人
	CreateTime BaseTime `json:"createTime,omitempty"`      //创建时间
	UpdateBy   int64    `json:"updateBy,string,omitempty"` //修改人
	UpdateTime BaseTime `json:"updateTime,omitempty"`      //修改时间
	Entity
	toNullCols []string `gorm:"-" json:"-"` //更新时需要置空(null)的字段列表
	toZeroCols []string `gorm:"-" json:"-"` //更新时需要设0值(0、"")的字段列表，仅支持int和string类型字段
}

type BaseEntityQuery struct {
	CurPage  int `gorm:"-" form:"curPage" json:"curPage,omitempty"`   //第几页
	PageSize int `gorm:"-" form:"pageSize" json:"pageSize,omitempty"` //数量
	Entity
	orders []string `gorm:"-" json:"-"`
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

// 设置需要设0值(0、"")的字段列表，仅支持int类型字段
func (p *BaseEntity) SetToZeroCols(cols ...string) {
	p.toZeroCols = append(p.toZeroCols, cols...)
}

func (p *BaseEntity) GetToZeroCols() []string {
	return p.toZeroCols
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
func (p *Entity) AddLikeAll(column, conditionVal string) *Entity {
	if len(strings.TrimSpace(column)) > 0 && len(strings.TrimSpace(conditionVal)) > 0 {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" LIKE ?")
		p.condVals = append(p.condVals, "%"+conditionVal+"%")
	}
	return p
}

// AND column LIKE %?
func (p *Entity) AddLikeLeft(column string, conditionVal string) *Entity {
	if len(strings.TrimSpace(column)) > 0 && len(strings.TrimSpace(conditionVal)) > 0 {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" LIKE ?")
		p.condVals = append(p.condVals, "%"+conditionVal)
	}
	return p
}

// AND column LIKE ?%
func (p *Entity) AddLikeRight(column string, conditionVal string) *Entity {
	if len(strings.TrimSpace(column)) > 0 && len(strings.TrimSpace(conditionVal)) > 0 {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" LIKE ?")
		p.condVals = append(p.condVals, conditionVal+"%")
	}
	return p
}

// AND column <> ?
func (p *Entity) AddNot(column string, value interface{}) *Entity {
	return p.buildCompare(column, value, "<>")
}

// AND column < ?
func (p *Entity) AddLt(column string, value interface{}) *Entity {
	return p.buildCompare(column, value, "<")
}

// AND column <= ?
func (p *Entity) AddLe(column string, value interface{}) *Entity {
	return p.buildCompare(column, value, "<=")
}

// AND column > ?
func (p *Entity) AddGt(column string, value interface{}) *Entity {
	return p.buildCompare(column, value, ">")
}

// AND column >= ?
func (p *Entity) AddGe(column string, value interface{}) *Entity {
	return p.buildCompare(column, value, ">=")
}

// 构建比较方法
func (p *Entity) buildCompare(column string, value interface{}, cond string) *Entity {
	if value != nil {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" "+cond+" ?")
		p.condVals = append(p.condVals, value)
	}
	return p
}

// AND column IN (?)
func (p *Entity) AddIn(column string, conditionVal ...interface{}) *Entity {
	if len(conditionVal) > 0 {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" IN ?")
		p.condVals = append(p.condVals, conditionVal)
	}
	return p
}

// AND column NOT IN (?)
func (p *Entity) AddNotIn(column string, conditionVal ...interface{}) *Entity {
	if len(conditionVal) > 0 {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" NOT IN ?")
		p.condVals = append(p.condVals, conditionVal)
	}
	return p
}

// AND column BETWEEN ? AND ?
func (p *Entity) AddBetween(column string, condStart, condEnd interface{}) *Entity {
	if condStart != nil && condEnd != nil {
		p.whereCols = append(p.whereCols, column)
		p.whereCond = append(p.whereCond, column+" BETWEEN ? AND ?")
		p.condVals = append(p.condVals, condStart, condEnd)
	}
	return p
}

// AND column IS NULL
func (p *Entity) AddIsNull(column string) *Entity {
	p.whereCols = append(p.whereCols, column)
	p.whereCond = append(p.whereCond, column+" IS NULL")
	return p
}

// AND column IS NOT NULL
func (p *Entity) AddIsNotNull(column string) *Entity {
	p.whereCols = append(p.whereCols, column)
	p.whereCond = append(p.whereCond, column+" IS NOT NULL")
	return p
}

func (p *Entity) GetConditions() ([]string, string, []interface{}) {
	if p.whereCols == nil || len(p.whereCols) <= 0 {
		return nil, "", nil
	}
	return p.whereCols, strings.Join(p.whereCond, " AND "), p.condVals
}
