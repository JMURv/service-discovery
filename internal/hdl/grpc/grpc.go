package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/JMURv/service-discovery/api/pb"
	"github.com/JMURv/service-discovery/internal/ctrl"
	md "github.com/JMURv/service-discovery/pkg/model"
	"github.com/JMURv/service-discovery/pkg/model/mapper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
)

type Ctrl interface {
	ListNames(ctx context.Context) ([]string, error)
	ListAddrsByName(ctx context.Context, name string) ([]string, error)
	ListServices(ctx context.Context) ([]md.Service, error)
	FindServiceByName(ctx context.Context, name string) (string, error)
	Register(ctx context.Context, name, addr string) error
	Deregister(ctx context.Context, name, addr string) error
}

type Handler struct {
	pb.ServiceDiscoveryServer
	srv  *grpc.Server
	ctrl Ctrl
}

func New(ctrl Ctrl) *Handler {
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
		zap.L().Fatal("failed to listen", zap.Error(err))
	}

	if err = h.srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		zap.L().Fatal("failed to serve", zap.Error(err))
	}
}

func (h *Handler) Close() error {
	h.srv.GracefulStop()
	return nil
}

func (h *Handler) ListNames(ctx context.Context, _ *pb.Empty) (*pb.ListNamesMsg, error) {
	res, err := h.ctrl.ListNames(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, ctrl.ErrInternalError.Error())
	}

	return &pb.ListNamesMsg{
		Name: res,
	}, nil
}

func (h *Handler) ListAddrsByName(ctx context.Context, req *pb.ServiceNameMsg) (*pb.ListAddrsMsg, error) {
	if req == nil || req.Name == "" {
		zap.L().Error("failed to decode request")
		return nil, status.Errorf(codes.InvalidArgument, ctrl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.ListAddrsByName(ctx, req.Name)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, ctrl.ErrInternalError.Error())
	}

	return &pb.ListAddrsMsg{
		Address: res,
	}, nil
}

func (h *Handler) ListServices(ctx context.Context, _ *pb.Empty) (*pb.ListServiceMsg, error) {
	res, err := h.ctrl.ListServices(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, ctrl.ErrInternalError.Error())
	}

	return &pb.ListServiceMsg{
		Service: mapper.ListSvcToProto(res),
	}, nil
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

func (h *Handler) FindServiceByName(ctx context.Context, req *pb.ServiceNameMsg) (*pb.ServiceAddressMsg, error) {
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
