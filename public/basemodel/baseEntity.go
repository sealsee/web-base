package basemodel

type BaseEntity struct {
	CreateBy   *int64    `json:"createBy"`   //创建人
	CreateTime *BaseTime `json:"createTime"` //创建时间
	UpdateBy   *int64    `json:"updateBy"`   //修改人
	UpdateTime *BaseTime `json:"updateTime"` //修改时间
}

type BaseEntityQuery struct {
	DataScope string `swaggerignore:"true"`         // 数据范围控制
	ExpType   int    `form:"expType" default:"1"`   // 导出类型 1-excel 2-dbf ...
	CurPage   int    `form:"curPage" default:"1"`   //第几页
	PageSize  int    `form:"pageSize" default:"10"` //数量
	OrderBy   string `form:"orderBy" `              //排序字段
	IsAsc     string `form:"isAsc" `                //排序规则  降序desc   asc升序
}
