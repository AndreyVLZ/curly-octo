package localstore

import (
	"context"
	"testing"

	"github.com/AndreyVLZ/curly-octo/internal/model"
)

type mockStore struct {
	userID string
	t      *testing.T
}

func (ms *mockStore) List(ctx context.Context, userID string) ([]*model.Data, error) {
	return nil, nil
}

func (ms *mockStore) SaveArray(ctx context.Context, userID string, arr []*model.Data) error {
	return nil
}

func (ms *mockStore) GetData(ctx context.Context, userID, dataID string) (model.Data, error) {
	return model.Data{}, nil
}

func (ms *mockStore) SaveData(ctx context.Context, userID string, data model.Data) error {
	return nil
}

func TestLocalStore(t *testing.T) {

}
