package handlers

import (
	"net/http"
	"strconv"

	appStore "stock-service-version-three/internal/application/store"
	"stock-service-version-three/internal/interfaces/dto"
)

type ApproveStoreHandler struct {
	useCase *appStore.ApproveStoreUseCase
}

func NewApproveStoreHandler(useCase *appStore.ApproveStoreUseCase) *ApproveStoreHandler {
	return &ApproveStoreHandler{useCase: useCase}
}

func (h *ApproveStoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.PathValue("storeID")
	storeID, err := strconv.ParseInt(storeIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid store ID")
		return
	}

	var req dto.ApproveStoreRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appStore.ApproveStoreRequest{
		StoreID:  storeID,
		Decision: req.Decision,
	}

	resp, err := h.useCase.Execute(r.Context(), appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.ApproveStoreResponse{
		StoreID: resp.StoreID,
		Status:  resp.Status,
	})
}

type AllowCategoryHandler struct {
	useCase *appStore.AllowCategoryUseCase
}

func NewAllowCategoryHandler(useCase *appStore.AllowCategoryUseCase) *AllowCategoryHandler {
	return &AllowCategoryHandler{useCase: useCase}
}

func (h *AllowCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.PathValue("storeID")
	storeID, err := strconv.ParseInt(storeIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid store ID")
		return
	}

	var req dto.AllowCategoryRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appStore.AllowCategoryRequest{
		StoreID:    storeID,
		CategoryID: req.CategoryID,
	}

	resp, err := h.useCase.Execute(r.Context(), appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.AllowCategoryResponse{
		StoreID:    resp.StoreID,
		CategoryID: resp.CategoryID,
	})
}
