package app

import (
	"grpc-auth/pkg/api/services"
	"grpc-auth/pkg/pb"
	"grpc-auth/pkg/utils"
)

type Server struct {
	jwt         utils.JwtWrapper
	authService services.AuthService
	pb.UnimplementedAuthServiceServer
	//mustEmbedUnimplementedAuthServiceServer
}

func NewServer(jwt utils.JwtWrapper, authService services.AuthService) *Server {
	return &Server{jwt: jwt, authService: authService}
}
