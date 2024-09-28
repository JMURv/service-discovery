package http

import (
	"bytes"
	"context"
	"errors"
	"github.com/JMURv/service-discovery/internal/ctrl"
	"github.com/JMURv/service-discovery/mocks"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
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

	payload, _ := json.Marshal(map[string]string{"name": name, "address": addr})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	hdl.register(w, req)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	// Test case 2: ErrAlreadyExists
	ctrlRepo.EXPECT().Register(gomock.Any(), name, addr).Return(ctrl.ErrAlreadyExists).Times(1)

	payload, _ = json.Marshal(map[string]string{"name": name, "address": addr})
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.register(w, req)
	assert.Equal(t, http.StatusConflict, w.Result().StatusCode)

	// Test case 3: ErrInternalError
	var ErrOther = errors.New("other error")
	ctrlRepo.EXPECT().Register(gomock.Any(), name, addr).Return(ErrOther).Times(1)

	payload, _ = json.Marshal(map[string]string{"name": name, "address": addr})
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.register(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	// Test case 4: ErrDecodeRequest
	payload, _ = json.Marshal(map[string]string{"name": name})
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.register(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Test case 5: ErrDecodeRequest
	payload, _ = json.Marshal(map[string]string{"address": addr})
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.register(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Test case 6: Invalid JSOM
	payload, _ = json.Marshal(map[string]any{"address": 123})
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.register(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
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

	payload, _ := json.Marshal(map[string]string{"name": name, "address": addr})
	req := httptest.NewRequest(http.MethodPost, "/deregister", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	hdl.deregister(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Test case 2: ErrNotFound
	ctrlRepo.EXPECT().Deregister(gomock.Any(), name, addr).Return(ctrl.ErrNotFound).Times(1)

	payload, _ = json.Marshal(map[string]string{"name": name, "address": addr})
	req = httptest.NewRequest(http.MethodPost, "/deregister", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.deregister(w, req)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

	// Test case 3: ErrInternalError
	var ErrOther = errors.New("other error")
	ctrlRepo.EXPECT().Deregister(gomock.Any(), name, addr).Return(ErrOther).Times(1)

	payload, _ = json.Marshal(map[string]string{"name": name, "address": addr})
	req = httptest.NewRequest(http.MethodPost, "/deregister", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.deregister(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	// Test case 4: ErrDecodeRequest
	payload, _ = json.Marshal(map[string]string{"name": name})
	req = httptest.NewRequest(http.MethodPost, "/deregister", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.deregister(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Test case 5: ErrDecodeRequest
	payload, _ = json.Marshal(map[string]string{"address": addr})
	req = httptest.NewRequest(http.MethodPost, "/deregister", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.deregister(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Test case 6: Invalid JSOM
	payload, _ = json.Marshal(map[string]any{"address": 123})
	req = httptest.NewRequest(http.MethodPost, "/deregister", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.deregister(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestFindSvc(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	name := "test-svc"
	addr := "http://localhost:8080"

	// Test case 1: Success
	ctrlRepo.EXPECT().FindServiceByName(gomock.Any(), name).Return(addr, nil).Times(1)

	payload, _ := json.Marshal(map[string]string{"name": name})
	req := httptest.NewRequest(http.MethodPost, "/find", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	hdl.find(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Test case 2: ErrNotFound
	ctrlRepo.EXPECT().FindServiceByName(gomock.Any(), name).Return("", ctrl.ErrNotFound).Times(1)

	payload, _ = json.Marshal(map[string]string{"name": name})
	req = httptest.NewRequest(http.MethodPost, "/find", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.find(w, req)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

	// Test case 3: ErrInternalError
	var ErrOther = errors.New("other error")
	ctrlRepo.EXPECT().FindServiceByName(gomock.Any(), name).Return("", ErrOther).Times(1)

	payload, _ = json.Marshal(map[string]string{"name": name})
	req = httptest.NewRequest(http.MethodPost, "/find", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.find(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	// Test case 4: ErrDecodeRequest
	payload, _ = json.Marshal(map[string]string{"name": ""})
	req = httptest.NewRequest(http.MethodPost, "/find", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.find(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Test case 5: ErrDecodeRequest
	payload, _ = json.Marshal(map[string]string{"address": addr})
	req = httptest.NewRequest(http.MethodPost, "/find", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.find(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Test case 6: Invalid JSOM
	payload, _ = json.Marshal(map[string]any{"address": 123})
	req = httptest.NewRequest(http.MethodPost, "/find", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.find(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestListAddrs(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	name := "test-svc"
	expRes := []string{"http://localhost:8080", "http://localhost:8081"}

	// Test case 1: Success
	ctrlRepo.EXPECT().ListAddrs(gomock.Any(), name).Return(expRes, nil).Times(1)

	payload, _ := json.Marshal(map[string]string{"name": name})
	req := httptest.NewRequest(http.MethodPost, "/list-addrs", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	hdl.listAddrs(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Test case 2: ErrAlreadyExists
	ctrlRepo.EXPECT().ListAddrs(gomock.Any(), name).Return([]string{}, ctrl.ErrAlreadyExists).Times(1)

	payload, _ = json.Marshal(map[string]string{"name": name})
	req = httptest.NewRequest(http.MethodPost, "/list-addrs", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.listAddrs(w, req)
	assert.Equal(t, http.StatusConflict, w.Result().StatusCode)

	// Test case 3: ErrInternalError
	var ErrOther = errors.New("other error")
	ctrlRepo.EXPECT().ListAddrs(gomock.Any(), name).Return([]string{}, ErrOther).Times(1)

	payload, _ = json.Marshal(map[string]string{"name": name})
	req = httptest.NewRequest(http.MethodPost, "/list-addrs", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.listAddrs(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	// Test case 4: ErrDecodeRequest
	payload, _ = json.Marshal(map[string]string{"name": ""})
	req = httptest.NewRequest(http.MethodPost, "/list-addrs", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.listAddrs(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Test case 5: ErrDecodeRequest
	payload, _ = json.Marshal(map[string]string{"address": "addr"})
	req = httptest.NewRequest(http.MethodPost, "/list-addrs", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.listAddrs(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Test case 6: Invalid JSOM
	payload, _ = json.Marshal(map[string]any{"address": 123})
	req = httptest.NewRequest(http.MethodPost, "/list-addrs", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.listAddrs(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestListSvcs(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	ctx := context.Background()
	expRes := []string{"http://localhost:8080", "http://localhost:8081"}

	// Test case 1: Success
	ctrlRepo.EXPECT().ListServices(gomock.Any()).Return(expRes, nil).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/list-svcs", nil)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	hdl.listSvcs(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Test case 2: ErrAlreadyExists
	ctrlRepo.EXPECT().ListServices(gomock.Any()).Return([]string{}, ctrl.ErrAlreadyExists).Times(1)

	req = httptest.NewRequest(http.MethodGet, "/list-svcs", nil)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.listSvcs(w, req)
	assert.Equal(t, http.StatusConflict, w.Result().StatusCode)

	// Test case 3: ErrInternalError
	var ErrOther = errors.New("other error")
	ctrlRepo.EXPECT().ListServices(gomock.Any()).Return([]string{}, ErrOther).Times(1)

	req = httptest.NewRequest(http.MethodGet, "/list-svcs", nil)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	hdl.listSvcs(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestStart(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	go hdl.Start(8080)
	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/health-check")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %v", resp.StatusCode)
	}
}

func TestClose(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	ctrlRepo := mocks.NewMockCtrl(ctrlMock)
	hdl := New(ctrlRepo)

	go hdl.Start(8080)
	time.Sleep(500 * time.Millisecond)

	if err := hdl.Close(); err != nil {
		t.Fatalf("Expected no error while closing, got %v", err)
	}

	resp, err := http.Get("http://localhost:8080/health-check")
	if err == nil {
		defer resp.Body.Close()
	}
}
