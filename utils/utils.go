package utils

import (
	"time"
)

type TransactionFile struct {
	Name             string
	Path             string
	LastModifiedTime time.Time // modification time
}