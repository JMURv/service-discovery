package ctrl

import (
	"context"
	"errors"
	"github.com/JMURv/service-discovery/internal/ctrl/mocks"
	"github.com/JMURv/service-discovery/internal/repo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRegister(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	svcRepo := mocks.NewMockServiceDiscoveryRepo(ctrlMock)
	ctrl := New(svcRepo)

	ctx := context.Background()
	name := "test-svc"
	addr := "http://localhost:8080"

	// Test case 1: User exists
	svcRepo.EXPECT().Register(gomock.Any(), name, addr).Return(nil).Times(1)

	err := ctrl.Register(ctx, name, addr)
	assert.Nil(t, err)

	// Test case 2: ErrAlreadyExists
	svcRepo.EXPECT().Register(gomock.Any(), name, addr).Return(repo.ErrAlreadyExists).Times(1)

	err = ctrl.Register(ctx, name, addr)
	assert.IsType(t, repo.ErrAlreadyExists, err)

	// Test case 3: Repo error (other than ErrAlreadyExists)
	var ErrOther = errors.New("other error")
	svcRepo.EXPECT().Register(gomock.Any(), name, addr).Return(ErrOther).Times(1)

	err = ctrl.Register(ctx, name, addr)
	assert.IsType(t, ErrOther, err)
}

func TestDeregister(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	svcRepo := mocks.NewMockServiceDiscoveryRepo(ctrlMock)
	ctrl := New(svcRepo)

	ctx := context.Background()
	name := "test-svc"
	addr := "http://localhost:8080"

	// Test case 1: User exists
	svcRepo.EXPECT().Deregister(gomock.Any(), name, addr).Return(nil).Times(1)

	err := ctrl.Deregister(ctx, name, addr)
	assert.Nil(t, err)

	// Test case 2: ErrNotFound
	svcRepo.EXPECT().Deregister(gomock.Any(), name, addr).Return(repo.ErrNotFound).Times(1)

	err = ctrl.Deregister(ctx, name, addr)
	assert.IsType(t, repo.ErrNotFound, err)

	// Test case 3: Repo error (other than ErrNotFound)
	var ErrOther = errors.New("other error")
	svcRepo.EXPECT().Deregister(gomock.Any(), name, addr).Return(ErrOther).Times(1)

	err = ctrl.Deregister(ctx, name, addr)
	assert.IsType(t, ErrOther, err)
}

func TestFindServiceByName(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	svcRepo := mocks.NewMockServiceDiscoveryRepo(ctrlMock)
	ctrl := New(svcRepo)

	ctx := context.Background()
	name := "test-svc"
	addr := "http://localhost:8080"

	// Test case 1: Success
	svcRepo.EXPECT().FindServiceByName(gomock.Any(), name).Return(addr, nil).Times(1)

	res, err := ctrl.FindServiceByName(ctx, name)
	assert.Equal(t, addr, res)

	// Test case 2: ErrNotFound
	svcRepo.EXPECT().FindServiceByName(gomock.Any(), name).Return("", repo.ErrNotFound).Times(1)

	res, err = ctrl.FindServiceByName(ctx, name)
	assert.Equal(t, "", res)
	assert.IsType(t, repo.ErrNotFound, err)

	// Test case 3: Repo error (other than ErrNotFound)
	var ErrOther = errors.New("other error")
	svcRepo.EXPECT().FindServiceByName(gomock.Any(), name).Return("", ErrOther).Times(1)

	res, err = ctrl.FindServiceByName(ctx, name)
	assert.Equal(t, "", res)
	assert.IsType(t, ErrOther, err)
}

func TestListServices(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	svcRepo := mocks.NewMockServiceDiscoveryRepo(ctrlMock)
	ctrl := New(svcRepo)

	ctx := context.Background()
	expectedRes := []string{"name1", "name2", "name3"}

	// Test case 1: Success
	svcRepo.EXPECT().ListServices(gomock.Any()).Return(expectedRes, nil).Times(1)

	res, err := ctrl.ListServices(ctx)
	assert.Equal(t, expectedRes, res)

	// Test case 2: Repo error
	var ErrOther = errors.New("other error")
	svcRepo.EXPECT().ListServices(gomock.Any()).Return([]string{}, ErrOther).Times(1)

	res, err = ctrl.ListServices(ctx)
	assert.Equal(t, []string{}, res)
	assert.IsType(t, ErrOther, err)
}

func TestListAddrs(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	svcRepo := mocks.NewMockServiceDiscoveryRepo(ctrlMock)
	ctrl := New(svcRepo)

	ctx := context.Background()
	expectedRes := []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:8082"}
	name := "test-svc"

	// Test case 1: Success
	svcRepo.EXPECT().ListAddrs(gomock.Any(), name).Return(expectedRes, nil).Times(1)
	res, err := ctrl.ListAddrs(ctx, name)
	assert.Equal(t, expectedRes, res)

	// Test case 2: ErrNotFound
	svcRepo.EXPECT().ListAddrs(gomock.Any(), name).Return([]string{}, repo.ErrNotFound).Times(1)

	res, err = ctrl.ListAddrs(ctx, name)
	assert.Equal(t, []string{}, res)
	assert.IsType(t, repo.ErrNotFound, err)

	// Test case 3: Repo error (other than ErrNotFound)
	var ErrOther = errors.New("other error")
	svcRepo.EXPECT().ListAddrs(gomock.Any(), name).Return([]string{}, ErrOther).Times(1)

	res, err = ctrl.ListAddrs(ctx, name)
	assert.Equal(t, []string{}, res)
	assert.IsType(t, ErrOther, err)
}
