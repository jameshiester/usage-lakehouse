package repository

import (
	"context"
	dbentity "usage-lakehouse/internal/db/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsageTransactionRepository interface {
	Create(ctx context.Context, t *dbentity.UsageTransaction) error
	GetByID(ctx context.Context, id string) (*dbentity.UsageTransaction, error)
	Update(ctx context.Context, t *dbentity.UsageTransaction) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]dbentity.UsageTransaction, error)
	SaveWithDetails(ctx context.Context, t *dbentity.UsageTransaction, details []dbentity.UsageTransactionDetail) error
}

type usageTransactionRepositorySQL struct {
	db *pgxpool.Pool
}

func NewUsageTransactionRepository(db *pgxpool.Pool) UsageTransactionRepository {
	return &usageTransactionRepositorySQL{db: db}
}

func (r *usageTransactionRepositorySQL) Create(ctx context.Context, t *dbentity.UsageTransaction) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO usage_transaction (id, transaction_id, transaction_type, transaction_sub_type, transaction_date, service_period_start, service_period_end, is_final, is_canceled, purpose, power_region_id, tdspid, premise_id, created, updated) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		t.ID, t.TransactionID, t.TransactionType, t.TransactionSubType, t.TransactionDate, t.ServicePeriodStart, t.ServicePeriodEnd, t.IsFinal, t.IsCanceled, t.Purpose, t.PowerRegionID, t.TDSPID, t.PremiseID, t.Created, t.Updated,
	)
	return err
}

func (r *usageTransactionRepositorySQL) GetByID(ctx context.Context, id string) (*dbentity.UsageTransaction, error) {
	var t dbentity.UsageTransaction
	err := r.db.QueryRow(ctx, `SELECT id, transaction_id, transaction_type, transaction_sub_type, transaction_date, service_period_start, service_period_end, is_final, is_canceled, purpose, power_region_id, tdspid, premise_id, created, updated FROM usage_transaction WHERE id=$1`, id).
		Scan(&t.ID, &t.TransactionID, &t.TransactionType, &t.TransactionSubType, &t.TransactionDate, &t.ServicePeriodStart, &t.ServicePeriodEnd, &t.IsFinal, &t.IsCanceled, &t.Purpose, &t.PowerRegionID, &t.TDSPID, &t.PremiseID, &t.Created, &t.Updated)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *usageTransactionRepositorySQL) Update(ctx context.Context, t *dbentity.UsageTransaction) error {
	_, err := r.db.Exec(ctx, `UPDATE usage_transaction SET transaction_id=$1, transaction_type=$2, transaction_sub_type=$3, transaction_date=$4, service_period_start=$5, service_period_end=$6, is_final=$7, is_canceled=$8, purpose=$9, power_region_id=$10, tdspid=$11, premise_id=$12, updated=$13 WHERE id=$14`, t.TransactionID, t.TransactionType, t.TransactionSubType, t.TransactionDate, t.ServicePeriodStart, t.ServicePeriodEnd, t.IsFinal, t.IsCanceled, t.Purpose, t.PowerRegionID, t.TDSPID, t.PremiseID, t.Updated, t.ID)
	return err
}

func (r *usageTransactionRepositorySQL) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM usage_transaction WHERE id=$1`, id)
	return err
}

func (r *usageTransactionRepositorySQL) List(ctx context.Context) ([]dbentity.UsageTransaction, error) {
	rows, err := r.db.Query(ctx, `SELECT id, transaction_id, transaction_type, transaction_sub_type, transaction_date, service_period_start, service_period_end, is_final, is_canceled, purpose, power_region_id, tdspid, premise_id, created, updated FROM usage_transaction`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var txs []dbentity.UsageTransaction
	for rows.Next() {
		var t dbentity.UsageTransaction
		if err := rows.Scan(&t.ID, &t.TransactionID, &t.TransactionType, &t.TransactionSubType, &t.TransactionDate, &t.ServicePeriodStart, &t.ServicePeriodEnd, &t.IsFinal, &t.IsCanceled, &t.Purpose, &t.PowerRegionID, &t.TDSPID, &t.PremiseID, &t.Created, &t.Updated); err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, nil
}

func (r *usageTransactionRepositorySQL) SaveWithDetails(ctx context.Context, t *dbentity.UsageTransaction, details []dbentity.UsageTransactionDetail) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	_, err = tx.Exec(ctx,
		`INSERT INTO usage_transaction (id, transaction_id, transaction_type, transaction_sub_type, transaction_date, service_period_start, service_period_end, is_final, is_canceled, purpose, power_region_id, tdspid, premise_id, created, updated) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		t.ID, t.TransactionID, t.TransactionType, t.TransactionSubType, t.TransactionDate, t.ServicePeriodStart, t.ServicePeriodEnd, t.IsFinal, t.IsCanceled, t.Purpose, t.PowerRegionID, t.TDSPID, t.PremiseID, t.Created, t.Updated,
	)
	if err != nil {
		return err
	}
	for _, d := range details {
		_, err = tx.Exec(ctx,
			`INSERT INTO usage_transaction_detail (usage_transaction_id, meter_id, meter_name, power_region_id, premise_id, is_canceled, service_period_start, service_period_end, consumption, production, created, updated) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
			d.UsageTransactionID, d.MeterID, d.MeterName, d.PowerRegionID, d.PremiseID, d.IsCanceled, d.ServicePeriodStart, d.ServicePeriodEnd, d.Consumption, d.Production, d.Created, d.Updated,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
