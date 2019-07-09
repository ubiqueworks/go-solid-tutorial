//go:generate mockgen -package rpchandler -source=rpchandler.go -destination rpchandler_mock.go

package rpchandler

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/caarlos0/env"
	"github.com/ubiqueworks/go-solid-tutorial/domain"
	"github.com/ubiqueworks/go-solid-tutorial/pb"
	"google.golang.org/grpc"
)

type config struct {
	BindAddr string `env:"BIND_ADDR" envDefault:"0.0.0.0" required:"true"`
	Port     int    `env:"RPC_PORT" envDefault:"9000" required:"true"`
}

func (c config) Validate() error {
	if c.BindAddr == "" {
		return fmt.Errorf("missing required HTTP_ADDR")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid RPC_PORT: %d", c.Port)
	}

	return nil
}

type Service interface {
	CreateUser(u *domain.User) error
	DeleteUser(userId string) error
}

func New(s Service) (*RpcHandler, error) {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	h := &RpcHandler{
		addr:    fmt.Sprintf("%s:%d", cfg.BindAddr, cfg.Port),
		service: s,
	}
	return h, nil
}

type RpcHandler struct {
	addr    string
	service Service
}

func (h RpcHandler) Bootstrap() (shutdownFn func() error, err error) {
	log.Printf("starting RPC server on %s", h.addr)

	server := grpc.NewServer()
	pb.RegisterUserAdminServer(server, &h)

	listener, err := net.Listen("tcp", h.addr)
	if err != nil {
		return
	}

	shutdownFn = func() error {
		defer func() {
			log.Println("RPC server shutdown")
		}()

		log.Println("shutting down RPC server...")
		server.GracefulStop()
		return nil
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Printf("RPC server error: %v", err)
		}
	}()

	return
}

func (h RpcHandler) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	user := &domain.User{
		Name: r.Name,
		Password: r.Password,
	}
	if err := h.service.CreateUser(user); err != nil {
		return nil, err
	}
	return &pb.CreateUserReply{
		User: &pb.User{
			Id:   user.ID,
			Name: user.Name,
		},
	}, nil
}

func (h RpcHandler) DeleteUser(ctx context.Context, r *pb.DeleteUserRequest) (*pb.EmptyReply, error) {
	if err := h.service.DeleteUser(r.Id); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
