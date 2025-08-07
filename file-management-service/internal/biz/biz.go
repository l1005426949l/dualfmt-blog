package biz

import (
	"context"
	"io"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewFileUsecase)

type File struct {
	Name    string
	Content io.Reader
	Size    int64
}

type FileRepo interface {
	Upload(ctx context.Context, file *File) (string, error)
	Download(ctx context.Context, path string) (*File, error)
	PublishUploadEvent(ctx context.Context, filePath string, format string) error
}

type FileUsecase struct {
	repo FileRepo
	log  *log.Helper
}

func NewFileUsecase(repo FileRepo, logger log.Logger) *FileUsecase {
	return &FileUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *FileUsecase) UploadFile(ctx context.Context, file *File, format string) (string, error) {
	uc.log.Infof("Uploading file: %s", file.Name)

	filePath, err := uc.repo.Upload(ctx, file)
	if err != nil {
		return "", err
	}

	err = uc.repo.PublishUploadEvent(ctx, filePath, format)
	if err != nil {
		uc.log.Errorf("failed to publish upload event for file %s: %v", filePath, err)
	}

	return filePath, nil
}

func (uc *FileUsecase) DownloadFile(ctx context.Context, path string) (*File, error) {
    uc.log.Infof("Downloading file: %s", path)
    return uc.repo.Download(ctx, path)
}
