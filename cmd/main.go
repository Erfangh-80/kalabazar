package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	appCampaign "stock-service-version-three/internal/application/campaign"
	appCommission "stock-service-version-three/internal/application/commission"
	appInventory "stock-service-version-three/internal/application/inventory"
	appPayout "stock-service-version-three/internal/application/payout"
	appProduct "stock-service-version-three/internal/application/product"
	appSeller "stock-service-version-three/internal/application/seller"
	appSettlement "stock-service-version-three/internal/application/settlement"
	appStore "stock-service-version-three/internal/application/store"
	appWarehouse "stock-service-version-three/internal/application/warehouse"

	handler "stock-service-version-three/internal/interfaces/http/handlers"
	"stock-service-version-three/internal/interfaces/http/router"
)

func main() {
	fmt.Println("Vendor Service v3 - Bootstrap")
	fmt.Println("Layer status: Domain \u2705 | Application \u2705 | Interfaces (HTTP) \u2705")
	fmt.Println("Infrastructure (repository implementations): \u23f3 Pending")
	fmt.Println()
	fmt.Println("The service is ready for repository implementation.")
	fmt.Println("Run 'go build ./...' to verify compilation.")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()
	fmt.Println("\nShutting down...")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Bye.")
	os.Exit(0)
}

var _ = router.NewRouter
var _ = handler.RegisterSellerHandler{}
var _ = appSeller.RegisterSellerUseCase{}
var _ = appStore.ApproveStoreUseCase{}
var _ = appWarehouse.CreateWarehouseUseCase{}
var _ = appProduct.CreateProductUseCase{}
var _ = appInventory.CreateInventoryUseCase{}
var _ = appCampaign.CreateCampaignUseCase{}
var _ = appCommission.CalculateCommissionUseCase{}
var _ = appSettlement.CreateSettlementUseCase{}
var _ = appPayout.ExecutePayoutUseCase{}
