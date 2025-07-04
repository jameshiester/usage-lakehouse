package main

import (
	"fmt"
	"time"
	dbentity "usage-lakehouse/internal/db/entity"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-pg/migrations/v8"
	"github.com/google/uuid"
)

func init() {
	gofakeit.Seed(0)
	premiseAccountStatuses := []dbentity.PremiseAccountStatus{
		{
			Code:   "ACTIVE",
			Name:   "Active",
			Active: true,
		},
		{
			Code:   "ENROLL_REQUESTED",
			Name:   "Enrollment Requested",
			Active: false,
		},
		{
			Code:   "ENROLL_REJECTED",
			Name:   "Enrollment Requested",
			Active: false,
		},
		{
			Code:   "PENDING_ENROLLMENT",
			Name:   "Pending Enrollment",
			Active: false,
		},
		{
			Code:   "DELETE_REQUESTED",
			Name:   "Delete Requested",
			Active: true,
		},
		{
			Code:   "PENDING DELETE",
			Name:   "Pending Delete",
			Active: true,
		},
		{
			Code:   "DELETED",
			Name:   "Deleted",
			Active: false,
		},
	}
	ercotID := uuid.New().String()
	pjmID := uuid.New().String()
	misoID := uuid.New().String()
	caisoID := uuid.New().String()
	nyisoID := uuid.New().String()
	isoneID := uuid.New().String()
	sppID := uuid.New().String()
	powerRegions := []dbentity.PowerRegion{
		{
			ID:   ercotID,
			Name: "ERCOT",
		},
		{
			ID:   pjmID,
			Name: "PJM",
		},
		{
			ID:   misoID,
			Name: "MISO",
		},
		{
			ID:   caisoID,
			Name: "CAISO",
		},
		{
			ID:   nyisoID,
			Name: "NYISO",
		},
		{
			ID:   isoneID,
			Name: "ISONE",
		},
		{
			ID:   sppID,
			Name: "SPP",
		},
	}
	// TDSP IDs
	oncorID := uuid.New().String()
	aepNorthID := uuid.New().String()
	aepCentralID := uuid.New().String()
	centerpointID := uuid.New().String()
	tnmpID := uuid.New().String()

	tdsps := []dbentity.TDSP{
		// ERCOT TDSPs
		{
			ID:                              oncorID,
			Name:                            "Oncor Electric Delivery Company LLC",
			Code:                            "ONCOR",
			LegalID:                         "1039940674000",
			PremiseCodeValidationExpression: "^(1044372|1017699)\\d{10}$",
		},
		{
			ID:                              aepCentralID,
			Name:                            "AEP Texas Central Company",
			Code:                            "AEP-C",
			LegalID:                         "007924772",
			PremiseCodeValidationExpression: "^1000288\\d{15}$",
		},
		{
			ID:                              aepNorthID,
			Name:                            "AEP Texas North Company",
			Code:                            "AEP-N",
			LegalID:                         "007923311",
			PremiseCodeValidationExpression: "^1000078\\d{15}$",
		},
		{
			ID:                              centerpointID,
			Name:                            "CenterPoint Energy Houston Electric LLC",
			Code:                            "CENTERPOINT",
			LegalID:                         "957877905",
			PremiseCodeValidationExpression: "^10089\\d{17}$",
		},
		{
			ID:                              tnmpID,
			Name:                            "Texas-New Mexico Power Company",
			Code:                            "TNMP",
			LegalID:                         "007929441",
			PremiseCodeValidationExpression: "^10168\\d{17}$",
		},
	}

	tdspPowerRegions := []dbentity.TDSPPowerRegion{
		// ERCOT TDSP-Power Region mappings
		{TDSPID: oncorID, PowerRegionID: ercotID},
		{TDSPID: aepNorthID, PowerRegionID: ercotID},
		{TDSPID: aepCentralID, PowerRegionID: ercotID},
		{TDSPID: centerpointID, PowerRegionID: ercotID},
		{TDSPID: tnmpID, PowerRegionID: ercotID},
	}

	usageTransactionCode := "USAGE"
	enrollmentTransactionCode := "ENROLLMENT"
	transactionTypes := []dbentity.TransactionType{
		{
			Code:        usageTransactionCode,
			Name:        "Usage",
			Description: "Historic or monthly usage transaction",
		},
		{
			Code:        enrollmentTransactionCode,
			Name:        "Enrollment",
			Description: "Enrollment requests and responses",
		},
	}

	powerRegionTransactionTypes := []dbentity.PowerRegionTransactionType{
		{
			Code:                "867",
			Name:                "Usage",
			PowerRegionID:       ercotID,
			TransactionTypeCode: usageTransactionCode,
			Description:         "Historic or monthly usage transaction",
		},
	}

	historicUsageCode := "UH"
	monthlyUsageCode := "UM"
	transactionSubTypes := []dbentity.TransactionSubType{
		{
			Code:                historicUsageCode,
			Name:                "Historic Usage",
			TransactionTypeCode: usageTransactionCode,
			Description:         "Historic usage transaction",
		},
		{
			Code:                monthlyUsageCode,
			Name:                "Monthly Usage",
			TransactionTypeCode: usageTransactionCode,
			Description:         "Monthly usage transaction",
		},
	}

	powerRegionTransactionSubTypes := []dbentity.PowerRegionTransactionSubType{
		{
			Code:                   "2",
			Name:                   "Historic Usage",
			TransactionSubTypeCode: historicUsageCode,
			PowerRegionID:          ercotID,
			Description:            "Historic Usage transmitted from ERCOT to CR",
		},
		{
			Code:                   "3",
			Name:                   "Monthly Usage",
			PowerRegionID:          ercotID,
			TransactionSubTypeCode: monthlyUsageCode,
			Description:            "Monthly or Final Usage transmitted from ERCOT to CR",
		},
	}
	usageTransactionPurposes := []dbentity.UsageTransactionPurpose{
		{
			Code:        "N",
			Name:        "New",
			IsCancel:    false,
			Description: "New usage transaction",
		},
		{
			Code:        "C",
			Name:        "Cancel",
			IsCancel:    true,
			Description: "Canceled usage transaction",
		},
		{
			Code:        "R",
			Name:        "Replace",
			IsCancel:    true,
			Description: "Used when the TDSP cancels and sends a replacement transaction for corrected data",
		},
	}

	powerRegionUsageTransactionPurposes := []dbentity.PowerRegionUsageTransactionPurpose{
		{
			Code:                        "00",
			Name:                        "Original",
			Description:                 "Conveys original readings for the account being reported.",
			PowerRegionID:               ercotID,
			UsageTransactionPurposeCode: "N",
		},
		{
			Code:                        "01",
			Name:                        "Cancellation",
			PowerRegionID:               ercotID,
			UsageTransactionPurposeCode: "C",
			Description:                 "Readings previously reported for the account are to be ignored.  This would cancel the entire period of usage for the period.",
		},
		{
			Code:                        "02",
			Name:                        "Replace",
			PowerRegionID:               ercotID,
			UsageTransactionPurposeCode: "R",
			Description:                 "Used when the TDSP cancels and sends a replacement transaction for corrected data.",
		},
	}
	premiseTypeCodes := []dbentity.PremiseType{
		{Code: "RESIDENTIAL", Name: "RESIDENTIAL", Description: "Residential premise"},
		{Code: "COMMERCIAL", Name: "COMMERCIAL", Description: "Commercial premise"},
	}

	faberID := uuid.New().String()
	acmeID := uuid.New().String()
	accounts := []dbentity.Account{
		{
			ID:      acmeID,
			LegalID: "1234567",
			Name:    "Acme Corporation",
		},
		{
			ID:      faberID,
			LegalID: "7654321",
			Name:    "Faber LLC",
		},
	}

	premise1ID := uuid.New().String()
	name := "10443720008808467"
	address := "04914 BAYONNE DR"
	city := "ROWLETT"
	state := "TX"
	zip := "750881851"
	country := "USA"

	premises := []dbentity.Premise{
		{
			ID:              premise1ID,
			PowerRegionID:   ercotID,
			PremiseTypeCode: "RESIDENTIAL",
			Name:            &name,
			Code:            "10443720008808467",
			CustomerName:    gofakeit.Name(),
			AddressLine1:    &address,
			City:            &city,
			State:           &state,
			Zip:             &zip,
			Country:         &country,
		},
	}

	meter1ID := uuid.New().String()
	meters := []dbentity.Meter{
		{ID: meter1ID,
			PremiseID:     premise1ID,
			PowerRegionID: ercotID,
			Name:          "LG12345",
			Type:          "IDR",
			LoadProfile:   "",
			CycleCode:     "13",
		},
	}

	t := time.Now().UTC()
	startTime := time.Date(t.Year()-4, t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	premiseAccountJunctions := []dbentity.PremiseAccountJunction{
		{
			AccountID:                acmeID,
			PremiseID:                premise1ID,
			PremiseAccountStatusCode: "ACTIVE",
			MinStart:                 startTime,
		},
	}
	premiseAccountHistory := []dbentity.PremiseAccountHistory{
		{
			AccountID:      acmeID,
			PremiseID:      premise1ID,
			EstimatedStart: startTime,
			Start:          &startTime,
		},
	}

	powerRegionUsageTransactionProductTransferDetailTypes := []dbentity.PowerRegionUsageTransactionProductTransferDetailType{
		{
			Code:          "PL",
			PowerRegionID: ercotID,
			Interval:      false,
			Meter:         true,
			Summary:       false,
			Name:          "Non-Interval Detail",
			Description:   "Non-Interval Detail",
		},
		{
			Code:          "SU",
			PowerRegionID: ercotID,
			Interval:      false,
			Meter:         true,
			Summary:       true,
			Name:          "Non-Interval Usage Summary",
			Description:   "Non-Interval Usage Summary",
		},
		{
			Code:          "BD",
			PowerRegionID: ercotID,
			Interval:      false,
			Meter:         false,
			Summary:       false,
			Name:          "Unmetered Services Detail",
			Description:   "Unmetered Services Detail",
		},
		{
			Code:          "BO",
			PowerRegionID: ercotID,
			Interval:      true,
			Meter:         true,
			Summary:       true,
			Name:          "Interval Summary",
			Description:   "Interval Summary",
		},
		{
			Code:          "IA",
			PowerRegionID: ercotID,
			Interval:      true,
			Meter:         true,
			Summary:       true,
			Name:          "Net Interval Usage Summary",
			Description:   "Net Interval Usage Summary",
		},
		{
			Code:          "PM",
			PowerRegionID: ercotID,
			Interval:      true,
			Meter:         true,
			Summary:       false,
			Name:          "Interval Detail",
			Description:   "Interval Detail",
		},
		{
			Code:          "PP",
			PowerRegionID: ercotID,
			Interval:      true,
			Meter:         true,
			Summary:       true,
			Name:          "Net Interval Usage Summary Across Meters",
			Description:   "Net Interval Usage Summary Across Meters",
		},
	}

	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("seeding premise_account_status table...")

		for _, status := range premiseAccountStatuses {
			_, err := db.Exec(`
			INSERT INTO public.premise_account_status (code, name, is_active)
			VALUES (?, ?, ?)
		`, status.Code, status.Name, status.Active)
			if err != nil {
				return fmt.Errorf("error inserting premise_account_status %s: %v", status.Code, err)
			}
		}
		fmt.Printf("inserted %d premise_account_status records\n", len(premiseAccountStatuses))

		fmt.Println("seeding power_region table...")
		for _, region := range powerRegions {
			_, err := db.Exec(`
				INSERT INTO public.power_region (id, name)
				VALUES (?, ?)
			`, region.ID, region.Name)
			if err != nil {
				return fmt.Errorf("error inserting power_region %s: %v", region.Name, err)
			}
		}
		fmt.Printf("inserted %d power_region records\n", len(powerRegions))

		fmt.Println("seeding tdsp table...")
		for _, tdsp := range tdsps {
			_, err := db.Exec(`
				INSERT INTO public.tdsp (id, legal_entity_name, name, legal_id, abbreviation, premise_code_validation_expression)
				VALUES (?, ?, ?, ?, ?, ?)
			`, tdsp.ID, tdsp.Name, tdsp.Name, tdsp.LegalID, tdsp.Code, tdsp.PremiseCodeValidationExpression)
			if err != nil {
				return fmt.Errorf("error inserting tdsp %s: %v", tdsp.Name, err)
			}
		}
		fmt.Printf("inserted %d tdsp records\n", len(tdsps))

		fmt.Println("seeding tdsp_power_region_junction table...")
		for _, tdspPowerRegion := range tdspPowerRegions {
			_, err := db.Exec(`
				INSERT INTO public.tdsp_power_region_junction (tdsp_id, power_region_id)
				VALUES (?, ?)
			`, tdspPowerRegion.TDSPID, tdspPowerRegion.PowerRegionID)
			if err != nil {
				return fmt.Errorf("error inserting tdsp_power_region_junction %s: %v", tdspPowerRegion.TDSPID, err)
			}
		}
		fmt.Printf("inserted %d tdsp_power_region_junction records\n", len(tdspPowerRegions))

		fmt.Println("seeding transaction_type table...")
		for _, transactionType := range transactionTypes {
			_, err := db.Exec(`
				INSERT INTO public.transaction_type (code, name, description)
				VALUES (?, ?, ?)
			`, transactionType.Code, transactionType.Name, transactionType.Description)
			if err != nil {
				return fmt.Errorf("error inserting transaction_type %s: %v", transactionType.Code, err)
			}
		}
		fmt.Printf("inserted %d transaction_type records\n", len(transactionTypes))

		fmt.Println("seeding power_region_transaction_type table...")
		for _, powerRegionTransactionType := range powerRegionTransactionTypes {
			_, err := db.Exec(`
				INSERT INTO public.power_region_transaction_type (code, name, power_region_id, transaction_type_code, description)
				VALUES (?, ?, ?, ?, ?)
			`, powerRegionTransactionType.Code, powerRegionTransactionType.Name, powerRegionTransactionType.PowerRegionID, powerRegionTransactionType.TransactionTypeCode, powerRegionTransactionType.Description)
			if err != nil {
				return fmt.Errorf("error inserting power_region_transaction_type %s: %v", powerRegionTransactionType.Code, err)
			}
		}
		fmt.Printf("inserted %d power_region_transaction_type records\n", len(powerRegionTransactionTypes))

		fmt.Println("seeding transaction_sub_type table...")
		for _, transactionSubType := range transactionSubTypes {
			_, err := db.Exec(`
				INSERT INTO public.transaction_sub_type (code, name, transaction_type_code, description)
				VALUES (?, ?, ?, ?)
			`, transactionSubType.Code, transactionSubType.Name, transactionSubType.TransactionTypeCode, transactionSubType.Description)
			if err != nil {
				return fmt.Errorf("error inserting transaction_sub_type %s: %v", transactionSubType.Code, err)
			}
		}
		fmt.Printf("inserted %d transaction_sub_type records\n", len(transactionSubTypes))
		fmt.Println("seeding power_region_transaction_sub_type table...")
		for _, powerRegionTransactionSubType := range powerRegionTransactionSubTypes {
			_, err := db.Exec(`
				INSERT INTO public.power_region_transaction_sub_type (code, name, power_region_id, transaction_sub_type_code, description)
				VALUES (?, ?, ?, ?, ?)
			`, powerRegionTransactionSubType.Code, powerRegionTransactionSubType.Name, powerRegionTransactionSubType.PowerRegionID, powerRegionTransactionSubType.TransactionSubTypeCode, powerRegionTransactionSubType.Description)
			if err != nil {
				return fmt.Errorf("error inserting power_region_transaction_sub_type %s: %v", powerRegionTransactionSubType.Code, err)
			}
		}
		fmt.Printf("inserted %d power_region_transaction_sub_type records\n", len(powerRegionTransactionSubTypes))
		fmt.Println("seeding usage_transaction_purpose table...")
		for _, usageTransactionPurpose := range usageTransactionPurposes {
			_, err := db.Exec(`
				INSERT INTO public.usage_transaction_purpose (code, name, is_cancel, description)
				VALUES (?, ?, ?, ?)
			`, usageTransactionPurpose.Code, usageTransactionPurpose.Name, usageTransactionPurpose.IsCancel, usageTransactionPurpose.Description)
			if err != nil {
				return fmt.Errorf("error inserting usage_transaction_purpose %s: %v", usageTransactionPurpose.Code, err)
			}
		}
		fmt.Printf("inserted %d usage_transaction_purpose records\n", len(usageTransactionPurposes))
		fmt.Println("seeding power_region_usage_transaction_purpose table...")
		for _, powerRegionUsageTransactionPurpose := range powerRegionUsageTransactionPurposes {
			_, err := db.Exec(`
				INSERT INTO public.power_region_usage_transaction_purpose (code, name, power_region_id, usage_transaction_purpose_code, description)
				VALUES (?, ?, ?, ?, ?)
			`, powerRegionUsageTransactionPurpose.Code, powerRegionUsageTransactionPurpose.Name, powerRegionUsageTransactionPurpose.PowerRegionID, powerRegionUsageTransactionPurpose.UsageTransactionPurposeCode, powerRegionUsageTransactionPurpose.Description)
			if err != nil {
				return fmt.Errorf("error inserting power_region_usage_transaction_purpose %s: %v", powerRegionUsageTransactionPurpose.Code, err)
			}
		}
		fmt.Printf("inserted %d power_region_usage_transaction_purpose records\n", len(powerRegionUsageTransactionPurposes))

		fmt.Println("seeding premise_type table...")
		for _, premiseType := range premiseTypeCodes {
			_, err := db.Exec(`
				INSERT INTO public.premise_type (code, name, description)
				VALUES (?, ?, ?)
			`, premiseType.Code, premiseType.Name, premiseType.Description)
			if err != nil {
				return fmt.Errorf("error inserting premise_type %s: %v", premiseType.Code, err)
			}
		}
		fmt.Printf("inserted %d premise_type records\n", len(premiseTypeCodes))

		fmt.Println("seeding account table...")
		for _, account := range accounts {
			_, err := db.Exec(`
				INSERT INTO public.account (id, legal_id, name)
				VALUES (?, ?, ?)
			`, account.ID, account.LegalID, account.Name)
			if err != nil {
				return fmt.Errorf("error inserting account %s: %v", account.ID, err)
			}
		}
		fmt.Printf("inserted %d account records\n", len(accounts))

		fmt.Println("seeding premise table...")
		for _, premise := range premises {
			_, err := db.Exec(`
				INSERT INTO public.premise (id, power_region_id, premise_type_code, name, code, customer_name, address_line_1, city, state, zip, country)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, premise.ID, premise.PowerRegionID, premise.PremiseTypeCode, premise.Name, premise.Code, premise.CustomerName, premise.AddressLine1, premise.City, premise.State, premise.Zip, premise.Country)
			if err != nil {
				return fmt.Errorf("error inserting premise %s: %v", premise.ID, err)
			}
		}
		fmt.Printf("inserted %d premise records\n", len(premises))

		fmt.Println("seeding meter table...")
		for _, meter := range meters {
			_, err := db.Exec(`
				INSERT INTO public.meter (id, premise_id, power_region_id, name, type, load_profile, cycle_code)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`, meter.ID, meter.PremiseID, meter.PowerRegionID, meter.Name, meter.Type, meter.LoadProfile, meter.CycleCode)
			if err != nil {
				return fmt.Errorf("error inserting meter %s: %v", meter.ID, err)
			}
		}
		fmt.Printf("inserted %d meter records\n", len(meters))

		fmt.Println("seeding premise_account_junction table...")
		for _, premiseAccountJunction := range premiseAccountJunctions {
			_, err := db.Exec(`
				INSERT INTO public.premise_account_junction (account_id, premise_id, premise_account_status_code, min_start_dt)
				VALUES (?, ?, ?, ?)
			`, premiseAccountJunction.AccountID, premiseAccountJunction.PremiseID, premiseAccountJunction.PremiseAccountStatusCode, premiseAccountJunction.MinStart)
			if err != nil {
				return fmt.Errorf("error inserting premise_account_junction %s: %v", premiseAccountJunction.AccountID, err)
			}
		}
		fmt.Printf("inserted %d premise_account_junction records\n", len(premiseAccountJunctions))

		fmt.Println("seeding premise_account_history table...")
		for _, premiseAccountHistory := range premiseAccountHistory {
			_, err := db.Exec(`
				INSERT INTO public.premise_account_history (account_id, premise_id, estimated_start_dt, start_dt)
				VALUES (?, ?, ?, ?)
			`, premiseAccountHistory.AccountID, premiseAccountHistory.PremiseID, premiseAccountHistory.EstimatedStart, premiseAccountHistory.Start)
			if err != nil {
				return fmt.Errorf("error inserting premise_account_history %s: %v", premiseAccountHistory.AccountID, err)
			}
		}
		fmt.Printf("inserted %d premise_account_history records\n", len(premiseAccountHistory))

		fmt.Println("seeding power_region_usage_transaction_product_transfer_detail_type table...")
		for _, powerRegionUsageTransactionProductTransferDetailType := range powerRegionUsageTransactionProductTransferDetailTypes {
			_, err := db.Exec(`
				INSERT INTO public.power_region_usage_transaction_product_transfer_detail_type (code, name, power_region_id, is_interval, is_meter, is_summary, description)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`, powerRegionUsageTransactionProductTransferDetailType.Code, powerRegionUsageTransactionProductTransferDetailType.Name, powerRegionUsageTransactionProductTransferDetailType.PowerRegionID, powerRegionUsageTransactionProductTransferDetailType.Interval, powerRegionUsageTransactionProductTransferDetailType.Meter, powerRegionUsageTransactionProductTransferDetailType.Summary, powerRegionUsageTransactionProductTransferDetailType.Description)
			if err != nil {
				return fmt.Errorf("error inserting power_region_usage_transaction_product_transfer_detail_type %s: %v", powerRegionUsageTransactionProductTransferDetailType.Code, err)
			}
		}
		fmt.Printf("inserted %d power_region_usage_transaction_product_transfer_detail_type records\n", len(powerRegionUsageTransactionProductTransferDetailTypes))

		return nil
	}, func(db migrations.DB) error {
		_, err := db.Exec(`DELETE * FROM public.premise_account_status`)
		return err
	})
}
