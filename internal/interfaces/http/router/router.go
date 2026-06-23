package router

import (
	"net/http"

	handler "stock-service-version-three/internal/interfaces/http/handlers"
)

func NewRouter(
	sellerRegister *handler.RegisterSellerHandler,
	sellerVerifyKYC *handler.VerifyKYCHandler,
	sellerUpdateRank *handler.UpdateRankHandler,
	storeApprove *handler.ApproveStoreHandler,
	storeAllowCategory *handler.AllowCategoryHandler,
	warehouseCreate *handler.CreateWarehouseHandler,
	warehouseLink *handler.LinkWarehouseHandler,
	productCreate *handler.CreateProductHandler,
	productApprove *handler.ApproveProductHandler,
	inventoryCreate *handler.CreateInventoryHandler,
	inventoryUpdatePrice *handler.UpdatePriceHandler,
	inventoryOrderPaid *handler.HandleOrderPaidHandler,
	inventoryOrderDelivered *handler.HandleOrderDeliveredHandler,
	inventoryResetPrice *handler.ResetPriceHandler,
	inventoryRecordRefPrice *handler.RecordReferencePriceHandler,
	campaignCreate *handler.CreateCampaignHandler,
	campaignLink *handler.LinkCampaignHandler,
	campaignApprove *handler.ApproveCampaignHandler,
	campaignActivate *handler.ActivateCampaignHandler,
	campaignEnd *handler.EndCampaignHandler,
	commissionCalc *handler.CalculateCommissionHandler,
	settlementCreate *handler.CreateSettlementHandler,
	payoutExecute *handler.ExecutePayoutHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /api/v1/sellers/register", sellerRegister)
	mux.Handle("POST /api/v1/sellers/verify-kyc", sellerVerifyKYC)
	mux.Handle("PUT /api/v1/sellers/rank", sellerUpdateRank)

	mux.Handle("POST /api/v1/stores/{id}/approve", storeApprove)
	mux.Handle("POST /api/v1/stores/{id}/categories", storeAllowCategory)

	mux.Handle("POST /api/v1/warehouses", warehouseCreate)
	mux.Handle("POST /api/v1/warehouses/link", warehouseLink)

	mux.Handle("POST /api/v1/products", productCreate)
	mux.Handle("POST /api/v1/products/{id}/approve", productApprove)

	mux.Handle("POST /api/v1/inventory", inventoryCreate)
	mux.Handle("PUT /api/v1/inventory/{id}/price", inventoryUpdatePrice)
	mux.Handle("POST /api/v1/inventory/{id}/reserve", inventoryOrderPaid)
	mux.Handle("POST /api/v1/inventory/{id}/finalize", inventoryOrderDelivered)
	mux.Handle("POST /api/v1/inventory/{id}/reset-price", inventoryResetPrice)
	mux.Handle("POST /api/v1/reference-prices", inventoryRecordRefPrice)

	mux.Handle("POST /api/v1/campaigns", campaignCreate)
	mux.Handle("POST /api/v1/campaigns/{id}/link", campaignLink)
	mux.Handle("POST /api/v1/campaigns/{id}/approve", campaignApprove)
	mux.Handle("POST /api/v1/campaigns/{id}/activate", campaignActivate)
	mux.Handle("POST /api/v1/campaigns/{id}/end", campaignEnd)

	mux.Handle("POST /api/v1/commissions", commissionCalc)

	mux.Handle("POST /api/v1/settlements", settlementCreate)

	mux.Handle("POST /api/v1/payouts", payoutExecute)

	return mux
}
