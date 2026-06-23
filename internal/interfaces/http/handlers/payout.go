package handlers

import (
	"net/http"

	appPayout "stock-service-version-three/internal/application/payout"
	"stock-service-version-three/internal/interfaces/dto"
)

type ExecutePayoutHandler struct {
	useCase *appPayout.ExecutePayoutUseCase
}

func NewExecutePayoutHandler(useCase *appPayout.ExecutePayoutUseCase) *ExecutePayoutHandler {
	return &ExecutePayoutHandler{useCase: useCase}
}

func (h *ExecutePayoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.ExecutePayoutRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appPayout.ExecutePayoutRequest{
		SellerID: req.SellerID,
		Amount:   req.Amount,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusCreated, dto.ExecutePayoutResponse{
		PayoutID: resp.PayoutID,
		Amount:   resp.Amount,
		Status:   string(resp.Status),
	})
}
