package repository

import (
	"context"
	dbentity "usage-lakehouse/internal/db/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PowerRegionRepository interface {
	Create(ctx context.Context, p *dbentity.PowerRegion) error
	GetByID(ctx context.Context, id string) (*dbentity.PowerRegion, error)
	GetByName(ctx context.Context, id string) (*dbentity.PowerRegion, error)
	Update(ctx context.Context, p *dbentity.PowerRegion) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]dbentity.PowerRegion, error)
}

type powerRegionRepositorySQL struct {
	db *pgxpool.Pool
}

func NewPowerRegionRepository(db *pgxpool.Pool) PowerRegionRepository {
	return &powerRegionRepositorySQL{db: db}
}

func (r *powerRegionRepositorySQL) Create(ctx context.Context, p *dbentity.PowerRegion) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO power_region (id, name) VALUES ($1, $2)`,
		p.ID, p.Name,
	)
	return err
}

func (r *powerRegionRepositorySQL) GetByID(ctx context.Context, id string) (*dbentity.PowerRegion, error) {
	var p dbentity.PowerRegion
	err := r.db.QueryRow(ctx, `SELECT id, name, created_dttm, updated_dttm FROM power_region WHERE id=$1`, id).
		Scan(&p.ID, &p.Name, &p.Created, &p.Updated)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *powerRegionRepositorySQL) GetByName(ctx context.Context, name string) (*dbentity.PowerRegion, error) {
	var p dbentity.PowerRegion
	err := r.db.QueryRow(ctx, `SELECT id, name, created_dttm, updated_dttm FROM power_region WHERE name=$1`, name).
		Scan(&p.ID, &p.Name, &p.Created, &p.Updated)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *powerRegionRepositorySQL) Update(ctx context.Context, p *dbentity.PowerRegion) error {
	_, err := r.db.Exec(ctx, `UPDATE power_region SET name=$1 WHERE id=$2`, p.Name, p.ID)
	return err
}

func (r *powerRegionRepositorySQL) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM power_region WHERE id=$1`, id)
	return err
}

func (r *powerRegionRepositorySQL) List(ctx context.Context) ([]dbentity.PowerRegion, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, created_dttm, updated_dttm FROM power_region`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var regions []dbentity.PowerRegion
	for rows.Next() {
		var p dbentity.PowerRegion
		if err := rows.Scan(&p.ID, &p.Name, &p.Created, &p.Updated); err != nil {
			return nil, err
		}
		regions = append(regions, p)
	}
	return regions, nil
}
