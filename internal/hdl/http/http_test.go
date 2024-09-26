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
