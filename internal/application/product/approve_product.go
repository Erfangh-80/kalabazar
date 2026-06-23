package product

import (
	"errors"

	domain "stock-service-version-three/internal/domain/product"
)

type ApproveProductUseCase struct {
	repo domain.ProductRepository
}

func NewApproveProductUseCase(repo domain.ProductRepository) *ApproveProductUseCase {
	return &ApproveProductUseCase{repo: repo}
}

func (uc *ApproveProductUseCase) Execute(req ApproveProductRequest) (*ApproveProductResponse, error) {
	p, err := uc.repo.FindByID(req.ProductID)
	if err != nil {
		return nil, err
	}

	switch req.Decision {
	case "approved":
		if err := p.Approve(); err != nil {
			return nil, err
		}
	case "rejected":
		if err := p.Reject(); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid decision: must be 'approved' or 'rejected'")
	}

	if err := uc.repo.Update(p); err != nil {
		return nil, err
	}

	return &ApproveProductResponse{
		ProductID: p.ID,
		Status:    string(p.Status),
	}, nil
}
