package localstore

import (
	"context"

	"github.com/AndreyVLZ/curly-octo/internal/model"
)

type storager interface {
	List(ctx context.Context, userID string) ([]*model.Data, error)
	SaveArray(ctx context.Context, userID string, arr []*model.Data) error
	GetData(ctx context.Context, userID, dataID string) (model.Data, error)
	SaveData(ctx context.Context, userID string, data model.Data) error
}

// LocalStore обертка над store без userID.
type LocalStore struct {
	userName string
	store    storager
}

func NewLocalStore(userName string, store storager) *LocalStore {
	return &LocalStore{userName: userName, store: store}
}

func (ls *LocalStore) List(ctx context.Context) ([]*model.Data, error) {
	return ls.store.List(ctx, ls.userName)
}

func (ls *LocalStore) SaveArray(ctx context.Context, arr []*model.Data) error {
	return ls.store.SaveArray(ctx, ls.userName, arr)
}

func (ls *LocalStore) GetData(ctx context.Context, dataID string) (model.Data, error) {
	return ls.store.GetData(ctx, ls.userName, dataID)
}

func (ls *LocalStore) SaveData(ctx context.Context, data model.Data) error {
	return ls.store.SaveData(ctx, ls.userName, data)
}
