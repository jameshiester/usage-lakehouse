package model

import "time"

type UsageData struct {
	AssetID  string  `json:"asset_id" parquet:"name=asset_id, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	UsageQty float64 `json:"usage_qty" parquet:"name=usage_qty, type=DOUBLE"`
}

type TransactionSetPurposeCode string

type ReportTypeCode string

type ProductTransferDetailTypeCode string

const (
	TransactionSetPurposeCodeOriginal                                TransactionSetPurposeCode     = "00"
	TransactionSetPurposeCodeCanceled                                TransactionSetPurposeCode     = "01"
	TransactionSetPurposeCodeReplace                                 TransactionSetPurposeCode     = "05"
	TransactionSetPurposeCodePolledServices                          TransactionSetPurposeCode     = "EX"
	ReportTypeCodeIDR                                                ReportTypeCode                = "C1"
	ReportTypeCodeNIDR                                               ReportTypeCode                = "DD"
	ProductTransferDetailTypeCodeNonIntervalDetail                   ProductTransferDetailTypeCode = "PL"
	ProductTransferDetailTypeCodeNonIntervalUsage                    ProductTransferDetailTypeCode = "SU"
	ProductTransferDetailTypeCodeUnmeteredServices                   ProductTransferDetailTypeCode = "BD"
	ProductTransferDetailTypeCodeIntervalSummary                     ProductTransferDetailTypeCode = "BO"
	ProductTransferDetailTypeCodeNetIntervalUsageSummary             ProductTransferDetailTypeCode = "IA"
	ProductTransferDetailTypeCodeIntervalDetail                      ProductTransferDetailTypeCode = "PM"
	ProductTransferDetailTypeCodeNetIntervalUsageSummaryAcrossMeters ProductTransferDetailTypeCode = "PP"
)

type ErcotQuantityDelivered struct {
	Quantity    float64   `json:"quantity"`
	IntervalEnd time.Time `json:"interval_end"`
}

type ErcotProductTransferDetail struct {
	TransferType       ProductTransferDetailTypeCode `json:"product_transfer_detail_type_code"   validate:"required"`
	ServicePeriodStart *time.Time                    `json:"service_period_start,omitempty"`
	ServicePeriodEnd   *time.Time                    `json:"service_period_end,omitempty"`
	ExchangeDate       *time.Time                    `json:"exchange_date,omitempty"`
	MeterRole          *string                       `json:"meter_role,omitempty"`
	MeterType          *string                       `json:"meter_type,omitempty"`
	Channel            *string                       `json:"channel,omitempty"`
	MeterName          *string                       `json:"meter_name,omitempty"`
	Quantities         *[]ErcotQuantityDelivered     `json:"quantity_delivered,omitempty"`
}

type ErcotMonthlyUsageTransaction struct {
	TransactionID          string                       `json:"transaction_id"  validate:"required"`
	Purpose                TransactionSetPurposeCode    `json:"transaction_set_purpose_code"   validate:"required"`
	Date                   time.Time                    `json:"date"  validate:"required"`
	PremiseCode            string                       `json:"premise_code"  validate:"required"`
	PowerRegion            string                       `json:"power_region"  validate:"required"`
	ReportType             ReportTypeCode               `json:"report_type_code"  validate:"required"`
	Final                  *string                      `json:"action_code,omitempty"`
	TdspName               string                       `json:"tdsp_name"  validate:"required"`
	TdspLegalID            string                       `json:"tdsp_legal_id"  validate:"required"`
	CrName                 string                       `json:"cr_name"  validate:"required"`
	CrLegalID              string                       `json:"cr_legal_id"  validate:"required"`
	ProductTransferDetails []ErcotProductTransferDetail `json:"product_transfer_details"  validate:"required"`
}

type UsageTransaction struct {
	TransactionID string                    `json:"transaction_id"`
	Purpose       TransactionSetPurposeCode `json:"transaction_set_purpose_code"`
	Date          time.Time                 `json:"date"`
	PowerRegion   string                    `json:"power_region"`
	ReportType    ReportTypeCode            `json:"report_type_code"`
	Final         *string                   `json:"action_code,omitempty"`
	TdspName      string                    `json:"tdsp_name"`
	TdspLegalID   string                    `json:"tdsp_legal_id"`
	CrName        string                    `json:"cr_name"`
	CrLegalID     string                    `json:"cr_legal_id"`
}

type UsageTransactionDetail struct {
	TransactionID   string
	Purpose         TransactionSetPurposeCode
	TransactionDate time.Time
	Consumption     float64
	Production      float64
	MeterCode       string
	IntervalStart   time.Time
	IntervalEnd     time.Time
	PowerRegion     string
	TdspName        string
	TdspLegalID     string
	Final           bool
	PremiseCode     string
}
