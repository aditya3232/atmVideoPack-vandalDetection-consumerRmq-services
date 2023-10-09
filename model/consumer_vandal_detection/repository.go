package consumer_vandal_detection

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aditya3232/gatewatchApp-services.git/config"
	"github.com/aditya3232/gatewatchApp-services.git/helper"
	libraryMinio "github.com/aditya3232/gatewatchApp-services.git/library/minio"
	"github.com/aditya3232/gatewatchApp-services.git/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type Repository interface {
	ConsumerQueueVandalDetection() (VandalDetection, error)
	Create(vandalDetection VandalDetection) (VandalDetection, error)
}

type repository struct {
	db       *gorm.DB
	rabbitmq *amqp.Connection
}

func NewRepository(db *gorm.DB, rabbitmq *amqp.Connection) *repository {
	return &repository{db, rabbitmq}
}

func (r *repository) ConsumerQueueVandalDetection() (VandalDetection, error) {

	// create channel
	channel, err := r.rabbitmq.Channel()
	if err != nil {
		return VandalDetection{}, err
	}
	defer channel.Close()

	// consume queue
	msgs, err := channel.Consume(
		"VandalDetectionQueue", // name queue
		"",                     // Consumer name (empty for random name)
		true,                   // Auto-acknowledgment (set to true for auto-ack)
		false,                  // Exclusive
		false,                  // No-local
		false,                  // No-wait
		nil,                    // Arguments
	)

	if err != nil {
		return VandalDetection{}, err
	}

	// get message
	for d := range msgs {
		vandalDetection := VandalDetection{}
		err := json.Unmarshal(d.Body, &vandalDetection)
		if err != nil {
			return VandalDetection{}, err
		}

		// konversi vandalDetection.FileNameCaptureVandalDetection string ke bytes
		bytesConvertedFile, err := base64.StdEncoding.DecodeString(vandalDetection.FileNameCaptureVandalDetection)
		if err != nil {
			return VandalDetection{}, err
		}

		// Mengunggah gambar ke MinIO
		objectName := "vandal-detection/" + helper.DateTimeToStringWithStrip(time.Now()) + ".jpg"
		FileNameCaptureVandalDetection := helper.DateTimeToStringWithStrip(time.Now()) + ".jpg"

		key, err := libraryMinio.UploadFileFromPutObject(config.CONFIG.MINIO_BUCKET, objectName, bytesConvertedFile)
		if err != nil {
			log.Error(fmt.Sprintf("Gambar gagal diunggah ke MinIO dengan nama objek: %s\n", key.Key))
			return VandalDetection{}, err
		}

		// insert Tid, DateTime, Person, File, ConvertedFile from message to db
		_, err = r.Create(
			VandalDetection{
				Tid:                            vandalDetection.Tid,
				DateTime:                       vandalDetection.DateTime,
				Person:                         vandalDetection.Person,
				FileNameCaptureVandalDetection: FileNameCaptureVandalDetection,
			},
		)
		if err != nil {
			return VandalDetection{}, err
		}
	}

	return VandalDetection{}, nil

}

func (r *repository) Create(vandalDetection VandalDetection) (VandalDetection, error) {
	err := r.db.Create(&vandalDetection).Error
	if err != nil {
		return VandalDetection{}, err
	}

	return vandalDetection, nil
}
