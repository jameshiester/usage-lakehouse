package dbentity

import "time"

type Meter struct {
	ID            string
	AccountID     string
	PremiseID     string
	PowerRegionID string
	Name          string
	Type          string
	LoadProfile   string
	CycleCode     string
	Active        bool
	Created       time.Time
	Updated       time.Time
}
