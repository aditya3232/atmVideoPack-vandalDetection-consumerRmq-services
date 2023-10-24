package tb_vandal_detection

import (
	"strconv"
	"time"
)

// entity TbVandalDetection
type TbVandalDetection struct {
	ID                             int       `gorm:"primaryKey" json:"id"`
	CreatedAt                      time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt                      time.Time `gorm:"column:updated_at;default:now()" json:"updated_at"`
	TidID                          *int      `json:"tid_id"`
	DateTime                       string    `json:"date_time"`
	Person                         string    `json:"Person"`
	FileNameCaptureVandalDetection string    `json:"file_name_capture_vandal_detection"`
}

func (m *TbVandalDetection) TableName() string {
	return "tb_vandal_detection"
}

func (e *TbVandalDetection) RedisKey() string {
	if e.ID == 0 {
		return "tb_vandal_detection"
	}

	return "tb_vandal_detection:" + strconv.Itoa(e.ID)
}
