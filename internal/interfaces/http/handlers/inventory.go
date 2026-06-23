package handlers

import (
	"net/http"
	"strconv"

	appInventory "stock-service-version-three/internal/application/inventory"
	"stock-service-version-three/internal/interfaces/dto"
)

type CreateInventoryHandler struct {
	useCase *appInventory.CreateInventoryUseCase
}

func NewCreateInventoryHandler(useCase *appInventory.CreateInventoryUseCase) *CreateInventoryHandler {
	return &CreateInventoryHandler{useCase: useCase}
}

func (h *CreateInventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInventoryRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appInventory.CreateInventoryRequest{
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		BasePrice:   req.BasePrice,
		Stock:       req.Stock,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusCreated, dto.CreateInventoryResponse{
		InventoryID:    resp.InventoryID,
		AvailableStock: resp.AvailableStock,
		FinalPrice:     resp.FinalPrice,
	})
}

type UpdatePriceHandler struct {
	useCase *appInventory.UpdatePriceUseCase
}

func NewUpdatePriceHandler(useCase *appInventory.UpdatePriceUseCase) *UpdatePriceHandler {
	return &UpdatePriceHandler{useCase: useCase}
}

func (h *UpdatePriceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	inventoryIDStr := r.PathValue("inventoryID")
	inventoryID, err := strconv.ParseInt(inventoryIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid inventory ID")
		return
	}

	var req dto.UpdatePriceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appInventory.UpdatePriceRequest{
		InventoryID:       inventoryID,
		DiscountPercentage: req.DiscountPercentage,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.UpdatePriceResponse{
		InventoryID: resp.InventoryID,
		OldPrice:    resp.OldPrice,
		NewPrice:    resp.NewPrice,
	})
}

type HandleOrderPaidHandler struct {
	useCase *appInventory.HandleOrderPaidUseCase
}

func NewHandleOrderPaidHandler(useCase *appInventory.HandleOrderPaidUseCase) *HandleOrderPaidHandler {
	return &HandleOrderPaidHandler{useCase: useCase}
}

func (h *HandleOrderPaidHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	inventoryIDStr := r.PathValue("inventoryID")
	inventoryID, err := strconv.ParseInt(inventoryIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid inventory ID")
		return
	}

	var req dto.HandleOrderPaidRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appInventory.HandleOrderPaidRequest{
		InventoryID: inventoryID,
		Quantity:    req.Quantity,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.HandleOrderPaidResponse{
		AvailableStock: resp.AvailableStock,
		ReservedStock:  resp.ReservedStock,
	})
}

type HandleOrderDeliveredHandler struct {
	useCase *appInventory.HandleOrderDeliveredUseCase
}

func NewHandleOrderDeliveredHandler(useCase *appInventory.HandleOrderDeliveredUseCase) *HandleOrderDeliveredHandler {
	return &HandleOrderDeliveredHandler{useCase: useCase}
}

func (h *HandleOrderDeliveredHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	inventoryIDStr := r.PathValue("inventoryID")
	inventoryID, err := strconv.ParseInt(inventoryIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid inventory ID")
		return
	}

	var req dto.HandleOrderDeliveredRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appInventory.HandleOrderDeliveredRequest{
		InventoryID: inventoryID,
		Quantity:    req.Quantity,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.HandleOrderDeliveredResponse{
		StockOut: resp.StockOut,
	})
}

type ResetPriceHandler struct {
	useCase *appInventory.ResetPriceUseCase
}

func NewResetPriceHandler(useCase *appInventory.ResetPriceUseCase) *ResetPriceHandler {
	return &ResetPriceHandler{useCase: useCase}
}

func (h *ResetPriceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	inventoryIDStr := r.PathValue("inventoryID")
	inventoryID, err := strconv.ParseInt(inventoryIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid inventory ID")
		return
	}

	appReq := appInventory.ResetPriceRequest{
		InventoryID: inventoryID,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.ResetPriceResponse{
		InventoryID: resp.InventoryID,
		FinalPrice:  resp.FinalPrice,
	})
}

type RecordReferencePriceHandler struct {
	useCase *appInventory.RecordReferencePriceUseCase
}

func NewRecordReferencePriceHandler(useCase *appInventory.RecordReferencePriceUseCase) *RecordReferencePriceHandler {
	return &RecordReferencePriceHandler{useCase: useCase}
}

func (h *RecordReferencePriceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("productID")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	var req dto.RecordReferencePriceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appInventory.RecordReferencePriceRequest{
		ProductID: productID,
		Price:     req.Price,
		Source:    req.Source,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.RecordReferencePriceResponse{
		ProductID: resp.ProductID,
		Price:     resp.Price,
	})
}
