package basemodel

import "time"

type BaseEntity struct {
	CreateBy   int64     `json:"createBy" db:"create_by"`     //创建人
	CreateTime *BaseTime `json:"createTime" db:"create_time"` //创建时间
	UpdateBy   int64     `json:"updateBy" db:"update_by"`     //修改人
	UpdateTime *BaseTime `json:"updateTime" db:"update_time"` //修改时间
}

type BaseEntityQuery struct {
	DataScope string `swaggerignore:"true"`         // 数据范围控制
	ExpType   int    `form:"expType" default:"1"`   // 导出类型 1-excel 2-dbf ...
	CurPage   int    `form:"curPage" default:"1"`   //第几页
	PageSize  int    `form:"pageSize" default:"10"` //数量
	OrderBy   string `form:"orderBy" `              //排序字段
	IsAsc     string `form:"isAsc" `                //排序规则  降序desc   asc升序
}

type BaseEntityAdd struct {
	CreateBy   int64    `db:"create_by"`
	CreateTime BaseTime `db:"create_time"`
	UpdateBy   int64    `db:"update_by"`
	UpdateTime BaseTime `db:"update_time"`
}

type BaseEntityEdit struct {
	UpdateBy   int64    `db:"update_by"`
	UpdateTime BaseTime `db:"update_time"`
}

func (b *BaseEntityAdd) SetCreateBy(userId int64) {
	b.CreateBy = userId
	kt := BaseTime(time.Now())
	b.CreateTime = kt
	b.UpdateBy = userId
	b.UpdateTime = kt
}

func (b *BaseEntityEdit) SetUpdateBy(userId int64) {
	b.UpdateBy = userId
	b.UpdateTime = BaseTime(time.Now())
}
