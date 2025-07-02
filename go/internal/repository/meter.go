package repository

import (
	"context"
	dbentity "usage-lakehouse/internal/db/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MeterRepository interface {
	Create(ctx context.Context, m *dbentity.Meter) error
	GetByID(ctx context.Context, id string) (*dbentity.Meter, error)
	GetByName(ctx context.Context, powerRegionID string, name string) (*dbentity.Meter, error)
	Update(ctx context.Context, m *dbentity.Meter) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]dbentity.Meter, error)
	GetNameIDMap(ctx context.Context, names []string) (map[string]string, error)
}

type meterRepositorySQL struct {
	db *pgxpool.Pool
}

func NewMeterRepository(db *pgxpool.Pool) MeterRepository {
	return &meterRepositorySQL{db: db}
}

func (r *meterRepositorySQL) Create(ctx context.Context, m *dbentity.Meter) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO meter (id, account_id, premise_id, power_region_id, name, type, load_profile, cycle_code, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		m.ID, m.AccountID, m.PremiseID, m.PowerRegionID, m.Name, m.Type, m.LoadProfile, m.CycleCode, m.Active,
	)
	return err
}

func (r *meterRepositorySQL) GetByID(ctx context.Context, id string) (*dbentity.Meter, error) {
	var m dbentity.Meter
	err := r.db.QueryRow(ctx, `SELECT id, account_id, premise_id, power_region_id, name, type, load_profile, cycle_code, active, created, updated FROM meter WHERE id=$1`, id).
		Scan(&m.ID, &m.AccountID, &m.PremiseID, &m.PowerRegionID, &m.Name, &m.Type, &m.LoadProfile, &m.CycleCode, &m.Active, &m.Created, &m.Updated)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *meterRepositorySQL) GetByName(ctx context.Context, powerRegionID string, name string) (*dbentity.Meter, error) {
	var m dbentity.Meter
	err := r.db.QueryRow(ctx, `SELECT id, account_id, premise_id, power_region_id, name, type, load_profile, cycle_code, active, created, updated FROM meter WHERE name=$1 AND power_region_id=$2`, name, powerRegionID).
		Scan(&m.ID, &m.AccountID, &m.PremiseID, &m.PowerRegionID, &m.Name, &m.Type, &m.LoadProfile, &m.CycleCode, &m.Active, &m.Created, &m.Updated)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *meterRepositorySQL) Update(ctx context.Context, m *dbentity.Meter) error {
	_, err := r.db.Exec(ctx, `UPDATE meter SET account_id=$1, premise_id=$2, power_region_id=$3, name=$4, type=$5, load_profile=$6, cycle_code=$7, active=$8 WHERE id=$9`, m.AccountID, m.PremiseID, m.PowerRegionID, m.Name, m.Type, m.LoadProfile, m.CycleCode, m.Active, m.ID)
	return err
}

func (r *meterRepositorySQL) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM meter WHERE id=$1`, id)
	return err
}

func (r *meterRepositorySQL) List(ctx context.Context) ([]dbentity.Meter, error) {
	rows, err := r.db.Query(ctx, `SELECT id, account_id, premise_id, power_region_id, name, type, load_profile, cycle_code, active, created, updated FROM meter`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var meters []dbentity.Meter
	for rows.Next() {
		var m dbentity.Meter
		if err := rows.Scan(&m.ID, &m.AccountID, &m.PremiseID, &m.PowerRegionID, &m.Name, &m.Type, &m.LoadProfile, &m.CycleCode, &m.Active, &m.Created, &m.Updated); err != nil {
			return nil, err
		}
		meters = append(meters, m)
	}
	return meters, nil
}

func (r *meterRepositorySQL) GetNameIDMap(ctx context.Context, names []string) (map[string]string, error) {
	if len(names) == 0 {
		return map[string]string{}, nil
	}

	// Build the SQL IN clause dynamically
	query := `SELECT DISTINCT name, id FROM meter WHERE name = ANY($1)`
	rows, err := r.db.Query(ctx, query, names)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]string)
	for rows.Next() {
		var name, id string
		if err := rows.Scan(&name, &id); err != nil {
			return nil, err
		}
		result[name] = id
	}
	return result, nil
}
