package repository

import (
	"github.com/google/uuid"
	"github.com/xos/Projects/stk-project/backend/internal/model"
	"gorm.io/gorm"
)

type MenuRepository interface {
	FindAll() ([]model.Menu, error)
	FindByID(id uuid.UUID) (*model.Menu, error)
	FindChildren(parentID *uuid.UUID) ([]model.Menu, error)
	Create(menu *model.Menu) error
	Update(menu *model.Menu) error
	Delete(id uuid.UUID) error
	GetMaxOrderIndex(parentID *uuid.UUID) (int, error)
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) FindAll() ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Order("depth ASC, order_index ASC").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindByID(id uuid.UUID) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.Preload("Children", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_index ASC")
	}).First(&menu, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) FindChildren(parentID *uuid.UUID) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Where("parent_id = ?", parentID).Order("order_index ASC").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Menu{}).Error
}

func (r *menuRepository) GetMaxOrderIndex(parentID *uuid.UUID) (int, error) {
	var max int
	err := r.db.Model(&model.Menu{}).
		Where("parent_id = ?", parentID).
		Select("COALESCE(MAX(order_index), -1)").
		Scan(&max).Error
	return max, err
}
