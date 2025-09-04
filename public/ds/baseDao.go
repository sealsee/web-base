package ds

import (
	"github.com/sealsee/web-base/public/basemodel"
)

type BaseModel[T basemodel.BaseEntity] struct {
}

// TODO 暂未调通，先不用
func (model *BaseModel[T]) GetById(id int64) *T {
	if id < 1 {
		return nil
	}
	var data T
	where := new(T)
	// where.ID = id
	// where.Deleted = common.Normal
	res := GetGDB().Find(&data, where)
	if res.Error != nil {
		panic(res.Error)
	}
	if res.RowsAffected < 1 {
		return nil
	}
	return &data
}

// func (syncAppStuprof *BaseModel) Count(q *BaseModel) int {
// 	q.Deleted = common.Normal
// 	return query.ExecGetQueryCount[EaSyncStuProfQuery, EaSyncStuProf](q)
// }
