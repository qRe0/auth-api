package app

import (
	"fmt"
	"log"
	"net"

	"github.com/qRe0/auth-api/configs"
	"github.com/qRe0/auth-api/internal/handlers"
	"github.com/qRe0/auth-api/internal/migrations"
	authRepo "github.com/qRe0/auth-api/internal/repository/auth"
	tokenRepo "github.com/qRe0/auth-api/internal/repository/token"
	authServ "github.com/qRe0/auth-api/internal/service/auth"
	tokenServ "github.com/qRe0/auth-api/internal/service/token"
	pb "github.com/qRe0/auth-api/proto/gen/go"
	"google.golang.org/grpc"
)

func Run() {
	cfg, err := configs.LoadEnv()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := authRepo.Init(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}

	cache, err := tokenRepo.Init(cfg.Redis)
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

	authRepository := authRepo.NewAuthRepository(db)
	tokenRepository := tokenRepo.NewTokenRepo(cfg.JWT, cache)
	tokenService := tokenServ.NewTokenService(tokenRepository)
	authService := authServ.NewAuthService(cfg.JWT, authRepository, tokenService)
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
