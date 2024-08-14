package internal

import (
	"errors"
	"fmt"

	"github.com/milovidov983/oms-temporal/pkg/models"
)

var ErrWrongStatus = errors.New("wrong order status")

func CreateError(err error, status models.OrderStatus) error {
	return fmt.Errorf("%w: order status - %s", err, status.String())
}
