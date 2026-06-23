package handlers

import (
	"net/http"
	"strconv"
	"time"

	appCampaign "stock-service-version-three/internal/application/campaign"
	"stock-service-version-three/internal/interfaces/dto"
)

type CreateCampaignHandler struct {
	useCase *appCampaign.CreateCampaignUseCase
}

func NewCreateCampaignHandler(useCase *appCampaign.CreateCampaignUseCase) *CreateCampaignHandler {
	return &CreateCampaignHandler{useCase: useCase}
}

func (h *CreateCampaignHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCampaignRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid start_at: must be RFC3339")
		return
	}
	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid end_at: must be RFC3339")
		return
	}

	appReq := appCampaign.CreateCampaignRequest{
		Title:        req.Title,
		DiscountType: req.DiscountType,
		Value:        req.Value,
		StartAt:      startAt,
		EndAt:        endAt,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusCreated, dto.CreateCampaignResponse{
		CampaignID:     resp.CampaignID,
		Status:         resp.Status,
		ApprovalStatus: resp.ApprovalStatus,
	})
}

type LinkCampaignHandler struct {
	useCase *appCampaign.LinkCampaignUseCase
}

func NewLinkCampaignHandler(useCase *appCampaign.LinkCampaignUseCase) *LinkCampaignHandler {
	return &LinkCampaignHandler{useCase: useCase}
}

func (h *LinkCampaignHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	campaignIDStr := r.PathValue("campaignID")
	campaignID, err := strconv.ParseInt(campaignIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid campaign ID")
		return
	}

	var req dto.LinkCampaignRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appCampaign.LinkCampaignRequest{
		CampaignID:  campaignID,
		InventoryID: req.InventoryID,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.LinkCampaignResponse{
		CampaignID:  resp.CampaignID,
		InventoryID: resp.InventoryID,
	})
}

type ApproveCampaignHandler struct {
	useCase *appCampaign.ApproveCampaignUseCase
}

func NewApproveCampaignHandler(useCase *appCampaign.ApproveCampaignUseCase) *ApproveCampaignHandler {
	return &ApproveCampaignHandler{useCase: useCase}
}

func (h *ApproveCampaignHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	campaignIDStr := r.PathValue("campaignID")
	campaignID, err := strconv.ParseInt(campaignIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid campaign ID")
		return
	}

	var req dto.ApproveCampaignRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	appReq := appCampaign.ApproveCampaignRequest{
		CampaignID: campaignID,
		Decision:   req.Decision,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.ApproveCampaignResponse{
		CampaignID:     resp.CampaignID,
		ApprovalStatus: resp.ApprovalStatus,
	})
}

type ActivateCampaignHandler struct {
	useCase *appCampaign.ActivateCampaignUseCase
}

func NewActivateCampaignHandler(useCase *appCampaign.ActivateCampaignUseCase) *ActivateCampaignHandler {
	return &ActivateCampaignHandler{useCase: useCase}
}

func (h *ActivateCampaignHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	campaignIDStr := r.PathValue("campaignID")
	campaignID, err := strconv.ParseInt(campaignIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid campaign ID")
		return
	}

	var req dto.ActivateCampaignRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	now, err := time.Parse(time.RFC3339, req.Now)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid now: must be RFC3339")
		return
	}

	appReq := appCampaign.ActivateCampaignRequest{
		CampaignID: campaignID,
		Now:        now,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.ActivateCampaignResponse{
		CampaignID: resp.CampaignID,
		Status:     resp.Status,
	})
}

type EndCampaignHandler struct {
	useCase *appCampaign.EndCampaignUseCase
}

func NewEndCampaignHandler(useCase *appCampaign.EndCampaignUseCase) *EndCampaignHandler {
	return &EndCampaignHandler{useCase: useCase}
}

func (h *EndCampaignHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	campaignIDStr := r.PathValue("campaignID")
	campaignID, err := strconv.ParseInt(campaignIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid campaign ID")
		return
	}

	var req dto.EndCampaignRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	now, err := time.Parse(time.RFC3339, req.Now)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid now: must be RFC3339")
		return
	}

	appReq := appCampaign.EndCampaignRequest{
		CampaignID: campaignID,
		Now:        now,
	}

	resp, err := h.useCase.Execute(appReq)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	encodeJSON(w, http.StatusOK, dto.EndCampaignResponse{
		CampaignID: resp.CampaignID,
		Status:     resp.Status,
	})
}
