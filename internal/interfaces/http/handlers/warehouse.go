package handlers

import (
	"net/http"
	"strconv"

	appWarehouse "stock-service-version-three/internal/application/warehouse"
	"stock-service-version-three/internal/interfaces/dto"
)

type CreateWarehouseHandler struct {
	useCase *appWarehouse.CreateWarehouseUseCase
}

func NewCreateWarehouseHandler(useCase *appWarehouse.CreateWarehouseUseCase) *CreateWarehouseHandler {
	return &CreateWarehouseHandler{useCase: useCase}
}

func (h *CreateWarehouseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWarehouseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appWarehouse.CreateWarehouseRequest{
		Name:     req.Name,
		Capacity: req.Capacity,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusCreated, dto.CreateWarehouseResponse{
		WarehouseID: resp.WarehouseID,
		Name:        resp.Name,
	})
}

type LinkWarehouseHandler struct {
	useCase *appWarehouse.LinkWarehouseToStoreUseCase
}

func NewLinkWarehouseHandler(useCase *appWarehouse.LinkWarehouseToStoreUseCase) *LinkWarehouseHandler {
	return &LinkWarehouseHandler{useCase: useCase}
}

func (h *LinkWarehouseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.PathValue("storeID")
	storeID, err := strconv.ParseInt(storeIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid store ID")
		return
	}

	var req dto.LinkWarehouseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appWarehouse.LinkWarehouseRequest{
		StoreID:     storeID,
		WarehouseID: req.WarehouseID,
		Type:        req.Type,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.LinkWarehouseResponse{
		StoreID:     resp.StoreID,
		WarehouseID: resp.WarehouseID,
	})
}
