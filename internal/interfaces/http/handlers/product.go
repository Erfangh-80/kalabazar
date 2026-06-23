package handlers

import (
	"net/http"
	"strconv"

	appProduct "stock-service-version-three/internal/application/product"
	"stock-service-version-three/internal/interfaces/dto"
)

type CreateProductHandler struct {
	useCase *appProduct.CreateProductUseCase
}

func NewCreateProductHandler(useCase *appProduct.CreateProductUseCase) *CreateProductHandler {
	return &CreateProductHandler{useCase: useCase}
}

func (h *CreateProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appProduct.CreateProductRequest{
		StoreID:    req.StoreID,
		Title:      req.Title,
		CategoryID: req.CategoryID,
		Brand:      req.Brand,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusCreated, dto.CreateProductResponse{
		ProductID: resp.ProductID,
		Status:    resp.Status,
	})
}

type ApproveProductHandler struct {
	useCase *appProduct.ApproveProductUseCase
}

func NewApproveProductHandler(useCase *appProduct.ApproveProductUseCase) *ApproveProductHandler {
	return &ApproveProductHandler{useCase: useCase}
}

func (h *ApproveProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("productID")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	var req dto.ApproveProductRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appProduct.ApproveProductRequest{
		ProductID: productID,
		Decision:  req.Decision,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.ApproveProductResponse{
		ProductID: resp.ProductID,
		Status:    resp.Status,
	})
}
