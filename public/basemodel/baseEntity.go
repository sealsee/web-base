package basemodel

import "github.com/sealsee/web-base/public/ds/page"

type BaseEntity struct {
	CreateBy   *int64    `json:"createBy,omitempty" db:"create_by"`     //创建人
	CreateTime *BaseTime `json:"createTime,omitempty" db:"create_time"` //创建时间
	UpdateBy   *int64    `json:"updateBy,omitempty" db:"update_by"`     //修改人
	UpdateTime *BaseTime `json:"updateTime,omitempty" db:"update_time"` //修改时间
}

type BaseEntityQuery struct {
	CurPage  int    `gorm:"-" form:"curPage" json:"curPage,omitempty" default:"1"` //第几页
	PageSize int    `gorm:"-" form:"pageSize" json:"pageSize,omitempty"`           //数量
	OrderBy  string `gorm:"-" form:"orderBy" json:"orderBy,omitempty"`             //排序字段
	IsAsc    string `gorm:"-" form:"isAsc" json:"isAsc,omitempty"`                 //排序规则  降序desc   asc升序
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
