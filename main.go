package main

import (
	"fmt"

	"github.com/aditya3232/gatewatchApp-services.git/config"
	"github.com/aditya3232/gatewatchApp-services.git/connection"
	"github.com/aditya3232/gatewatchApp-services.git/helper"
	"github.com/aditya3232/gatewatchApp-services.git/model/consumer_vandal_detection"
	"github.com/gin-gonic/gin"
)

func main() {
	defer helper.RecoverPanic()

	forever := make(chan bool)
	go func() {
		consumerVandalDetectionRepository := consumer_vandal_detection.NewRepository(connection.DatabaseMysql(), connection.RabbitMQ())
		consumerVandalDetectionService := consumer_vandal_detection.NewService(consumerVandalDetectionRepository)

		_, err := consumerVandalDetectionService.ConsumerQueueVandalDetection()
		if err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println(" [*] - waiting for messages")
	<-forever

	router := gin.Default()
	if config.CONFIG.DEBUG == 0 {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Run(fmt.Sprintf("%s:%s", config.CONFIG.APP_HOST, config.CONFIG.APP_PORT))

}
