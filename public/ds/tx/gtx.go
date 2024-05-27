package tx

import (
	"errors"
	"fmt"

	"github.com/sealsee/web-base/public/basemodel"
	"github.com/sealsee/web-base/public/utils/jsonUtils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var gormdb *gorm.DB

func InitGTx(gdb *gorm.DB) {
	gormdb = gdb
}

// 用到gorm的方法签名
type GTx interface {
	Create(value interface{}) (tx *gorm.DB)
	CreateInBatches(value interface{}, batchSize int) (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	Updates(values interface{}) (tx *gorm.DB)
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
	Raw(sql string, values ...interface{}) (tx *gorm.DB)
	Exec(sql string, values ...interface{}) (tx *gorm.DB)
	Model(value interface{}) (tx *gorm.DB)
	Omit(columns ...string) (tx *gorm.DB)
}

// 扩展类型
type Db struct {
	*gorm.DB
}

// 扩展方法-更新，支持自定义条件、支持所有字段置NULL、支持数字类型置零、支持条件传参
func (db *Db) UpdatesNewWithCondition(uWhere interface{}, uData basemodel.IEntidy, condition interface{}, args ...interface{}) (tx *gorm.DB) {
	dataMap, _ := jsonUtils.StructToDbMap(uData)
	// fmt.Printf("---> uData map: %v \n", dataMap)
	for _, col := range uData.GetToNullCols() {
		if _, ok := dataMap[col]; ok && dataMap[col] != "" {
			continue
		}
		dataMap[col] = nil
	}
	for _, col := range uData.GetToZeroCols() {
		if _, ok := dataMap[col]; ok {
			continue
		}
		dataMap[col] = 0
	}
	// fmt.Printf("---> dataMap: %v \n", dataMap)
	whereMap, cond, condArgs := convertWhereEntity(uWhere)
	gdb := db.Model(&uWhere).Where(whereMap)
	if cond != "" {
		gdb.Where(cond, condArgs...)
	}
	if condition != nil {
		gdb.Where(condition, args...)
	}
	return gdb.Updates(dataMap)
}

// 扩展方法-更新，支持自定义条件、支持所有字段置NULL、支持数字类型置零
func (db *Db) UpdatesNew(uWhere interface{}, uData basemodel.IEntidy) (tx *gorm.DB) {
	return db.UpdatesNewWithCondition(uWhere, uData, nil)
}

// 执行写操作
func ExecGTx(exec interface{}) bool {
	if exec == nil {
		return false
	}

	err := gormdb.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
				fmt.Println(err)
			}
		}()

		db := Db{tx}
		switch fn := exec.(type) {
		case func(Db):
			fn(db)
		case func(GTx):
			fn(tx)
		case func(Db) bool:
			rlt := fn(db)
			if !rlt {
				tx.Rollback()
				zap.L().Error("手动事务回滚")
				return errors.New("手动事务回滚")
			}
		case func(GTx) bool:
			rlt := fn(tx)
			if !rlt {
				tx.Rollback()
				zap.L().Error("手动事务回滚")
				return errors.New("手动事务回滚")
			}
		default:
			panic(errors.New("param is invalid fun"))
		}

		return nil
	})

	if err != nil {
		zap.L().Error("事务执行失败[-]", zap.Error(err))
		return false
	}
	return true
}

// 处理where条件, 转换合并成map条件+自定义条件
func convertWhereEntity(where interface{}) (map[string]interface{}, string, []interface{}) {
	if iface, ok := where.(basemodel.IEntidy); ok {
		columns, conditions, args := iface.GetConditions()
		whereMap, _ := jsonUtils.StructToDbMap(where)
		for k, v := range whereMap {
			// 删除旧key
			delete(whereMap, k)
			var hasCol bool
			// 判断condition column是否在where条件里，如果包含则map里去除
			for _, col := range columns {
				if col == k {
					hasCol = true
				}
			}
			if !hasCol && k != "curPage" && k != "pageSize" {
				// 添加新key
				whereMap[k] = v
			}
		}
		return whereMap, conditions, args
	} else {
		panic(errors.New("where entity is invalid"))
	}
}
