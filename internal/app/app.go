package app

import (
	"fmt"
	"log"
	"net"

	"github.com/qRe0/auth-api/configs"
	"github.com/qRe0/auth-api/internal/handlers"
	"github.com/qRe0/auth-api/internal/migrations"
	repository "github.com/qRe0/auth-api/internal/repository/auth"
	service "github.com/qRe0/auth-api/internal/service/auth"
	pb "github.com/qRe0/auth-api/proto/gen/go"
	"google.golang.org/grpc"
)

func Run() {
	cfg, err := configs.LoadEnv()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := repository.Init(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}

	migrator, err := migrations.NewMigrator(db)
	if err != nil {
		log.Fatalln(err)
	}

	err = migrator.Latest()
	if err != nil {
		log.Fatalln(err)
	}

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(cfg.JWT, authRepo)
	handler := handlers.NewAuthHandler(authService, cfg.JWT, ":50051")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSignUpServer(grpcServer, handler)

	fmt.Printf("gRPC server is running on %v\n", lis.Addr().String())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
