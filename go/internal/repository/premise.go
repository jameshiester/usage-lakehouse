package repository

import (
	"context"
	"usage-lakehouse/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PremiseRepository interface {
	Create(ctx context.Context, p *model.Premise) error
	GetByID(ctx context.Context, id string) (*model.Premise, error)
	GetByCode(ctx context.Context, code string) (*model.Premise, error)
	Update(ctx context.Context, p *model.Premise) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.Premise, error)
}

type premiseRepositorySQL struct {
	db *pgxpool.Pool
}

func NewPremiseRepository(db *pgxpool.Pool) PremiseRepository {
	return &premiseRepositorySQL{db: db}
}

func (r *premiseRepositorySQL) Create(ctx context.Context, p *model.Premise) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO premise (id, code, name, customer_name, address_line_1, city, state, zip, country, premise_type_id, power_region_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		p.ID, p.Code, p.Name, p.CustomerName, p.Address.AddressLine1, p.Address.City, p.Address.State, p.Address.Zip, p.Address.Country, p.PremiseTypeID, p.PowerRegionID,
	)
	return err
}

func (r *premiseRepositorySQL) GetByID(ctx context.Context, id string) (*model.Premise, error) {
	var p model.Premise
	err := r.db.QueryRow(ctx, `SELECT id, code, name, customer_name, address_line_1, city, state, zip, country, premise_type_id, power_region_id, created_dttm, updated_dttm FROM premise WHERE id=$1`, id).
		Scan(&p.ID, &p.Code, &p.Name, &p.CustomerName, &p.Address.AddressLine1, &p.Address.City, &p.Address.State, &p.Address.Zip, &p.Address.Country, &p.PremiseTypeID, &p.PowerRegionID, &p.Created, &p.Updated)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *premiseRepositorySQL) GetByCode(ctx context.Context, code string) (*model.Premise, error) {
	var p model.Premise
	err := r.db.QueryRow(ctx, `SELECT id, code, name, customer_name, address_line_1, city, state, zip, country, premise_type_id, power_region_id, created_dttm, updated_dttm FROM premise WHERE code=$1`, code).
		Scan(&p.ID, &p.Code, &p.Name, &p.CustomerName, &p.Address.AddressLine1, &p.Address.City, &p.Address.State, &p.Address.Zip, &p.Address.Country, &p.PremiseTypeID, &p.PowerRegionID, &p.Created, &p.Updated)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *premiseRepositorySQL) Update(ctx context.Context, p *model.Premise) error {
	_, err := r.db.Exec(ctx, `UPDATE premise SET code=$1, name=$2, customer_name=$3, address_line_1=$4, city=$5, state=$6, zip=$7, country=$8, premise_type_id=$9, power_region_id=$10 WHERE id=$11`, p.Code, p.Name, p.CustomerName, p.Address.AddressLine1, p.Address.City, p.Address.State, p.Address.Zip, p.Address.Country, p.PremiseTypeID, p.PowerRegionID, p.ID)
	return err
}

func (r *premiseRepositorySQL) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM premise WHERE id=$1`, id)
	return err
}

func (r *premiseRepositorySQL) List(ctx context.Context) ([]model.Premise, error) {
	rows, err := r.db.Query(ctx, `SELECT id, code, name, customer_name, address_line_1, city, state, zip, country, premise_type_id, power_region_id, created_dttm, updated_dttm FROM premise`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var premises []model.Premise
	for rows.Next() {
		var p model.Premise
		if err := rows.Scan(&p.ID, &p.Code, &p.Name, &p.CustomerName, &p.Address.AddressLine1, &p.Address.City, &p.Address.State, &p.Address.Zip, &p.Address.Country, &p.PremiseTypeID, &p.PowerRegionID, &p.Created, &p.Updated); err != nil {
			return nil, err
		}
		premises = append(premises, p)
	}
	return premises, nil
}
