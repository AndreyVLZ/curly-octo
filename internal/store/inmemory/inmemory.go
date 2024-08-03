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
	// fmt.Printf("userStore_find: %v\n", s.userStore)
	for _, el := range s.userStore {
		if el.Login() == name {
			return el, nil
		}
	}

	return nil, ErrUserNotFind
}

func (s *Store) Add(ctx context.Context, user model.User) error {
	// fmt.Printf("addUSER: %v\n", user)
	for _, el := range s.userStore {
		// fmt.Printf("add login: %s\n", el.Login())
		if el.Login() == user.Login() {
			return ErrLoginIsExist
		}
	}

	id := user.ID()
	s.userStore[id] = &user
	// fmt.Printf("userStore_add: %v\n", s.userStore)

	return nil
}

func (s *Store) GetData(_ context.Context, userID, id string) (model.Data, error) {
	ids := ids{
		userID: userID,
		id:     id,
	}

	//fmt.Printf("stre: %v\n", s.store)
	//fmt.Printf("ids: %v\n", ids)

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

	//fmt.Printf("save ok: %s[%d]\n", id, len(data.Data()))

	return nil
}

func (s *Store) SaveArray(_ context.Context, userID string, arr []*model.Data) error {
	//fmt.Printf("store-userID: %s\n", userID)
	//fmt.Printf("сохранено: [%d] от %s\n", len(arr), userID)
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
	//fmt.Printf("найдено: [%d] от %s\n", len(arr), userID)

	return arr, nil
}
