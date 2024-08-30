package handlers

import (
	"context"
	"log"
	"net"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/qRe0/auth-api/configs"
	"github.com/qRe0/auth-api/internal/models"
	authService "github.com/qRe0/auth-api/internal/service/auth"
	pb "github.com/qRe0/auth-api/proto/gen/go"
	"google.golang.org/grpc"
	md "google.golang.org/grpc/metadata"
)

type AuthHandler struct {
	service authService.AuthServiceInterface
	cfg     configs.JWTConfig
	pb.UnimplementedSignUpServer
	pb.UnimplementedLogInServer
	pb.UnimplementedRefreshServer
	pb.UnimplementedRevokeServer
	pb.UnimplementedAuthMiddlewareServer
}

func NewAuthHandler(service authService.AuthServiceInterface, cfg configs.JWTConfig, address string) *AuthHandler {
	handler := &AuthHandler{
		service: service,
		cfg:     cfg,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterSignUpServer(grpcServer, handler)
	pb.RegisterLogInServer(grpcServer, handler)
	pb.RegisterRefreshServer(grpcServer, handler)
	pb.RegisterRevokeServer(grpcServer, handler)
	pb.RegisterAuthMiddlewareServer(grpcServer, handler)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %v: %v", address, err)
	}

	log.Printf("gRPC server is running on %v", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	return handler
}

func (a *AuthHandler) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
	}

	tokens, err := a.service.SignUp(ctx, user)
	if err != nil {
		return nil, err
	}

	metadata := md.Pairs(
		"Authorization", tokens.AccessToken,
		"Refresh-Token", tokens.RefreshToken,
	)
	err = grpc.SendHeader(ctx, metadata)
	if err != nil {
		return nil, err
	}

	resp := &pb.SignUpResponse{
		Message: "User created successfully!",
	}
	return resp, nil
}

func (a *AuthHandler) LogIn(ctx context.Context, req *pb.LogInRequest) (*pb.LogInResponse, error) {
	user := &models.User{
		Phone:    req.Phone,
		Password: req.Password,
	}

	tokens, err := a.service.LogIn(ctx, user)
	if err != nil {
		return nil, err
	}

	metadata := md.Pairs(
		"Authorization", tokens.AccessToken,
		"Refresh-Token", tokens.RefreshToken,
	)
	err = grpc.SendHeader(ctx, metadata)
	if err != nil {
		return nil, err
	}

	resp := &pb.LogInResponse{
		Message: "User logged in successfully!",
	}

	return resp, nil
}

func (a *AuthHandler) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	refreshToken := req.RefreshToken

	tokens, err := a.service.Refresh(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	metadata := md.Pairs(
		"Authorization", tokens.AccessToken,
		"Refresh-Token", tokens.RefreshToken,
	)
	err = grpc.SendHeader(ctx, metadata)
	if err != nil {
		return nil, err
	}

	resp := &pb.RefreshResponse{
		Message: "Token refreshed successfully!",
	}

	return resp, nil
}

func (a *AuthHandler) Revoke(ctx context.Context, req *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	user := &models.User{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
	}

	err := a.service.RevokeTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	resp := &pb.RevokeResponse{
		Message: "Tokens revoked!",
	}

	return resp, nil
}

func (a *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	token := req.Token

	userID, err := a.service.ValidateToken(token, a.cfg.SecretKey)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, err
	}

	blacklisted, err := a.service.TokenBlacklisted(ctx, token)
	if err != nil && !errors.Is(err, redis.Nil) {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, err
	}

	if blacklisted {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	resp := &pb.ValidateTokenResponse{
		UserId: userID,
		Valid:  true,
	}

	return resp, nil
}
