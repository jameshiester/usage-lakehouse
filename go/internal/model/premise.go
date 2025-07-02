package model

import "time"

type PremiseTypeCode string

const (
	PremiseTypeCodeResidential         PremiseTypeCode = "01"
	PremiseTypeCodeSmallNonResidential PremiseTypeCode = "02"
	PremiseTypeCodeLargeNonResidential PremiseTypeCode = "03"
)

type PremiseType struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created_dttm"`
	Updated     time.Time `json:"updated_dttm"`
}

type PowerRegionPremiseType struct {
	ID            string    `json:"id"`
	PowerRegionID string    `json:"power_region_id"`
	Code          string    `json:"power_region_premise_type_code"`
	Description   string    `json:"description"`
	PremiseTypeID string    `json:"premise_type_id"`
	Created       time.Time `json:"created_dttm"`
	Updated       time.Time `json:"updated_dttm"`
}

type Premise struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	CustomerName  string    `json:"customer_name"`
	Address       Address   `json:"address"`
	PremiseTypeID *string   `json:"premise_type_id,omitempty"`
	PowerRegionID string    `json:"power_region_id"`
	Created       time.Time `json:"created_dttm"`
	Updated       time.Time `json:"updated_dttm"`
}

type PremiseAccountJunction struct {
	ID        string     `json:"id"`
	PremiseID string     `json:"premise_id"`
	AccountID string     `json:"account_id"`
	Status    string     `json:"status"`
	MinStart  time.Time  `json:"min_start_dt"`
	MaxEnd    *time.Time `json:"max_end_dt,omitempty"`
	Created   time.Time  `json:"created_dttm"`
	Updated   time.Time  `json:"updated_dttm"`
}

type PremiseAccountHistory struct {
	ID             string     `json:"id"`
	PremiseID      string     `json:"premise_id"`
	AccountID      string     `json:"account_id"`
	EstimatedStart time.Time  `json:"estimated_start_dt"`
	EstimatedEnd   *time.Time `json:"estimated_end_dt,omitempty"`
	Start          *time.Time `json:"start_dt,omitempty"`
	End            *time.Time `json:"end_dt,omitempty"`
	Created        time.Time  `json:"created_dttm"`
	Updated        time.Time  `json:"updated_dttm"`
}
