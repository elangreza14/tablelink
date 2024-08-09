package repository

import (
	"context"

	"github.com/elangreza14/tablelink/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(
	db *pgxpool.Pool,
) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (a *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := `SELECT id, email, "password", "name", "role_id" FROM public.users
			WHERE email=$1;`
	res := &domain.User{}

	err := a.db.QueryRow(ctx, q, email).Scan(&res.ID, &res.Email, &res.Password, &res.Name, &res.RoleID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *AuthRepo) CreateUser(
	ctx context.Context,
	roleId int,
	name,
	email,
	password string,
) error {
	q := `INSERT INTO public.users
			(role_id, email, "password", "name")
			VALUES($1,$2,$3,$4);`

	_, err := a.db.Exec(ctx, q, roleId, name, email, password)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthRepo) UpdateUser(
	ctx context.Context,
	email string,
	name string,
) error {
	q := `UPDATE public.users
			SET name=$2
			WHERE email=$1;`

	_, err := a.db.Exec(ctx, q, email, name)
	if err != nil {
		return err
	}

	return nil
}
