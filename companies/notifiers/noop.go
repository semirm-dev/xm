package notifiers

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Noop struct{}

func NewNoopNotifier() *Noop {
	return &Noop{}
}

func (n *Noop) Notify(ctx context.Context, event string, message any) error {
	logrus.Infof("received event [%s] with payload: %v", event, message)
	return nil
}
