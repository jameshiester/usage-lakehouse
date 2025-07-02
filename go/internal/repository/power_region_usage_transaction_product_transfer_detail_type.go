package repository

import (
	"context"
	dbentity "usage-lakehouse/internal/db/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PowerRegionUsageTransactionProductTransferDetailTypeRepository interface {
	Create(ctx context.Context, t *dbentity.PowerRegionUsageTransactionProductTransferDetailType) error
	GetByCode(ctx context.Context, code string) (*dbentity.PowerRegionUsageTransactionProductTransferDetailType, error)
	Update(ctx context.Context, t *dbentity.PowerRegionUsageTransactionProductTransferDetailType) error
	Delete(ctx context.Context, code string) error
	List(ctx context.Context) ([]dbentity.PowerRegionUsageTransactionProductTransferDetailType, error)
	MapByCode(ctx context.Context) (map[string]dbentity.PowerRegionUsageTransactionProductTransferDetailType, error)
}

type powerRegionUsageTransactionProductTransferDetailTypeRepositorySQL struct {
	db *pgxpool.Pool
}

func NewPowerRegionUsageTransactionProductTransferDetailTypeRepository(db *pgxpool.Pool) PowerRegionUsageTransactionProductTransferDetailTypeRepository {
	return &powerRegionUsageTransactionProductTransferDetailTypeRepositorySQL{db: db}
}

func (r *powerRegionUsageTransactionProductTransferDetailTypeRepositorySQL) Create(ctx context.Context, t *dbentity.PowerRegionUsageTransactionProductTransferDetailType) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO power_region_usage_transaction_product_transfer_detail_type (code, power_region_id, interval, meter, summary, name, description) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		t.Code, t.PowerRegionID, t.Interval, t.Meter, t.Summary, t.Name, t.Description,
	)
	return err
}

func (r *powerRegionUsageTransactionProductTransferDetailTypeRepositorySQL) GetByCode(ctx context.Context, code string) (*dbentity.PowerRegionUsageTransactionProductTransferDetailType, error) {
	var t dbentity.PowerRegionUsageTransactionProductTransferDetailType
	err := r.db.QueryRow(ctx, `SELECT code, power_region_id, interval, meter, summary, name, description, created, updated FROM power_region_usage_transaction_product_transfer_detail_type WHERE code=$1`, code).
		Scan(&t.Code, &t.PowerRegionID, &t.Interval, &t.Meter, &t.Summary, &t.Name, &t.Description, &t.Created, &t.Updated)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *powerRegionUsageTransactionProductTransferDetailTypeRepositorySQL) Update(ctx context.Context, t *dbentity.PowerRegionUsageTransactionProductTransferDetailType) error {
	_, err := r.db.Exec(ctx, `UPDATE power_region_usage_transaction_product_transfer_detail_type SET power_region_id=$1, interval=$2, meter=$3, summary=$4, name=$5, description=$6 WHERE code=$7`, t.PowerRegionID, t.Interval, t.Meter, t.Summary, t.Name, t.Description, t.Code)
	return err
}

func (r *powerRegionUsageTransactionProductTransferDetailTypeRepositorySQL) Delete(ctx context.Context, code string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM power_region_usage_transaction_product_transfer_detail_type WHERE code=$1`, code)
	return err
}

func (r *powerRegionUsageTransactionProductTransferDetailTypeRepositorySQL) List(ctx context.Context) ([]dbentity.PowerRegionUsageTransactionProductTransferDetailType, error) {
	rows, err := r.db.Query(ctx, `SELECT code, power_region_id, interval, meter, summary, name, description, created, updated FROM power_region_usage_transaction_product_transfer_detail_type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var types []dbentity.PowerRegionUsageTransactionProductTransferDetailType
	for rows.Next() {
		var t dbentity.PowerRegionUsageTransactionProductTransferDetailType
		if err := rows.Scan(&t.Code, &t.PowerRegionID, &t.Interval, &t.Meter, &t.Summary, &t.Name, &t.Description, &t.Created, &t.Updated); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, nil
}

func (r *powerRegionUsageTransactionProductTransferDetailTypeRepositorySQL) MapByCode(ctx context.Context) (map[string]dbentity.PowerRegionUsageTransactionProductTransferDetailType, error) {
	rows, err := r.db.Query(ctx, `SELECT code, power_region_id, interval, meter, summary, name, description, created, updated FROM power_region_usage_transaction_product_transfer_detail_type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]dbentity.PowerRegionUsageTransactionProductTransferDetailType)
	for rows.Next() {
		var t dbentity.PowerRegionUsageTransactionProductTransferDetailType
		if err := rows.Scan(&t.Code, &t.PowerRegionID, &t.Interval, &t.Meter, &t.Summary, &t.Name, &t.Description, &t.Created, &t.Updated); err != nil {
			return nil, err
		}
		result[t.Code] = t
	}
	return result, nil
}
