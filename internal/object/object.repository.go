package object

import (
	"github.com/isd-sgcu/rpkm67-store/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindByKey(key string, result *model.Object) error
	Upload(file *model.Object) error
	DeleteByKey(key string) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) DeleteByKey(key string) error {
	return r.db.Where("object_key = ?", key).Delete(&model.Object{}).Error
}

func (r *repositoryImpl) FindByKey(key string, result *model.Object) error {
	return r.db.Where("object_key = ?", key).First(result).Error
}

func (r *repositoryImpl) Upload(file *model.Object) error {
	return r.db.Create(&file).Error
}
