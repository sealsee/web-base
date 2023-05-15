package basemodel

import (
	"strings"
	"time"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds/page"
)

type IQuery interface {
	GetOrders() string
}

type BaseEntity struct {
	Deleted    int      `json:"-"`
	CreateBy   int64    `json:"createBy,string,omitempty"` //创建人
	CreateTime BaseTime `json:"createTime,omitempty"`      //创建时间
	UpdateBy   int64    `json:"updateBy,string,omitempty"` //修改人
	UpdateTime BaseTime `json:"updateTime,omitempty"`      //修改时间
}

type BaseEntityQuery struct {
	Deleted  int      `json:"-"`
	CurPage  int      `gorm:"-" form:"curPage" json:"curPage,omitempty"`   //第几页
	PageSize int      `gorm:"-" form:"pageSize" json:"pageSize,omitempty"` //数量
	orders   []string `gorm:"-" json:"-"`
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
