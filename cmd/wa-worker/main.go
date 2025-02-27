package main

import (
	"fmt"
	"log"

	"github.com/ghazlabs/idn-remote-entry/internal/shared/rmq"
	"github.com/ghazlabs/idn-remote-entry/internal/wa-worker/core"
	"github.com/ghazlabs/idn-remote-entry/internal/wa-worker/driven/publisher/email"
	"github.com/ghazlabs/idn-remote-entry/internal/wa-worker/driven/publisher/wa"
	"github.com/ghazlabs/idn-remote-entry/internal/wa-worker/driver/worker"
	"github.com/go-resty/resty/v2"
	"github.com/riandyrn/go-env"
)

const (
	envKeyPublisherType       = "PUBLISHER_TYPE"
	envKeySmtpHost            = "SMTP_HOST"
	envKeySmtpPort            = "SMTP_PORT"
	envKeySmtpFrom            = "SMTP_FROM"
	envKeySmtpTo              = "SMTP_TO"
	envKeyWhatsappApiUser     = "WHATSAPP_API_USER"
	envKeyWhatsappApiPass     = "WHATSAPP_API_PASS"
	envKeyWhatsappApiBaseUrl  = "WHATSAPP_API_BASE_URL"
	envKeyRabbitMQConn        = "RABBITMQ_CONN"
	envKeyRabbitMQWaQueueName = "RABBITMQ_WA_QUEUE_NAME"
)

func initPublisher() (core.Publisher, error) {
	switch env.GetString(envKeyPublisherType) {
	case "wa":
		return wa.NewWaPublisher(wa.WaPublisherConfig{
			HttpClient:   resty.New(),
			Username:     env.GetString(envKeyWhatsappApiUser),
			Password:     env.GetString(envKeyWhatsappApiPass),
			WaApiBaseUrl: env.GetString(envKeyWhatsappApiBaseUrl),
		})
	case "email":
		return email.NewEmailPublisher(email.EmailPublisherConfig{
			Host: env.GetString(envKeySmtpHost),
			Port: env.GetInt(envKeySmtpPort),
			From: env.GetString(envKeySmtpFrom),
			To:   env.GetString(envKeySmtpTo),
		})
	default:
		return nil, fmt.Errorf("invalid publisher type: %s", env.GetString(envKeyPublisherType))
	}
}

func main() {
	// initialize publisher
	pub, err := initPublisher()
	if err != nil {
		log.Fatalf("failed to initialize publisher: %v", err)
	}

	// initialize service
	svc, err := core.NewService(core.ServiceConfig{
		Publisher: pub,
	})
	if err != nil {
		log.Fatalf("failed to initialize service: %v", err)
	}

	// initialize consumer
	rmqConsumer, err := rmq.NewConsumer(rmq.ConsumerConfig{
		QueueName:          env.GetString(envKeyRabbitMQWaQueueName),
		RabbitMQConnString: env.GetString(envKeyRabbitMQConn),
	})
	if err != nil {
		log.Fatalf("failed to initialize consumer: %v", err)
	}
	defer rmqConsumer.Close()

	// initialize worker
	w, err := worker.New(worker.Config{
		Service:     svc,
		RmqConsumer: rmqConsumer,
	})
	if err != nil {
		log.Fatalf("failed to initialize worker: %v", err)
	}

	// run worker
	log.Printf("wa-worker is running...")
	err = w.Run()
	if err != nil {
		log.Fatalf("failed to run worker: %v", err)
	}
}
