package data

import (
	"context"
	"encoding/json"
	"time"

	"file-management-service/internal/biz"
	"file-management-service/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/segmentio/kafka-go"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewMinioClient, NewKafkaProducer, NewFileRepo)

// Data .
type Data struct {
	minio *minio.Client
	kafka *kafka.Writer
	log   *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, minioClient *minio.Client, kafkaProducer *kafka.Writer) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		if err := kafkaProducer.Close(); err != nil {
			log.NewHelper(logger).Error(err)
		}
	}
	return &Data{minio: minioClient, kafka: kafkaProducer, log: log.NewHelper(logger)}, cleanup, nil
}

func NewMinioClient(c *conf.Data, logger log.Logger) (*minio.Client, error) {
	log := log.NewHelper(logger)
	log.Infof("connecting to minio at %s", c.Minio.Endpoint)
	return minio.New(c.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKeyId, c.Minio.SecretAccessKey, ""),
		Secure: c.Minio.UseSsl,
	})
}

func NewKafkaProducer(c *conf.Data, logger log.Logger) (*kafka.Writer, error) {
	log := log.NewHelper(logger)
	log.Infof("connecting to kafka at %s", c.Kafka.Addrs)
	writer := &kafka.Writer{
		Addr:     kafka.TCP(c.Kafka.Addrs),
		Topic:    "file-uploads", // This should be configurable
		Balancer: &kafka.LeastBytes{},
	}
	return writer, nil
}

// fileRepo is a repository for files.
type fileRepo struct {
	data *Data
}

// NewFileRepo creates a new file repository.
func NewFileRepo(data *Data) biz.FileRepo {
	return &fileRepo{
		data: data,
	}
}

const (
	bucketName = "articles" // This should be configurable
)

func (r *fileRepo) Upload(ctx context.Context, file *biz.File) (string, error) {
	// Create bucket if it does not exist
	exists, err := r.data.minio.BucketExists(ctx, bucketName)
	if err != nil {
		return "", err
	}
	if !exists {
		err = r.data.minio.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", err
		}
	}

	// Upload the file
	info, err := r.data.minio.PutObject(ctx, bucketName, file.Name, file.Content, file.Size, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	r.data.log.Infof("Successfully uploaded %s of size %d", file.Name, info.Size)
	return info.Key, nil
}

func (r *fileRepo) Download(ctx context.Context, path string) (*biz.File, error) {
	object, err := r.data.minio.GetObject(ctx, bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	stat, err := object.Stat()
	if err != nil {
		return nil, err
	}

	return &biz.File{
		Name:    path,
		Content: object,
		Size:    stat.Size,
	}, nil
}

type FileUploadEvent struct {
	FilePath string `json:"file_path"`
	Format   string `json:"format"`
	Timestamp int64 `json:"timestamp"`
}

func (r *fileRepo) PublishUploadEvent(ctx context.Context, filePath string, format string) error {
	event := FileUploadEvent{
		FilePath: filePath,
		Format: format,
		Timestamp: time.Now().Unix(),
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return r.data.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte(filePath),
		Value: eventBytes,
	})
}
