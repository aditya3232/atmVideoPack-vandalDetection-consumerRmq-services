package add_vandal_detection_to_elastic

// ini entity data yg akan dikirim ke elastic

type ElasticVandalDetection struct {
	ID                             string `json:"id"`
	TidID                          *int   `json:"tid_id"`
	DateTime                       string `json:"date_time"`
	Person                         string `json:"person"`
	FileNameCaptureVandalDetection string `json:"file_name_capture_vandal_detection"`
}
