package database

import (
	"github.com/isd-sgcu/rpkm67-store/config"
	"github.com/isd-sgcu/rpkm67-store/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func InitDatabase(conf *config.DB, isDebug bool) (db *gorm.DB, err error) {
	gormConf := &gorm.Config{TranslateError: true}

	if !isDebug {
		gormConf.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	db, err = gorm.Open(postgres.Open(conf.Url), gormConf)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Object{})
	if err != nil {
		return nil, err
	}

	return
}
