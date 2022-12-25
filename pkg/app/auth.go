package app

import (
	"context"
	"grpc-auth/pkg/api/models"
	"grpc-auth/pkg/pb"
	"net/http"
)

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	//var user models.User
	user := new(models.User)
	var err error

	user.UserName = req.Username
	user.Password = req.Password

	if err = s.authService.Login(user); err != nil {
		return &pb.LoginResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	token, err := s.jwt.GenerateToken(user)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotAcceptable,
			Error:  err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Error:  "",
		Token:  token,
	}, nil
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User
	var err error

	user.UserName = req.Username
	user.Password = req.Password

	if err = s.authService.Register(user); err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}
	return &pb.RegisterResponse{
		Status: http.StatusOK,
		Error:  "",
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	//var user models.User
	user := new(models.User)
	var err error

	token, err := s.jwt.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	user.UserName = token.UserName

	if err = s.authService.Validate(user); err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}
	return &pb.ValidateResponse{
		Status: http.StatusOK,
		Error:  "",
		UserId: int32(token.Id),
	}, nil
}
