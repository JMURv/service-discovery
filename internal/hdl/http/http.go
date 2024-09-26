package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/JMURv/service-discovery/internal/ctrl"
	"github.com/JMURv/service-discovery/internal/validation"
	md "github.com/JMURv/service-discovery/pkg/model"
	utils "github.com/JMURv/service-discovery/pkg/utils/http"
	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Handler struct {
	srv  *http.Server
	ctrl *ctrl.Controller
}

func New(ctrl *ctrl.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) Start(port int) {
	r := mux.NewRouter()
	r.HandleFunc("/health-check", h.healthCheck).Methods(http.MethodGet)
	r.HandleFunc("/register", h.register).Methods(http.MethodPost)
	r.HandleFunc("/deregister", h.deregister).Methods(http.MethodPost)
	r.HandleFunc("/find", h.find).Methods(http.MethodPost)

	r.HandleFunc("/list-svcs", h.listSvcs).Methods(http.MethodGet)
	r.HandleFunc("/list-addrs", h.listAddrs).Methods(http.MethodPost)

	h.srv = &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	err := h.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		zap.L().Debug("Server error", zap.Error(err))
	}
}

func (h *Handler) Close() error {
	if err := h.srv.Shutdown(context.Background()); err != nil {
		return err
	}
	return nil
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	utils.SuccessResponse(w, http.StatusOK, "OK")
}

func (h *Handler) listSvcs(w http.ResponseWriter, r *http.Request) {
	svcs, err := h.ctrl.ListServices(r.Context())
	if err != nil && errors.Is(err, ctrl.ErrAlreadyExists) {
		utils.ErrResponse(w, http.StatusConflict, err)
		return
	} else if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, svcs)
}

func (h *Handler) listAddrs(w http.ResponseWriter, r *http.Request) {
	req := &md.Service{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		zap.L().Debug("failed to decode request", zap.Error(err))
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if req.Name == "" {
		zap.L().Debug("failed to decode request", zap.Error(validation.ErrMissingName))
		utils.ErrResponse(w, http.StatusBadRequest, validation.ErrMissingName)
		return
	}

	svcs, err := h.ctrl.ListAddrs(r.Context(), req.Name)
	if err != nil && errors.Is(err, ctrl.ErrAlreadyExists) {
		utils.ErrResponse(w, http.StatusConflict, err)
		return
	} else if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, svcs)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	req := &md.Service{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		zap.L().Debug("failed to decode request", zap.Error(err))
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if req.Name == "" {
		zap.L().Debug("failed to decode request", zap.Error(validation.ErrMissingName))
		utils.ErrResponse(w, http.StatusBadRequest, validation.ErrMissingName)
		return
	} else if req.Address == "" {
		zap.L().Debug("failed to decode request", zap.Error(validation.ErrMissingAddress))
		utils.ErrResponse(w, http.StatusBadRequest, validation.ErrMissingAddress)
		return
	}

	err := h.ctrl.Register(r.Context(), req.Name, req.Address)
	if err != nil && errors.Is(err, ctrl.ErrAlreadyExists) {
		utils.ErrResponse(w, http.StatusConflict, err)
		return
	} else if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "OK")
}

func (h *Handler) deregister(w http.ResponseWriter, r *http.Request) {
	req := &md.Service{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		zap.L().Debug("failed to decode request", zap.Error(err))
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if req.Name == "" {
		zap.L().Debug("failed to decode request", zap.Error(validation.ErrMissingName))
		utils.ErrResponse(w, http.StatusBadRequest, validation.ErrMissingName)
		return
	} else if req.Address == "" {
		zap.L().Debug("failed to decode request", zap.Error(validation.ErrMissingAddress))
		utils.ErrResponse(w, http.StatusBadRequest, validation.ErrMissingAddress)
		return
	}

	err := h.ctrl.Deregister(r.Context(), req.Name, req.Address)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		utils.ErrResponse(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "OK")
}

func (h *Handler) find(w http.ResponseWriter, r *http.Request) {
	req := &md.Service{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		zap.L().Debug("failed to decode request", zap.Error(err))
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if req.Name == "" {
		zap.L().Debug("failed to decode request", zap.Error(validation.ErrMissingName))
		utils.ErrResponse(w, http.StatusBadRequest, validation.ErrMissingName)
		return
	}

	res, err := h.ctrl.FindServiceByName(r.Context(), req.Name)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		utils.ErrResponse(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, res)
}
