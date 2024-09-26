package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/JMURv/service-discovery/api/pb"
	"github.com/JMURv/service-discovery/internal/ctrl"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type Handler struct {
	pb.ServiceDiscoveryServer
	srv  *grpc.Server
	ctrl *ctrl.Controller
}

func New(ctrl *ctrl.Controller) *Handler {
	srv := grpc.NewServer()
	reflection.Register(srv)
	return &Handler{
		ctrl: ctrl,
		srv:  srv,
	}
}

func (h *Handler) Start(port int) {
	pb.RegisterServiceDiscoveryServer(h.srv, h)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Fatal(h.srv.Serve(lis))
}

func (h *Handler) Close() error {
	h.srv.GracefulStop()
	return nil
}

func (h *Handler) Register(ctx context.Context, req *pb.NameAndAddressMsg) (*pb.Empty, error) {
	if req == nil || req.Name == "" || req.Address == "" {
		zap.L().Error("failed to decode request")
		return nil, status.Errorf(codes.InvalidArgument, ctrl.ErrDecodeRequest.Error())
	}

	err := h.ctrl.Register(ctx, req.Name, req.Address)
	if err != nil && errors.Is(err, ctrl.ErrAlreadyExists) {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, ctrl.ErrInternalError.Error())
	}

	return &pb.Empty{}, nil
}

func (h *Handler) Deregister(ctx context.Context, req *pb.NameAndAddressMsg) (*pb.Empty, error) {
	if req == nil || req.Name == "" || req.Address == "" {
		zap.L().Error("failed to decode request")
		return nil, status.Errorf(codes.InvalidArgument, ctrl.ErrDecodeRequest.Error())
	}

	err := h.ctrl.Deregister(ctx, req.Name, req.Address)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, ctrl.ErrInternalError.Error())
	}

	return &pb.Empty{}, nil
}

func (h *Handler) FindService(ctx context.Context, req *pb.ServiceNameMsg) (*pb.ServiceAddressMsg, error) {
	if req == nil || req.Name == "" {
		zap.L().Error("failed to decode request")
		return nil, status.Errorf(codes.InvalidArgument, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.FindServiceByName(ctx, req.Name)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, ctrl.ErrInternalError.Error())
	}

	return &pb.ServiceAddressMsg{
		Address: res,
	}, nil
}

func (h *Handler) ListServices(ctx context.Context, req *pb.Empty) (*pb.ListNamesMsg, error) {
	res, err := h.ctrl.ListServices(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, ctrl.ErrInternalError.Error())
	}

	return &pb.ListNamesMsg{
		Name: res,
	}, nil
}

func (h *Handler) ListAddrs(ctx context.Context, req *pb.ServiceNameMsg) (*pb.ListAddrsMsg, error) {
	if req == nil || req.Name == "" {
		zap.L().Error("failed to decode request")
		return nil, status.Errorf(codes.InvalidArgument, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.ListAddrs(ctx, req.Name)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, ctrl.ErrInternalError.Error())
	}

	return &pb.ListAddrsMsg{
		Address: res,
	}, nil
}
