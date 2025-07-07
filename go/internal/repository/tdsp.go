package repository

import (
	"context"
	dbentity "usage-lakehouse/internal/db/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TDSPRepository interface {
	Create(ctx context.Context, t *dbentity.TDSP) error
	GetByID(ctx context.Context, id string) (*dbentity.TDSP, error)
	GetByLegalID(ctx context.Context, legalID string) (*dbentity.TDSP, error)
	GetByName(ctx context.Context, name string) (*dbentity.TDSP, error)
	Update(ctx context.Context, t *dbentity.TDSP) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]dbentity.TDSP, error)
}

type tdspRepositorySQL struct {
	db *pgxpool.Pool
}

func NewTDSPRepository(db *pgxpool.Pool) TDSPRepository {
	return &tdspRepositorySQL{db: db}
}

func (r *tdspRepositorySQL) Create(ctx context.Context, t *dbentity.TDSP) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO tdsp (id, name, code, legal_id, premise_code_validation_expression) VALUES ($1, $2, $3, $4, $5)`,
		t.ID, t.Name, t.Code, t.LegalID, t.PremiseCodeValidationExpression,
	)
	return err
}

func (r *tdspRepositorySQL) GetByID(ctx context.Context, id string) (*dbentity.TDSP, error) {
	var t dbentity.TDSP
	err := r.db.QueryRow(ctx, `SELECT id, name, code, legal_id, premise_code_validation_expression, created, updated FROM tdsp WHERE id=$1`, id).
		Scan(&t.ID, &t.Name, &t.Code, &t.LegalID, &t.PremiseCodeValidationExpression, &t.Created, &t.Updated)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tdspRepositorySQL) GetByLegalID(ctx context.Context, legalID string) (*dbentity.TDSP, error) {
	var t dbentity.TDSP
	err := r.db.QueryRow(ctx, `SELECT id, name, code, legal_id, premise_code_validation_expression, created, updated FROM tdsp WHERE legal_id=$1`, legalID).
		Scan(&t.ID, &t.Name, &t.Code, &t.LegalID, &t.PremiseCodeValidationExpression, &t.Created, &t.Updated)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tdspRepositorySQL) GetByName(ctx context.Context, name string) (*dbentity.TDSP, error) {
	var t dbentity.TDSP
	err := r.db.QueryRow(ctx, `SELECT id, name, code, legal_id, premise_code_validation_expression, created, updated FROM tdsp WHERE name=$1`, name).
		Scan(&t.ID, &t.Name, &t.Code, &t.LegalID, &t.PremiseCodeValidationExpression, &t.Created, &t.Updated)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tdspRepositorySQL) Update(ctx context.Context, t *dbentity.TDSP) error {
	_, err := r.db.Exec(ctx, `UPDATE tdsp SET name=$1, code=$2, legal_id=$3, premise_code_validation_expression=$4 WHERE id=$5`, t.Name, t.Code, t.LegalID, t.PremiseCodeValidationExpression, t.ID)
	return err
}

func (r *tdspRepositorySQL) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tdsp WHERE id=$1`, id)
	return err
}

func (r *tdspRepositorySQL) List(ctx context.Context) ([]dbentity.TDSP, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, code, legal_id, premise_code_validation_expression, created, updated FROM tdsp`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tdsps []dbentity.TDSP
	for rows.Next() {
		var t dbentity.TDSP
		if err := rows.Scan(&t.ID, &t.Name, &t.Code, &t.LegalID, &t.PremiseCodeValidationExpression, &t.Created, &t.Updated); err != nil {
			return nil, err
		}
		tdsps = append(tdsps, t)
	}
	return tdsps, nil
}
