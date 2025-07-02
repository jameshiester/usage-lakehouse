// MOVE FILE to go/internal/repository/account_repository.go
// (no code changes needed, just move the file and ensure the package name is 'repository')

package repository

import (
	"context"
	"usage-lakehouse/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository interface {
	Create(ctx context.Context, a *model.Account) error
	GetByID(ctx context.Context, id string) (*model.Account, error)
	Update(ctx context.Context, a *model.Account) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.Account, error)
	ExistsByName(ctx context.Context, name string, accountID *string) (bool, error)
	ExistsByLegalID(ctx context.Context, name string, accountID *string) (bool, error)
}

type accountRepositorySQL struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) AccountRepository {
	return &accountRepositorySQL{db: db}
}

func (r *accountRepositorySQL) Create(ctx context.Context, a *model.Account) error {
	a.ID = uuid.New().String()
	return r.db.QueryRow(ctx,
		`INSERT INTO account (id, legal_id, name) VALUES ($1, $2, $3) RETURNING id`,
		a.ID, a.LegalID, a.Name,
	).Scan(&a.ID)
}

func (r *accountRepositorySQL) GetByID(ctx context.Context, id string) (*model.Account, error) {
	var a model.Account
	err := r.db.QueryRow(ctx, `SELECT id, legal_id, name, created_dttm, updated_dttm FROM account WHERE id=$1`, id).Scan(&a.ID, &a.LegalID, &a.Name, &a.Created, &a.Updated)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *accountRepositorySQL) Update(ctx context.Context, a *model.Account) error {
	_, err := r.db.Exec(ctx, `UPDATE account SET name=$1 WHERE id=$2`, a.Name, a.ID)
	return err
}

func (r *accountRepositorySQL) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM account WHERE id=$1`, id)
	return err
}

func (r *accountRepositorySQL) List(ctx context.Context) ([]model.Account, error) {
	rows, err := r.db.Query(ctx, `SELECT id, legal_id, name, created_dttm, updated_dttm FROM public.account`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []model.Account
	for rows.Next() {
		var a model.Account
		if err := rows.Scan(&a.ID, &a.LegalID, &a.Name, &a.Created, &a.Updated); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (r *accountRepositorySQL) ExistsByName(ctx context.Context, name string, accountID *string) (bool, error) {
	var exists bool
	var err error
	sql := `SELECT EXISTS(SELECT 1 FROM account WHERE name=$1)`

	if accountID != nil {
		sql += ` AND id != $2`
		err = r.db.QueryRow(ctx, sql, name, accountID).Scan(&exists)
	} else {
		err = r.db.QueryRow(ctx, sql, name).Scan(&exists)
	}
	return exists, err
}

func (r *accountRepositorySQL) ExistsByLegalID(ctx context.Context, legalID string, accountID *string) (bool, error) {
	var exists bool
	var err error
	sql := `SELECT EXISTS(SELECT 1 FROM account WHERE legal_id=$1)`

	if accountID != nil {
		sql += ` AND id != $2`
		err = r.db.QueryRow(ctx, sql, legalID, accountID).Scan(&exists)
	} else {
		err = r.db.QueryRow(ctx, sql, legalID).Scan(&exists)
	}
	return exists, err
}
