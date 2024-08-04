package syncservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/AndreyVLZ/curly-octo/internal/model"
)

type iClient interface {
	SendData(ctx context.Context, arr []*model.Data) error
	SendFiles(ctx context.Context, files []*model.EncFile) error
	RecvData(ctx context.Context) ([]*model.Data, error)
	RecvFiles(ctx context.Context, files []*model.DecFile) error
}

type iStoreService interface {
	GetAll(ctx context.Context) ([]*model.Data, []*model.EncFile, error)
	SaveArray(ctx context.Context, arr []*model.Data) ([]*model.DecFile, error)
}

// SyncService Отвечает за отправку/приём данный.
type SyncService struct {
	storeService iStoreService
	client       iClient
}

func NewSyncService(client iClient, storeSrv iStoreService) *SyncService {
	return &SyncService{
		storeService: storeSrv,
		client:       client,
	}
}

// Send Отправляет данные.
func (ss *SyncService) Send(ctx context.Context) error {
	var (
		err   error
		arr   []*model.Data
		files []*model.EncFile
	)

	// получаем все данные
	arr, files, err = ss.storeService.GetAll(ctx)
	if err != nil {
		err = fmt.Errorf("storeService GetAll: %w", err)

		return err
	}

	// закрываем все открытые файлы при выходе
	defer func() {
		err = errors.Join(err, closeEncFiles(files))
	}()

	// отправляем данные
	if err = ss.client.SendData(ctx, arr); err != nil {
		err = fmt.Errorf("send data: %w", err)

		return err
	}

	// отправляем файлы
	if err = ss.client.SendFiles(ctx, files); err != nil {
		err = fmt.Errorf("send files: %w", err)

		return err
	}

	return nil
}

// Recv Принимает файлы.
func (ss *SyncService) Recv(ctx context.Context) error {
	var (
		err      error
		arr      []*model.Data
		decFiles []*model.DecFile
	)

	// получаем данные
	arr, err = ss.client.RecvData(ctx)
	if err != nil {
		err = fmt.Errorf("recvData: %w", err)

		return err
	}

	// сохраняем данные
	decFiles, err = ss.storeService.SaveArray(ctx, arr)
	if err != nil {
		err = fmt.Errorf("saveArray: %w", err)

		return err
	}

	// закрываем все файлы при выходе
	defer func() {
		err = errors.Join(err, closeFiles(decFiles))
	}()

	// получаем файлы
	if err = ss.client.RecvFiles(ctx, decFiles); err != nil {
		err = fmt.Errorf("recvFiles: %w", err)

		return err
	}

	return err
}

// closeEncFiles Закрывает файлы.
func closeEncFiles(files []*model.EncFile) error {
	errs := make([]error, 0, len(files))
	for i := range files {
		errs = append(errs, files[i].Close())
	}

	return errors.Join(errs...)
}

// closeFiles Закрывает файлы.
func closeFiles(files []*model.DecFile) error {
	errs := make([]error, 0, len(files))
	for i := range files {
		errs = append(errs, files[i].Close())
	}

	return errors.Join(errs...)
}
