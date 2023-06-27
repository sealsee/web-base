package tx

import (
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
}

func ExecGTx(exec func(GTx)) bool {
	err := gormdb.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
				zap.L().Error("执行失败", zap.Any("-", err))
				fmt.Println(err)
				// panic(err)
			}
		}()

		exec(tx)
		return nil
	})

	if err != nil {
		zap.L().Error("事务执行失败[-]", zap.Error(err))
		return false
	}
	return true
}
