package notifiers

import "context"

type Kafka struct{}

func NewKafkaNotifier() *Kafka {
	return &Kafka{}
}

func (n *Kafka) Notify(ctx context.Context, event string, message any) error {
	//TODO: finish implementation with Kafka, send data to Kafka topic
	return nil
}
