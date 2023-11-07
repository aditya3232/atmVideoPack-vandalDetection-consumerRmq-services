package consumer_vandal_detection

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aditya3232/atmVideoPack-vandalDetection-consumerRmq-services.git/config"
	"github.com/aditya3232/atmVideoPack-vandalDetection-consumerRmq-services.git/helper"
	libraryMinio "github.com/aditya3232/atmVideoPack-vandalDetection-consumerRmq-services.git/library/minio"
	log_function "github.com/aditya3232/atmVideoPack-vandalDetection-consumerRmq-services.git/log"
	"github.com/aditya3232/atmVideoPack-vandalDetection-consumerRmq-services.git/model/add_vandal_detection_to_elastic"
	esv7 "github.com/elastic/go-elasticsearch/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type Repository interface {
	ConsumerQueueVandalDetection() (RmqConsumerVandalDetection, error)
}

type repository struct {
	db            *gorm.DB
	rabbitmq      *amqp.Connection
	elasticsearch *esv7.Client
}

func NewRepository(db *gorm.DB, rabbitmq *amqp.Connection, elasticsearch *esv7.Client) *repository {
	return &repository{db, rabbitmq, elasticsearch}
}

func (r *repository) ConsumerQueueVandalDetection() (RmqConsumerVandalDetection, error) {
	var rmqConsumerVandalDetection RmqConsumerVandalDetection

	// create channel
	ch, err := r.rabbitmq.Channel()
	if err != nil {
		return rmqConsumerVandalDetection, err
	}
	defer ch.Close()

	// consume queue
	msgs, err := ch.Consume(
		"VandalDetectionQueue", // name queue
		"",                     // Consumer name (empty for random name)
		true,                   // Auto-acknowledgment (set to true for auto-ack)
		false,                  // Exclusive
		false,                  // No-local
		false,                  // No-wait
		nil,                    // Arguments
	)

	if err != nil {
		return rmqConsumerVandalDetection, err
	}

	// get message
	for d := range msgs {
		newVandalDetection := rmqConsumerVandalDetection
		err := json.Unmarshal(d.Body, &newVandalDetection)
		if err != nil {
			return rmqConsumerVandalDetection, err
		}

		// konversi VandalDetection.FileNameCaptureVandalDetection string ke bytes
		bytesConvertedFile, err := base64.StdEncoding.DecodeString(newVandalDetection.ConvertedFileCaptureVandalDetection)
		if err != nil {
			return rmqConsumerVandalDetection, err
		}

		// Mengunggah gambar ke MinIO
		objectName := "vandal-detection/" + helper.DateTimeToStringWithStrip(time.Now()) + ".jpg"
		FileNameCaptureVandalDetection := helper.DateTimeToStringWithStrip(time.Now()) + ".jpg"

		key, err := libraryMinio.UploadFileFromPutObject(config.CONFIG.MINIO_BUCKET, objectName, bytesConvertedFile)
		if err != nil {
			log_function.Error(fmt.Sprintf("Gambar gagal diunggah ke MinIO dengan nama objek: %s\n", key.Key))
			return rmqConsumerVandalDetection, err
		}

		// add data newVandalDtection to elasticsearch with CreateElasticVandalDetection
		repoElastic := add_vandal_detection_to_elastic.NewRepository(r.elasticsearch)
		resultElastic, err := repoElastic.CreateElasticVandalDetection(
			add_vandal_detection_to_elastic.ElasticVandalDetection{
				ID:                             helper.DateTimeToStringWithStrip(time.Now()),
				Tid:                            newVandalDetection.Tid,
				DateTime:                       newVandalDetection.DateTime,
				Person:                         newVandalDetection.Person,
				FileNameCaptureVandalDetection: FileNameCaptureVandalDetection,
			},
		)
		if err != nil {
			return rmqConsumerVandalDetection, err
		}
		// log result elastic
		log_function.Info(fmt.Sprintf("Result elastic: %v\n", resultElastic))

		// create data tb_vandal_detection
		// repo := tb_vandal_detection.NewRepository(r.db)
		// _, err = repo.Create(
		// 	tb_vandal_detection.TbVandalDetection{
		// 		TidID:                         newVandalDetection.TidID,
		// 		DateTime:                      newVandalDetection.DateTime,
		// 		Person:                        newVandalDetection.Person,
		// 		FileNameCaptureVandalDetection: FileNameCaptureVandalDetection,
		// 	},
		// )
		// if err != nil {
		// 	return rmqConsumerVandalDetection, err
		// }

	}

	return rmqConsumerVandalDetection, nil

}
