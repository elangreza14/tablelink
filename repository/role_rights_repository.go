package repository

import (
	"context"

	"github.com/elangreza14/tablelink/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRightRepo struct {
	db *pgxpool.Pool
}

func NewRoleRightRepo(
	db *pgxpool.Pool,
) *RoleRightRepo {
	return &RoleRightRepo{
		db: db,
	}
}

func (a *RoleRightRepo) GetRoleRightsByRoleID(ctx context.Context, roleID int) (*domain.RoleRights, error) {
	q := ` SELECT id, role_id, route, "section", "path", r_create, r_read, r_update, r_delete
			FROM public.role_rights
			WHERE role_id=$1;`
	res := &domain.RoleRights{}

	err := a.db.QueryRow(ctx, q, roleID).Scan(
		&res.ID,
		&res.RoleID,
		&res.Route,
		&res.Section,
		&res.Path,
		&res.RCreate,
		&res.RRead,
		&res.RUpdate,
		&res.RDelete,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
