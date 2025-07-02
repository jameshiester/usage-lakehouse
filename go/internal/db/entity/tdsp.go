package dbentity

import "time"

type TDSP struct {
	ID                              string
	Name                            string
	Code                            string
	LegalID                         string
	PremiseCodeValidationExpression string
	Created                         time.Time
	Updated                         time.Time
}

type TDSPPowerRegion struct {
	TDSPID        string
	PowerRegionID string
	Created       time.Time
	Updated       time.Time
}
