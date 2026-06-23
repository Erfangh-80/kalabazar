package handlers

import (
	"net/http"

	appSettlement "stock-service-version-three/internal/application/settlement"
	"stock-service-version-three/internal/interfaces/dto"
)

type CreateSettlementHandler struct {
	useCase *appSettlement.CreateSettlementUseCase
}

func NewCreateSettlementHandler(useCase *appSettlement.CreateSettlementUseCase) *CreateSettlementHandler {
	return &CreateSettlementHandler{useCase: useCase}
}

func (h *CreateSettlementHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSettlementRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appSettlement.CreateSettlementRequest{
		SellerID:   req.SellerID,
		GrossSales: req.GrossSales,
		Commission: req.Commission,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusCreated, dto.CreateSettlementResponse{
		SettlementID: resp.SettlementID,
		GrossSales:   resp.GrossSales,
		Commission:   resp.Commission,
		NetAmount:    resp.NetAmount,
	})
}
