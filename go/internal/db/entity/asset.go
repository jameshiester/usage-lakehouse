package dbentity

import "time"

type Asset struct {
	ID        string
	Name      string
	Code      string
	AccountID string
	MeterID   *string
	PremiseID string
	Created   time.Time
	Updated   time.Time
}
