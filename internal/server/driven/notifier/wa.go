package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ghazlabs/idn-remote-entry/internal/shared/core"
	"github.com/ghazlabs/idn-remote-entry/internal/shared/rmqutil"
	"github.com/wagslane/go-rabbitmq"
	"gopkg.in/validator.v2"
)

type WhatsappNotifier struct {
	WhatsappNotifierConfig
}

func NewWhatsappNotifier(cfg WhatsappNotifierConfig) (*WhatsappNotifier, error) {
	// validate config
	err := validator.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &WhatsappNotifier{
		WhatsappNotifierConfig: cfg,
	}, nil
}

type WhatsappNotifierConfig struct {
	RmqPublisher         *rabbitmq.Publisher `validate:"nonnil"`
	WhatsappRecipientIDs []string            `validate:"nonzero"`
	QueueName            string              `validate:"nonzero"`
}

func (n *WhatsappNotifier) Notify(ctx context.Context, v core.VacancyRecord) error {
	for _, waID := range n.WhatsappRecipientIDs {
		ntf := core.WhatsappNotification{
			RecipientID:   waID,
			VacancyRecord: v,
		}
		data, _ := json.Marshal(ntf)
		err := rmqutil.Publish(ctx, rmqutil.PublishParams{
			Publisher:   n.RmqPublisher,
			ContentType: "application/json",
			QueueName:   n.QueueName,
			Data:        data,
		})
		if err != nil {
			return fmt.Errorf("failed to publish notification %+v: %w", ntf, err)
		}
	}

	return nil
}
