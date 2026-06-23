package store

import "context"

type StoreRepository interface {
	Save(ctx context.Context, store *Store) error
	FindByID(ctx context.Context, id int64) (*Store, error)
	FindBySellerID(ctx context.Context, sellerID int64) (*Store, error)
	Update(ctx context.Context, store *Store) error
}
