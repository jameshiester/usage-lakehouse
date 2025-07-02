package dbentity

import "time"

type TransactionType struct {
	Code        string
	Name        string
	Description string
	Created     time.Time
	Updated     time.Time
}

type PowerRegionTransactionType struct {
	PowerRegionID       string
	TransactionTypeCode string
	Name                string
	Code                string
	Description         string
	Created             time.Time
	Updated             time.Time
}

type TransactionSubType struct {
	Code                string
	TransactionTypeCode string
	Name                string
	Description         string
	Created             time.Time
	Updated             time.Time
}

type UsageTransactionPurpose struct {
	Code        string
	Name        string
	IsCancel    bool
	Description string
	Created     time.Time
	Updated     time.Time
}

type PowerRegionTransactionSubType struct {
	Code                   string
	PowerRegionID          string
	TransactionSubTypeCode string
	Name                   string
	Description            string
	Created                time.Time
	Updated                time.Time
}

type PowerRegionUsageTransactionPurpose struct {
	Code                        string
	PowerRegionID               string
	UsageTransactionPurposeCode string
	Name                        string
	Description                 string
	Created                     time.Time
	Updated                     time.Time
}

type PowerRegionUsageTransactionProductTransferDetailType struct {
	Code          string
	PowerRegionID string
	Interval      bool
	Meter         bool
	Summary       bool
	Name          string
	Description   string
	Created       time.Time
	Updated       time.Time
}

type UsageTransaction struct {
	ID                 string
	TransactionID      string
	TransactionType    string
	TransactionSubType string
	TransactionDate    time.Time
	ServicePeriodStart time.Time
	ServicePeriodEnd   time.Time
	IsFinal            bool
	IsCanceled         bool
	Purpose            string
	PowerRegionID      string
	TDSPID             string
	PremiseID          string
	Created            time.Time
	Updated            time.Time
}

type UsageTransactionDetail struct {
	UsageTransactionID string
	Start              time.Time
	End                time.Time
	ServicePeriodStart time.Time
	ServicePeriodEnd   time.Time
	IsCanceled         bool
	MeterID            *string
	MeterName          string
	PowerRegionID      string
	PremiseID          string
	Consumption        *float64
	Production         *float64
	Created            time.Time
	Updated            time.Time
}
