package inventory

type RecordReferencePriceUseCase struct{}

func NewRecordReferencePriceUseCase() *RecordReferencePriceUseCase {
	return &RecordReferencePriceUseCase{}
}

func (uc *RecordReferencePriceUseCase) Execute(req RecordReferencePriceRequest) (*RecordReferencePriceResponse, error) {
	if req.Price <= 0 {
		return nil, ErrInvalidReferencePrice
	}
	if req.Source == "" {
		return nil, ErrInvalidReferenceSource
	}

	return &RecordReferencePriceResponse{
		ProductID: req.ProductID,
		Price:     req.Price,
	}, nil
}
