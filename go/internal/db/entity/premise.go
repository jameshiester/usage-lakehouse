package dbentity

import "time"

type PremiseType struct {
	Code        string
	Name        string
	Description string
	Created     time.Time
	Updated     time.Time
}

type PremiseAccountStatus struct {
	Code        string
	Name        string
	Description string
	Active      bool
	Created     time.Time
	Updated     time.Time
}

type Premise struct {
	ID              string
	AccountID       string
	Code            string
	Name            *string
	CustomerName    string
	AddressLine1    *string
	AddressLine2    *string
	City            *string
	State           *string
	Zip             *string
	Country         *string
	PremiseTypeCode string
	PowerRegionID   string
	Created         time.Time
	Updated         time.Time
}

type PremiseAccountJunction struct {
	AccountID                string
	PremiseID                string
	PremiseAccountStatusCode string
	MinStart                 time.Time
	MaxEnd                   *time.Time
	Created                  time.Time
	Updated                  time.Time
}

type PremiseAccountHistory struct {
	AccountID      string
	PremiseID      string
	EstimatedStart time.Time
	Start          *time.Time
	End            *time.Time
	EstimatedEnd   *time.Time
	Created        time.Time
	Updated        time.Time
}
