package utils

import (
	"errors"

	"github.com/go-jet/jet/v2/qrm"
)

func ErrIsNoRows(err error) bool {
	return errors.Is(err, qrm.ErrNoRows)
}
