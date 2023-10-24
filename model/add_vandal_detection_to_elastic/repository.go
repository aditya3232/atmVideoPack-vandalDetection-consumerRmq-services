package add_vandal_detection_to_elastic

import (
	"strconv"
	"strings"

	esv7 "github.com/elastic/go-elasticsearch/v7"
)

type Repository interface {
	CreateElasticVandalDetection(elasticVandalDetection ElasticVandalDetection) (ElasticVandalDetection, error)
}

type repository struct {
	elasticsearch *esv7.Client
}

func NewRepository(elasticsearch *esv7.Client) *repository {
	return &repository{elasticsearch}
}

func (r *repository) CreateElasticVandalDetection(elasticVandalDetection ElasticVandalDetection) (ElasticVandalDetection, error) {

	// Menggunakan library "github.com/elastic/go-elasticsearch" untuk melakukan operasi penyimpanan
	// Gantilah `indexName` dengan nama index Elasticsearch yang sesuai
	indexName := "vandal_detection_index"

	// Anda dapat membuat body dokumen yang akan disimpan di Elasticsearch
	// Misalnya, jika Anda ingin menyimpan data deteksi manusia yang diberikan sebagai JSON:
	body := []byte(`{
		"id": "` + elasticVandalDetection.ID + `",
		"tid_id": "` + strconv.Itoa(*elasticVandalDetection.TidID) + `",
		"date_time": "` + elasticVandalDetection.DateTime + `",
		"person": "` + elasticVandalDetection.Person + `",
		"file_name_capture_vandal_detection": "` + elasticVandalDetection.FileNameCaptureVandalDetection + `"
	}`)

	// Mengirimkan data ke Elasticsearch untuk disimpan
	_, err := r.elasticsearch.Index(indexName, strings.NewReader(string(body)))
	if err != nil {
		return elasticVandalDetection, err
	}

	// Jika operasi berhasil, Anda dapat mengembalikan data yang sama yang Anda terima sebagai argumen.
	return elasticVandalDetection, nil

}
