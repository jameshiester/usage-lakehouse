package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	dbentity "usage-lakehouse/internal/db/entity"
	"usage-lakehouse/internal/model"
	"usage-lakehouse/internal/repository"

	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type EDIMonthlyUsageHandler struct {
	repo                                                     repository.AccountRepository
	powerRegionRepo                                          repository.PowerRegionRepository
	tdspRepo                                                 repository.TDSPRepository
	premiseRepo                                              repository.PremiseRepository
	meterRepo                                                repository.MeterRepository
	usageTransactionPurposeRepo                              repository.UsageTransactionPurposeRepository
	transactionTypeRepo                                      repository.TransactionTypeRepository
	transactionSubTypeRepo                                   repository.TransactionSubTypeRepository
	usageTransactionRepo                                     repository.UsageTransactionRepository
	powerRegionUsageTransactionProductTransferDetailTypeRepo repository.PowerRegionUsageTransactionProductTransferDetailTypeRepository
}

func NewEDIMonthlyUsageHandler(repo repository.AccountRepository, powerRegionRepo repository.PowerRegionRepository, tdspRepo repository.TDSPRepository, premiseRepo repository.PremiseRepository, meterRepo repository.MeterRepository, usageTransactionPurposeRepo repository.UsageTransactionPurposeRepository, transactionTypeRepo repository.TransactionTypeRepository, transactionSubTypeRepo repository.TransactionSubTypeRepository, powerRegionUsageTransactionProductTransferDetailTypeRepo repository.PowerRegionUsageTransactionProductTransferDetailTypeRepository, usageTransactionRepo repository.UsageTransactionRepository) *EDIMonthlyUsageHandler {
	return &EDIMonthlyUsageHandler{repo: repo, powerRegionRepo: powerRegionRepo, tdspRepo: tdspRepo, premiseRepo: premiseRepo, meterRepo: meterRepo, usageTransactionPurposeRepo: usageTransactionPurposeRepo, transactionTypeRepo: transactionTypeRepo, transactionSubTypeRepo: transactionSubTypeRepo, powerRegionUsageTransactionProductTransferDetailTypeRepo: powerRegionUsageTransactionProductTransferDetailTypeRepo, usageTransactionRepo: usageTransactionRepo}
}

func (h *EDIMonthlyUsageHandler) CreateEDIMonthlyUsage(w http.ResponseWriter, r *http.Request) {
	powerRegions, err := h.powerRegionRepo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	powerRegionMap := make(map[string]dbentity.PowerRegion)
	for _, region := range powerRegions {
		powerRegionMap[region.Name] = region
	}
	validate := validator.New()
	var input model.ErcotMonthlyUsageTransaction
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"validationErrors": validationErrors.Error(),
			})
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	powerRegion := powerRegionMap[input.PowerRegion]
	tdsp, err := h.tdspRepo.GetByName(r.Context(), input.TdspName)
	if err != nil || tdsp == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	premise, err := h.premiseRepo.GetByCode(r.Context(), input.EsiID)
	if err != nil || premise == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	purpose, err := h.usageTransactionPurposeRepo.GetByPowerRegionAndCode(r.Context(), powerRegion.ID, string(input.Purpose))
	if err != nil || purpose == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	transactionType, err := h.transactionTypeRepo.GetByPowerRegionAndCode(r.Context(), powerRegion.ID, string(input.ReportType))
	if err != nil || transactionType == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	transactionSubType, err := h.transactionSubTypeRepo.GetByPowerRegionSubTypeCode(r.Context(), powerRegion.ID, string(input.ReportType))
	if err != nil || transactionSubType == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := uuid.New().String()
	isCanceled := purpose.IsCancel
	usageTransaction := dbentity.UsageTransaction{
		ID:                 id,
		TransactionID:      input.TransactionID,
		TransactionDate:    input.Date,
		PowerRegionID:      powerRegion.ID,
		TDSPID:             tdsp.ID,
		PremiseID:          premise.ID,
		Purpose:            purpose.Code,
		IsFinal:            input.Final != nil && *input.Final == "F",
		IsCanceled:         isCanceled,
		TransactionType:    transactionType.Code,
		TransactionSubType: transactionSubType.Code,
	}
	meterNameSet := make(map[string]struct{})
	uniqueMeterNames := make([]string, 0, len(input.ProductTransferDetails))
	for _, productTransferDetail := range input.ProductTransferDetails {
		if productTransferDetail.MeterName != nil {
			name := *productTransferDetail.MeterName
			if _, exists := meterNameSet[name]; !exists {
				meterNameSet[name] = struct{}{}
				uniqueMeterNames = append(uniqueMeterNames, name)
			}
		}
	}
	meterNameIDMap, err := h.meterRepo.GetNameIDMap(r.Context(), uniqueMeterNames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	powerRegionUsageTransactionProductTransferDetailTypeMap, err := h.powerRegionUsageTransactionProductTransferDetailTypeRepo.MapByCode(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type MeterTransferTypeKey struct {
		MeterName    string
		TransferType string
	}
	type GroupedProductTransferDetails struct {
		Details []model.ErcotProductTransferDetail
		Type    dbentity.PowerRegionUsageTransactionProductTransferDetailType
	}
	grouped := make(map[MeterTransferTypeKey]GroupedProductTransferDetails)
	var usageTransactionDetails []dbentity.UsageTransactionDetail
	for _, productTransferDetail := range input.ProductTransferDetails {
		meterName := *productTransferDetail.MeterName
		transferType := string(productTransferDetail.TransferType)
		key := MeterTransferTypeKey{MeterName: meterName, TransferType: transferType}
		g := grouped[key]
		g.Details = append(g.Details, productTransferDetail)
		if t, ok := powerRegionUsageTransactionProductTransferDetailTypeMap[transferType]; ok {
			g.Type = t
		}
		grouped[key] = g
	}
	for key, details := range grouped {
		if details.Type.Interval && details.Type.Meter && !details.Type.Summary {
			var consumptionDetail *model.ErcotProductTransferDetail
			var generationDetail *model.ErcotProductTransferDetail
			for i := range details.Details {
				detail := &details.Details[i]
				if detail.Channel != nil && *detail.Channel == "1" {
					consumptionDetail = detail
				}
				if detail.Channel != nil && *detail.Channel == "4" {
					generationDetail = detail
				}
			}
			uniqueTimes := make(map[time.Time]struct {
				Consumption float64
				Generation  float64
			})
			for _, d := range []*model.ErcotProductTransferDetail{consumptionDetail, generationDetail} {
				if d == nil {
					continue
				}
				if d.Quantities != nil {
					for _, q := range *d.Quantities {
						val := uniqueTimes[q.IntervalEnd]
						if d == consumptionDetail {
							val.Consumption += q.Quantity
						}
						if d == generationDetail {
							val.Generation += q.Quantity
						}
						uniqueTimes[q.IntervalEnd] = val
					}
				}
			}
			for t, v := range uniqueTimes {
				fmt.Printf("Time: %v, Consumption: %f, Generation: %f\n", t, v.Consumption, v.Generation)
				meterUUID := meterNameIDMap[key.MeterName]
				var servicePeriodStart time.Time
				var servicePeriodEnd time.Time
				if consumptionDetail != nil {
					servicePeriodStart = *consumptionDetail.ServicePeriodStart
					servicePeriodEnd = *consumptionDetail.ServicePeriodEnd
				}
				if generationDetail != nil {
					servicePeriodStart = *generationDetail.ServicePeriodStart
					servicePeriodEnd = *generationDetail.ServicePeriodEnd
				}
				newDetail := dbentity.UsageTransactionDetail{
					UsageTransactionID: id,
					MeterID:            &meterUUID,
					MeterName:          key.MeterName,
					PowerRegionID:      powerRegion.ID,
					PremiseID:          premise.ID,
					IsCanceled:         isCanceled,
					ServicePeriodStart: servicePeriodStart,
					ServicePeriodEnd:   servicePeriodEnd,
					Consumption:        &v.Consumption,
					Production:         &v.Generation,
				}
				usageTransactionDetails = append(usageTransactionDetails, newDetail)
			}
		} else if details.Type.Meter && !details.Type.Summary {
			for _, detail := range details.Details {
				meterUUID := meterNameIDMap[key.MeterName]
				if detail.Quantities != nil && len(*detail.Quantities) > 0 {
					for _, q := range *detail.Quantities {
						newDetail := dbentity.UsageTransactionDetail{
							MeterID:            &meterUUID,
							MeterName:          key.MeterName,
							PowerRegionID:      powerRegion.ID,
							PremiseID:          premise.ID,
							IsCanceled:         isCanceled,
							ServicePeriodStart: *detail.ServicePeriodStart,
							ServicePeriodEnd:   *detail.ServicePeriodEnd,
							Start:              *detail.ServicePeriodStart,
							End:                *detail.ServicePeriodEnd,
							Consumption:        &q.Quantity,
							UsageTransactionID: id,
						}
						usageTransactionDetails = append(usageTransactionDetails, newDetail)
					}
				}
			}
		}
	}
	err = h.usageTransactionRepo.SaveWithDetails(r.Context(), &usageTransaction, usageTransactionDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usageTransaction)
}
