package grpc

import (
	"context"
	"errors"
	pb "github.com/JMURv/service-discovery/api/pb"
	"github.com/JMURv/service-discovery/internal/ctrl"
	"github.com/JMURv/service-discovery/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	name := "test-svc"
	addr := "http://localhost:8080"

	// Test case 1: Success
	ctrlRepo.EXPECT().Register(gomock.Any(), name, addr).Return(nil).Times(1)

	_, err := hdl.Register(ctx, &pb.NameAndAddressMsg{Name: name, Address: addr})
	assert.Nil(t, err)

	// Test case 2: ErrAlreadyExists
	ctrlRepo.EXPECT().Register(gomock.Any(), name, addr).Return(ctrl.ErrAlreadyExists).Times(1)

	_, err = hdl.Register(ctx, &pb.NameAndAddressMsg{Name: name, Address: addr})
	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.AlreadyExists)
	assert.Equal(t, s.Message(), ctrl.ErrAlreadyExists.Error())

	// Test case 3: ErrInternalError
	var ErrOther = errors.New("other error")
	ctrlRepo.EXPECT().Register(gomock.Any(), name, addr).Return(ErrOther).Times(1)

	_, err = hdl.Register(ctx, &pb.NameAndAddressMsg{Name: name, Address: addr})
	s, ok = status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.Internal)
	assert.Equal(t, s.Message(), ctrl.ErrInternalError.Error())

	// Test case 4: ErrDecodeRequest
	_, err = hdl.Register(ctx, &pb.NameAndAddressMsg{Name: ""})
	s, ok = status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.InvalidArgument)
	assert.Equal(t, s.Message(), ctrl.ErrDecodeRequest.Error())

	// Test case 5: ErrDecodeRequest
	_, err = hdl.Register(ctx, &pb.NameAndAddressMsg{Address: ""})
	s, ok = status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.InvalidArgument)
	assert.Equal(t, s.Message(), ctrl.ErrDecodeRequest.Error())
}

func TestDeregister(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	name := "test-svc"
	addr := "http://localhost:8080"

	// Test case 1: Success
	ctrlRepo.EXPECT().Deregister(gomock.Any(), name, addr).Return(nil).Times(1)

	_, err := hdl.Deregister(ctx, &pb.NameAndAddressMsg{Name: name, Address: addr})
	assert.Nil(t, err)

	// Test case 2: ErrAlreadyExists
	ctrlRepo.EXPECT().Deregister(gomock.Any(), name, addr).Return(ctrl.ErrNotFound).Times(1)

	_, err = hdl.Deregister(ctx, &pb.NameAndAddressMsg{Name: name, Address: addr})
	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.NotFound)
	assert.Equal(t, s.Message(), ctrl.ErrNotFound.Error())

	// Test case 3: ErrInternalError
	var ErrOther = errors.New("other error")
	ctrlRepo.EXPECT().Deregister(gomock.Any(), name, addr).Return(ErrOther).Times(1)

	_, err = hdl.Deregister(ctx, &pb.NameAndAddressMsg{Name: name, Address: addr})
	s, ok = status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.Internal)
	assert.Equal(t, s.Message(), ctrl.ErrInternalError.Error())

	// Test case 4: ErrDecodeRequest - missing name
	_, err = hdl.Deregister(ctx, &pb.NameAndAddressMsg{Name: ""})
	s, ok = status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.InvalidArgument)
	assert.Equal(t, s.Message(), ctrl.ErrDecodeRequest.Error())

	// Test case 5: ErrDecodeRequest - missing address
	_, err = hdl.Deregister(ctx, &pb.NameAndAddressMsg{Address: ""})
	s, ok = status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.InvalidArgument)
	assert.Equal(t, s.Message(), ctrl.ErrDecodeRequest.Error())
}

func TestFindService(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	name := "test-svc"
	addr := "http://localhost:8080"

	// Test case 1: Success
	ctrlRepo.EXPECT().FindServiceByName(gomock.Any(), name).Return(addr, nil).Times(1)

	_, err := hdl.FindService(ctx, &pb.ServiceNameMsg{Name: name})
	assert.Nil(t, err)

	// Test case 2: ErrNotFound
	ctrlRepo.EXPECT().FindServiceByName(gomock.Any(), name).Return("", ctrl.ErrNotFound).Times(1)

	_, err = hdl.FindService(ctx, &pb.ServiceNameMsg{Name: name})
	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.NotFound)
	assert.Equal(t, s.Message(), ctrl.ErrNotFound.Error())

	// Test case 3: ErrInternalError
	var ErrOther = errors.New("other error")
	ctrlRepo.EXPECT().FindServiceByName(gomock.Any(), name).Return("", ErrOther).Times(1)

	_, err = hdl.FindService(ctx, &pb.ServiceNameMsg{Name: name})
	s, ok = status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.Internal)
	assert.Equal(t, s.Message(), ctrl.ErrInternalError.Error())

	// Test case 4: ErrDecodeRequest - missing name
	_, err = hdl.FindService(ctx, &pb.ServiceNameMsg{Name: ""})
	s, ok = status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.InvalidArgument)
	assert.Equal(t, s.Message(), ctrl.ErrDecodeRequest.Error())
}

func TestListServices(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	names := []string{"name-1", "name-1"}
	expectedRes := &pb.ListNamesMsg{Name: names}

	// Test case 1: Success
	ctrlRepo.EXPECT().ListServices(gomock.Any()).Return(names, nil).Times(1)

	res, err := hdl.ListServices(ctx, &pb.Empty{})
	assert.Nil(t, err)
	assert.Equal(t, expectedRes, res)

	// Test case 2: ErrInternalError
	var ErrOther = errors.New("other error")
	ctrlRepo.EXPECT().ListServices(gomock.Any()).Return([]string{}, ErrOther).Times(1)

	_, err = hdl.ListServices(ctx, &pb.Empty{})
	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.Internal)
	assert.Equal(t, s.Message(), ctrl.ErrInternalError.Error())

}

func TestListAddrs(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	name := "test-svc"
	addrs := []string{"http://localhost:8080", "http://localhost:8081"}
	expectedRes := &pb.ListAddrsMsg{Address: addrs}

	// Test case 1: Success
	ctrlRepo.EXPECT().ListAddrs(gomock.Any(), name).Return(addrs, nil).Times(1)

	res, err := hdl.ListAddrs(ctx, &pb.ServiceNameMsg{Name: name})
	assert.Nil(t, err)
	assert.Equal(t, expectedRes, res)

	// Test case 2: ErrNotFound
	ctrlRepo.EXPECT().ListAddrs(gomock.Any(), name).Return([]string{}, ctrl.ErrNotFound).Times(1)

	_, err = hdl.ListAddrs(ctx, &pb.ServiceNameMsg{Name: name})
	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.NotFound)
	assert.Equal(t, s.Message(), ctrl.ErrNotFound.Error())

	// Test case 3: ErrInternalError
	var ErrOther = errors.New("other error")
	ctrlRepo.EXPECT().ListAddrs(gomock.Any(), name).Return([]string{}, ErrOther).Times(1)

	_, err = hdl.ListAddrs(ctx, &pb.ServiceNameMsg{Name: name})
	s, ok = status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.Internal)
	assert.Equal(t, s.Message(), ctrl.ErrInternalError.Error())

	// Test case 4: ErrDecodeRequest - missing name
	_, err = hdl.ListAddrs(ctx, &pb.ServiceNameMsg{Name: ""})
	s, ok = status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}
	assert.Equal(t, s.Code(), codes.InvalidArgument)
	assert.Equal(t, s.Message(), ctrl.ErrDecodeRequest.Error())
}

func TestStart(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	go hdl.Start(8080)
	time.Sleep(500 * time.Millisecond)

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()
	assert.Nil(t, err)

	if err := hdl.Close(); err != nil {
		t.Fatalf("expected no error while closing, got %v", err)
	}
}
