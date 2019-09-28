package accounts

import (
	"time"
)

type Transaction struct {
	Date    time.Time
	Amount  float32
	Type    string
	Balance float32
}