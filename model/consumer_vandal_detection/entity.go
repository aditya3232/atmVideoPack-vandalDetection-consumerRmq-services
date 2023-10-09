package consumer_vandal_detection

// json di struct ini disesuaikan dengan key payload rmq
type VandalDetection struct {
	Tid                            string `json:"Tid"`
	DateTime                       string `json:"DateTime"`
	Person                         string `json:"Person"`
	FileNameCaptureVandalDetection string `json:"ConvertedFile"`
}

// table name
func (m *VandalDetection) TableName() string {
	return "tb_vandal_detection"
}
