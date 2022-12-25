package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"grpc-auth/pkg/api/models"
	"grpc-auth/pkg/api/services"
	"grpc-auth/pkg/app"
	"grpc-auth/pkg/pb"
	"grpc-auth/pkg/repository"
	"grpc-auth/pkg/utils"
	"log"
	"net"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var err error

	if err = godotenv.Load("../../configs/.env"); err != nil {
		return err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	db, err := setupDatabase(dsn)
	if err != nil {
		return err
	}
	fmt.Println("server[auth]-database connection established successfully.")

	if err = db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	fmt.Println("server[auth]-migrations completed.")

	lis, err := setupNet(os.Getenv("PORT"))
	if err != nil {
		return err
	}
	fmt.Println("server[auth]-started on ", os.Getenv("PORT"))

	jwt := utils.JwtWrapper{
		SecretKey: os.Getenv("SECRET_KEY"),
		Issuer:    "grpc-auth",
		Expire:    60,
	}

	storage := repository.NewStorage(db)
	authService := services.NewAuthService(storage)

	s := app.NewServer(jwt, authService)

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, s)

	if err = grpcServer.Serve(lis); err != nil {
		return err
	}

	fmt.Println(storage)
	return nil
}

func setupDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupNet(port string) (net.Listener, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	return lis, nil
}
