package rpchandler

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ubiqueworks/go-interface-usage/domain"
	"github.com/ubiqueworks/go-interface-usage/pb"
	"github.com/ubiqueworks/go-interface-usage/service"
	"google.golang.org/grpc"
)

func New(s service.Service, bindAddr string, bindPort int) *RpcHandler {
	h := &RpcHandler{
		addr:    fmt.Sprintf("%s:%d", bindAddr, bindPort),
		service: s,
	}
	return h
}

type RpcHandler struct {
	addr    string
	service service.Service
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
