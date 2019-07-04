package bus

import (
	"github.com/sirupsen/logrus"
)

type BusController interface {
	ConsumeMessages(*logrus.Entry) error
}
