package inmemory

import (
	"context"
	"errors"

	"github.com/AndreyVLZ/curly-octo/internal/model"
)

var (
	ErrExist        = errors.New("data is exist")
	ErrNotFind      = errors.New("not find")
	ErrLoginIsExist = errors.New("login уже есть")
	ErrUserNotFind  = errors.New("not find user by name")
)

type ids struct {
	userID string
	id     string
}

type Store struct {
	store     map[ids]*model.Data
	userStore map[string]*model.User
}

func New() *Store {
	return &Store{
		store:     make(map[ids]*model.Data),
		userStore: make(map[string]*model.User),
	}
}

func (s *Store) FindByName(name string) (*model.User, error) {
	for _, el := range s.userStore {
		if el.Login() == name {
			return el, nil
		}
	}

	return nil, ErrUserNotFind
}

func (s *Store) Add(ctx context.Context, user model.User) error {
	for _, el := range s.userStore {
		if el.Login() == user.Login() {
			return ErrLoginIsExist
		}
	}

	id := user.ID()
	s.userStore[id] = &user

	return nil
}

func (s *Store) GetData(_ context.Context, userID, id string) (model.Data, error) {
	ids := ids{
		userID: userID,
		id:     id,
	}

	d, isExist := s.store[ids]
	if !isExist {
		return model.Data{}, ErrNotFind
	}

	return *d, nil
}

func (s *Store) SaveData(_ context.Context, userID string, data model.Data) error {
	id := data.ID()

	ids := ids{
		userID: userID,
		id:     id,
	}

	if _, ok := s.store[ids]; ok {
		return ErrExist
	}

	s.store[ids] = &data

	return nil
}

func (s *Store) SaveArray(_ context.Context, userID string, arr []*model.Data) error {
	for i := range arr {
		s.store[ids{userID: userID, id: arr[i].ID()}] = arr[i]
	}

	return nil
}

func (s *Store) List(_ context.Context, userID string) ([]*model.Data, error) {
	arr := make([]*model.Data, 0)

	for ids := range s.store {
		if ids.userID != userID {
			continue
		}

		el := s.store[ids]

		arr = append(arr, el)
	}

	return arr, nil
}
