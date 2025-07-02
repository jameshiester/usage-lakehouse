package repository

import (
	"context"
	dbentity "usage-lakehouse/internal/db/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsageTransactionPurposeRepository interface {
	GetByCode(ctx context.Context, code string) (*dbentity.UsageTransactionPurpose, error)
	List(ctx context.Context) ([]dbentity.UsageTransactionPurpose, error)
	GetByPowerRegionAndCode(ctx context.Context, powerRegionID string, code string) (*dbentity.UsageTransactionPurpose, error)
}

type usageTransactionPurposeRepositorySQL struct {
	db *pgxpool.Pool
}

func NewUsageTransactionPurposeRepository(db *pgxpool.Pool) UsageTransactionPurposeRepository {
	return &usageTransactionPurposeRepositorySQL{db: db}
}

func (r *usageTransactionPurposeRepositorySQL) GetByCode(ctx context.Context, code string) (*dbentity.UsageTransactionPurpose, error) {
	var p dbentity.UsageTransactionPurpose
	err := r.db.QueryRow(ctx, `SELECT code, name, is_cancel, description, created_dttm, updated_dttm FROM usage_transaction_purpose WHERE code=$1`, code).
		Scan(&p.Code, &p.Name, &p.IsCancel, &p.Description, &p.Created, &p.Updated)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *usageTransactionPurposeRepositorySQL) List(ctx context.Context) ([]dbentity.UsageTransactionPurpose, error) {
	rows, err := r.db.Query(ctx, `SELECT code, name, is_cancel, description, created_dttm, updated_dttm FROM usage_transaction_purpose`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var purposes []dbentity.UsageTransactionPurpose
	for rows.Next() {
		var p dbentity.UsageTransactionPurpose
		if err := rows.Scan(&p.Code, &p.Name, &p.IsCancel, &p.Description, &p.Created, &p.Updated); err != nil {
			return nil, err
		}
		purposes = append(purposes, p)
	}
	return purposes, nil
}

func (r *usageTransactionPurposeRepositorySQL) GetByPowerRegionAndCode(ctx context.Context, powerRegionID string, code string) (*dbentity.UsageTransactionPurpose, error) {
	var p dbentity.UsageTransactionPurpose
	err := r.db.QueryRow(ctx, `
		SELECT utp.code, utp.name, utp.is_cancel, utp.description, utp.created_dttm, utp.updated_dttm
		FROM usage_transaction_purpose utp
		JOIN power_region_usage_transaction_purpose prutp ON utp.code = prutp.usage_transaction_purpose_code
		WHERE prutp.power_region_id = $1 AND prutp.code = $2
	`, powerRegionID, code).Scan(&p.Code, &p.Name, &p.IsCancel, &p.Description, &p.Created, &p.Updated)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
