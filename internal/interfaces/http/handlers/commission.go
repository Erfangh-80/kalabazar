package handlers

import (
	"net/http"

	appCommission "stock-service-version-three/internal/application/commission"
	"stock-service-version-three/internal/interfaces/dto"
)

type CalculateCommissionHandler struct {
	useCase *appCommission.CalculateCommissionUseCase
}

func NewCalculateCommissionHandler(useCase *appCommission.CalculateCommissionUseCase) *CalculateCommissionHandler {
	return &CalculateCommissionHandler{useCase: useCase}
}

func (h *CalculateCommissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.CalculateCommissionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appCommission.CalculateCommissionRequest{
		SellerID:    req.SellerID,
		SalesAmount: req.SalesAmount,
		Rate:        req.Rate,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusCreated, dto.CalculateCommissionResponse{
		CommissionID: resp.CommissionID,
		Amount:       resp.Amount,
		SalesAmount:  resp.SalesAmount,
	})
}
