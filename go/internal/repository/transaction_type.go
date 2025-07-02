package repository

import (
	"context"
	dbentity "usage-lakehouse/internal/db/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionTypeRepository interface {
	Create(ctx context.Context, t *dbentity.TransactionType) error
	GetByCode(ctx context.Context, code string) (*dbentity.TransactionType, error)
	Update(ctx context.Context, t *dbentity.TransactionType) error
	Delete(ctx context.Context, code string) error
	List(ctx context.Context) ([]dbentity.TransactionType, error)
	GetByPowerRegionAndCode(ctx context.Context, powerRegionID string, powerRegionTransactionTypeCode string) (*dbentity.TransactionType, error)
}

type transactionTypeRepositorySQL struct {
	db *pgxpool.Pool
}

func NewTransactionTypeRepository(db *pgxpool.Pool) TransactionTypeRepository {
	return &transactionTypeRepositorySQL{db: db}
}

func (r *transactionTypeRepositorySQL) Create(ctx context.Context, t *dbentity.TransactionType) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO transaction_type (code, name, description) VALUES ($1, $2, $3)`,
		t.Code, t.Name, t.Description,
	)
	return err
}

func (r *transactionTypeRepositorySQL) GetByCode(ctx context.Context, code string) (*dbentity.TransactionType, error) {
	var t dbentity.TransactionType
	err := r.db.QueryRow(ctx, `SELECT code, name, description, created, updated FROM transaction_type WHERE code=$1`, code).
		Scan(&t.Code, &t.Name, &t.Description, &t.Created, &t.Updated)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *transactionTypeRepositorySQL) Update(ctx context.Context, t *dbentity.TransactionType) error {
	_, err := r.db.Exec(ctx, `UPDATE transaction_type SET name=$1, description=$2 WHERE code=$3`, t.Name, t.Description, t.Code)
	return err
}

func (r *transactionTypeRepositorySQL) Delete(ctx context.Context, code string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM transaction_type WHERE code=$1`, code)
	return err
}

func (r *transactionTypeRepositorySQL) List(ctx context.Context) ([]dbentity.TransactionType, error) {
	rows, err := r.db.Query(ctx, `SELECT code, name, description, created, updated FROM transaction_type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var types []dbentity.TransactionType
	for rows.Next() {
		var t dbentity.TransactionType
		if err := rows.Scan(&t.Code, &t.Name, &t.Description, &t.Created, &t.Updated); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, nil
}

func (r *transactionTypeRepositorySQL) GetByPowerRegionAndCode(ctx context.Context, powerRegionID string, powerRegionTransactionTypeCode string) (*dbentity.TransactionType, error) {
	var t dbentity.TransactionType
	query := `
		SELECT tt.code, tt.name, tt.description, tt.created, tt.updated
		FROM power_region_transaction_type prtt
		JOIN transaction_type tt ON prtt.transaction_type_code = tt.code
		WHERE prtt.power_region_id = $1 AND prtt.code = $2
	`
	err := r.db.QueryRow(ctx, query, powerRegionID, powerRegionTransactionTypeCode).
		Scan(&t.Code, &t.Name, &t.Description, &t.Created, &t.Updated)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
