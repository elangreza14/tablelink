package main

import (
	"context"
	"fmt"
	"log"
	"net"

	gen "github.com/elangreza14/tablelink/gen/go"
	"github.com/elangreza14/tablelink/repository"
	"github.com/elangreza14/tablelink/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {

	conn := "postgres://tablelink:tablelink@localhost:5432/tablelink?sslmode=disable"

	config, err := pgxpool.ParseConfig(conn)
	if err != nil {
		log.Fatal(err)
	}

	pgxPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	authRepo := repository.NewAuthRepo(pgxPool)
	roleRightsRepo := repository.NewRoleRightRepo(pgxPool)
	authSrv := services.NewAuthService(authRepo, roleRightsRepo)
	grpcServer := grpc.NewServer()
	gen.RegisterAuthServer(grpcServer, authSrv)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("grpc run @ port 50051")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
