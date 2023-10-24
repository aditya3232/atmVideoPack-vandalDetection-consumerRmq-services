package consumer_vandal_detection

// json di struct ini disesuaikan dengan key payload rmq
type RmqConsumerVandalDetection struct {
	TidID                               *int   `json:"tid_id"`
	DateTime                            string `json:"date_time"`
	Person                              string `json:"Person"`
	ConvertedFileCaptureVandalDetection string `json:"converted_file_capture_vandal_detection"`
}
