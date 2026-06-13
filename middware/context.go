package middware

import (
	"context"
	"errors"

	"github.com/shawnyu5/debate_dragon_2.0/db"
)

type contextKey struct{}

var storeKey = contextKey{}

// ContextWithStore inject a db.Store into the context.
//
// Returns a new context.Context with the db store injected
func ContextWithStore(ctx context.Context, store *db.Store) context.Context {
	return context.WithValue(ctx, storeKey, store)
}

// StoreFromContext extract a db.Store from context
func StoreFromContext(ctx context.Context) (*db.Store, error) {
	store, ok := ctx.Value(storeKey).(*db.Store)
	if !ok {
		return nil, errors.New("database store not found in context")
	}
	return store, nil
}
