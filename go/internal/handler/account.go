package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"usage-lakehouse/internal/model"
	"usage-lakehouse/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type AccountHandler struct {
	repo     repository.AccountRepository
	validate *validator.Validate
}

func uniqueAccountNameValidator(repo repository.AccountRepository) validator.Func {
	return func(fl validator.FieldLevel) bool {
		ctx := context.Background()
		name := fl.Parent().FieldByName("Name").String()
		accountID := fl.Parent().FieldByName("AccountID").Interface().(*string)
		exists, err := repo.ExistsByName(ctx, name, accountID)
		return err == nil && !exists
	}
}

func uniqueAccountLegalIDValidator(repo repository.AccountRepository) validator.Func {
	return func(fl validator.FieldLevel) bool {
		ctx := context.Background()
		legalID := fl.Parent().FieldByName("LegalID").String()
		accountID := fl.Parent().FieldByName("AccountID").Interface().(*string)
		exists, err := repo.ExistsByLegalID(ctx, legalID, accountID)
		return err == nil && !exists
	}
}

func NewAccountHandler(repo repository.AccountRepository) *AccountHandler {
	validate := validator.New()
	validate.RegisterValidation("unique_account_name", uniqueAccountNameValidator(repo))
	validate.RegisterValidation("unique_account_legal_id", uniqueAccountLegalIDValidator(repo))
	return &AccountHandler{repo: repo, validate: validate}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	type accountInput struct {
		LegalID *string `json:"legal_id" validate:"required,unique_account_legal_id"`
		Name    string  `json:"name" validate:"required,unique_account_name"`
	}
	var input accountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(input); err != nil {
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
	a := model.Account{
		LegalID: input.LegalID,
		Name:    input.Name,
	}
	if err := h.repo.Create(r.Context(), &a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(a)
}

func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var a model.Account
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	a.ID = id
	if err := h.repo.Update(r.Context(), &a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a)
}

func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AccountHandler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(accounts)
}
