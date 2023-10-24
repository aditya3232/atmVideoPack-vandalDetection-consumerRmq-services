package tb_vandal_detection

import "gorm.io/gorm"

type Repository interface {
	Create(tbVandalDetection TbVandalDetection) (TbVandalDetection, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(tbVandalDetection TbVandalDetection) (TbVandalDetection, error) {
	err := r.db.Create(&tbVandalDetection).Error
	if err != nil {
		return tbVandalDetection, err
	}

	return tbVandalDetection, nil
}
