package repository

import (
	"context"
	dbentity "usage-lakehouse/internal/db/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionSubTypeRepository interface {
	Create(ctx context.Context, t *dbentity.TransactionSubType) error
	GetByCode(ctx context.Context, code string) (*dbentity.TransactionSubType, error)
	Update(ctx context.Context, t *dbentity.TransactionSubType) error
	Delete(ctx context.Context, code string) error
	List(ctx context.Context) ([]dbentity.TransactionSubType, error)
	GetByPowerRegionSubTypeCode(ctx context.Context, powerRegionID string, powerRegionTransactionSubTypeCode string) (*dbentity.TransactionSubType, error)
}

type transactionSubTypeRepositorySQL struct {
	db *pgxpool.Pool
}

func NewTransactionSubTypeRepository(db *pgxpool.Pool) TransactionSubTypeRepository {
	return &transactionSubTypeRepositorySQL{db: db}
}

func (r *transactionSubTypeRepositorySQL) Create(ctx context.Context, t *dbentity.TransactionSubType) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO transaction_sub_type (code, transaction_type_code, name, description) VALUES ($1, $2, $3, $4)`,
		t.Code, t.TransactionTypeCode, t.Name, t.Description,
	)
	return err
}

func (r *transactionSubTypeRepositorySQL) GetByCode(ctx context.Context, code string) (*dbentity.TransactionSubType, error) {
	var t dbentity.TransactionSubType
	err := r.db.QueryRow(ctx, `SELECT code, transaction_type_code, name, description, created, updated FROM transaction_sub_type WHERE code=$1`, code).
		Scan(&t.Code, &t.TransactionTypeCode, &t.Name, &t.Description, &t.Created, &t.Updated)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *transactionSubTypeRepositorySQL) Update(ctx context.Context, t *dbentity.TransactionSubType) error {
	_, err := r.db.Exec(ctx, `UPDATE transaction_sub_type SET transaction_type_code=$1, name=$2, description=$3 WHERE code=$4`, t.TransactionTypeCode, t.Name, t.Description, t.Code)
	return err
}

func (r *transactionSubTypeRepositorySQL) Delete(ctx context.Context, code string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM transaction_sub_type WHERE code=$1`, code)
	return err
}

func (r *transactionSubTypeRepositorySQL) List(ctx context.Context) ([]dbentity.TransactionSubType, error) {
	rows, err := r.db.Query(ctx, `SELECT code, transaction_type_code, name, description, created, updated FROM transaction_sub_type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subTypes []dbentity.TransactionSubType
	for rows.Next() {
		var t dbentity.TransactionSubType
		if err := rows.Scan(&t.Code, &t.TransactionTypeCode, &t.Name, &t.Description, &t.Created, &t.Updated); err != nil {
			return nil, err
		}
		subTypes = append(subTypes, t)
	}
	return subTypes, nil
}

func (r *transactionSubTypeRepositorySQL) GetByPowerRegionSubTypeCode(ctx context.Context, powerRegionID string, powerRegionTransactionSubTypeCode string) (*dbentity.TransactionSubType, error) {
	var t dbentity.TransactionSubType
	query := `
		SELECT tst.code, tst.transaction_type_code, tst.name, tst.description, tst.created, tst.updated
		FROM power_region_transaction_sub_type prtst
		JOIN transaction_sub_type tst ON prtst.transaction_sub_type_code = tst.code
		WHERE prtst.power_region_id = $1 AND prtst.code = $2
	`
	err := r.db.QueryRow(ctx, query, powerRegionID, powerRegionTransactionSubTypeCode).
		Scan(&t.Code, &t.TransactionTypeCode, &t.Name, &t.Description, &t.Created, &t.Updated)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
