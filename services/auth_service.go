package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/elangreza14/tablelink/domain"
	gen "github.com/elangreza14/tablelink/gen/go"
	"google.golang.org/grpc/metadata"
)

type (
	AuthRepo interface {
		GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
		CreateUser(
			ctx context.Context,
			roleId int,
			name,
			email,
			password string,
		) error
		UpdateUser(
			ctx context.Context,
			id string,
			name string,
		) error
	}
	RoleRightsRepo interface {
		GetRoleRightsByRoleID(ctx context.Context, roleID int) (*domain.RoleRights, error)
	}

	AuthService struct {
		gen.UnimplementedAuthServer
		roleRightsRepo RoleRightsRepo
		authRepo       AuthRepo
	}
)

func NewAuthService(authRepo AuthRepo, roleRightsRepo RoleRightsRepo) *AuthService {
	return &AuthService{
		UnimplementedAuthServer: gen.UnimplementedAuthServer{},
		roleRightsRepo:          roleRightsRepo,
		authRepo:                authRepo,
	}
}

func (a *AuthService) LoginUser(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	user, err := a.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user.Password != req.Password {
		return nil, errors.New("password incorrect")
	}

	return &gen.LoginResponse{
		Status:  true,
		Message: "Successfully",
		Data: &gen.Data{
			// TODO jwt
			AccessToken: req.Email,
		},
	}, nil
}

func (a *AuthService) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {

	user, err := a.validateToken(ctx)
	if err != nil {
		return nil, err
	}

	roleRights, err := a.roleRightsRepo.GetRoleRightsByRoleID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	if roleRights.RCreate == 0 {
		return nil, errors.New("this user cannot create user")
	}

	err = a.authRepo.CreateUser(ctx, int(req.RoleId), req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &gen.CreateUserResponse{
		Status:  true,
		Message: "Successfully",
	}, nil
}

func (a *AuthService) UpdUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {

	user, err := a.validateToken(ctx)
	if err != nil {
		return nil, err
	}

	roleRights, err := a.roleRightsRepo.GetRoleRightsByRoleID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	if roleRights.RCreate == 0 {
		return nil, errors.New("this user cannot create user")
	}

	err = a.authRepo.CreateUser(ctx, int(req.RoleId), req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &gen.CreateUserResponse{
		Status:  true,
		Message: "Successfully",
	}, nil
}

func (a *AuthService) UdateUser(ctx context.Context, req *gen.UpdateUserRequest) (*gen.UpdateUserResponse, error) {
	user, err := a.validateToken(ctx)
	if err != nil {
		return nil, err
	}

	roleRights, err := a.roleRightsRepo.GetRoleRightsByRoleID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	if roleRights.RUpdate == 0 {
		return nil, errors.New("this user cannot update user")
	}

	err = a.authRepo.UpdateUser(ctx, user.Email, req.Name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &gen.UpdateUserResponse{
		Status:  true,
		Message: "Successfully",
	}, nil
}

func (a *AuthService) validateToken(ctx context.Context) (*domain.User, error) {

	tokenMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("error reading metadata")
	}

	tokenRaw, ok := tokenMetadata["x-link-service"]
	if !ok {
		return nil, errors.New("error reading X-Link-Service")
	}

	if len(tokenRaw) == 0 {
		return nil, errors.New("not valid token")
	}

	tokens := strings.Split(tokenRaw[0], " ")

	if len(tokens) != 2 {
		return nil, errors.New("not valid token")
	}

	return a.authRepo.GetUserByEmail(ctx, tokens[1])
}
