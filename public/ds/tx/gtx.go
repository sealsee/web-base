package tx

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var gormdb *gorm.DB

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

		switch fn := exec.(type) {
		case func(GTx):
			fn(tx)
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
