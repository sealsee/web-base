package basemodel

import (
	"strings"
	"time"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds/page"
)

type BaseEntity struct {
	Deleted    int      `json:"-"`
	CreateBy   int64    `json:"createBy,omitempty"`   //创建人
	CreateTime BaseTime `json:"createTime,omitempty"` //创建时间
	UpdateBy   int64    `json:"updateBy,omitempty"`   //修改人
	UpdateTime BaseTime `json:"updateTime,omitempty"` //修改时间
}

type BaseEntityQuery struct {
	Deleted  int      `json:"-"`
	CurPage  int      `gorm:"-" form:"curPage" json:"curPage,omitempty"`   //第几页
	PageSize int      `gorm:"-" form:"pageSize" json:"pageSize,omitempty"` //数量
	Orders   []string `gorm:"-" json:"-"`
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
	if p.Orders == nil {
		p.Orders = make([]string, 0)
	}
	p.Orders = append(p.Orders, column+" ASC")
}

func (p *BaseEntityQuery) AddOrderDesc(column string) {
	if p.Orders == nil {
		p.Orders = make([]string, 0)
	}
	p.Orders = append(p.Orders, column+" DESC")
}

func (p *BaseEntityQuery) GetOrders() string {
	if p.Orders == nil || len(p.Orders) <= 0 {
		return ""
	}
	return strings.Join(p.Orders, ",")
}
