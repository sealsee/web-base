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

type Db struct {
	*gorm.DB
}

func (db *Db) UpdatesNew(uWhere interface{}, uData basemodel.IEntidy) (tx *gorm.DB) {
	dataMap, _ := jsonUtils.StructToDbMap(uData)
	// fmt.Printf("---> uData map: %v \n", dataMap)
	for _, col := range uData.GetToNullCols() {
		if _, ok := dataMap[col]; ok && dataMap[col] != "" {
			continue
		}
		dataMap[col] = nil
	}
	// fmt.Printf("---> dataMap: %v \n", dataMap)
	return db.Model(&uWhere).Where(uWhere).Updates(dataMap)
}

func InitGTx(gdb *gorm.DB) {
	gormdb = gdb
}

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
