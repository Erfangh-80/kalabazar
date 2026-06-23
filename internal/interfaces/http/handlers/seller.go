package handlers

import (
	"net/http"
	"strconv"

	appSeller "stock-service-version-three/internal/application/seller"
	"stock-service-version-three/internal/interfaces/dto"
)

type RegisterSellerHandler struct {
	useCase *appSeller.RegisterSellerUseCase
}

func NewRegisterSellerHandler(useCase *appSeller.RegisterSellerUseCase) *RegisterSellerHandler {
	return &RegisterSellerHandler{useCase: useCase}
}

func (h *RegisterSellerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterSellerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appSeller.RegisterSellerRequest{
		UserID:    req.UserID,
		StoreName: req.StoreName,
		Phone:     req.Phone,
	}

	resp, err := h.useCase.Execute(r.Context(), appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusCreated, dto.RegisterSellerResponse{
		SellerID: resp.SellerID,
		StoreID:  resp.StoreID,
	})
}

type VerifyKYCHandler struct {
	useCase *appSeller.VerifyKYCUseCase
}

func NewVerifyKYCHandler(useCase *appSeller.VerifyKYCUseCase) *VerifyKYCHandler {
	return &VerifyKYCHandler{useCase: useCase}
}

func (h *VerifyKYCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyKYCRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appSeller.VerifyKYCRequest{
		SellerID:  req.SellerID,
		KYCStatus: req.KYCStatus,
	}

	resp, err := h.useCase.Execute(r.Context(), appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.VerifyKYCResponse{
		SellerID: resp.SellerID,
		Status:   resp.Status,
	})
}

type UpdateRankHandler struct {
	useCase *appSeller.UpdateRankUseCase
}

func NewUpdateRankHandler(useCase *appSeller.UpdateRankUseCase) *UpdateRankHandler {
	return &UpdateRankHandler{useCase: useCase}
}

func (h *UpdateRankHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sellerIDStr := r.PathValue("sellerID")
	sellerID, err := strconv.ParseInt(sellerIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid seller ID")
		return
	}

	var req dto.UpdateRankRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appSeller.UpdateRankRequest{
		SellerID: sellerID,
		Score:    req.Score,
		Rank:     req.Rank,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.UpdateRankResponse{
		SellerID: resp.SellerID,
		Score:    resp.Score,
		Rank:     resp.Rank,
	})
}
