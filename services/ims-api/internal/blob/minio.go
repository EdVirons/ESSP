package blob

import (
	"context"
	"time"

	"github.com/edvirons/ssp/ims/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIO struct {
	Client *minio.Client
	Bucket string
	Expiry time.Duration
}

func NewMinIO(cfg config.Config) (*MinIO, error) {
	cl, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
		Region: cfg.MinIORegion,
	})
	if err != nil {
		return nil, err
	}

	exp := time.Duration(cfg.MinIOPresignExpirySeconds) * time.Second
	if exp <= 0 {
		exp = 15 * time.Minute
	}

	return &MinIO{
		Client: cl,
		Bucket: cfg.AttachmentsBucket,
		Expiry: exp,
	}, nil
}

func (m *MinIO) EnsureBucket(ctx context.Context) error {
	exists, err := m.Client.BucketExists(ctx, m.Bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return m.Client.MakeBucket(ctx, m.Bucket, minio.MakeBucketOptions{Region: ""})
}

func (m *MinIO) PresignPut(ctx context.Context, objectKey string, contentType string) (string, error) {
	reqParams := map[string][]string{}
	if contentType != "" {
		reqParams["response-content-type"] = []string{contentType}
	}
	u, err := m.Client.PresignedPutObject(ctx, m.Bucket, objectKey, m.Expiry)
	if err != nil {
		return "", err
	}
	_ = reqParams // reserved for future, PresignedPutObject doesn't accept query params.
	return u.String(), nil
}

func (m *MinIO) PresignGet(ctx context.Context, objectKey string) (string, error) {
	u, err := m.Client.PresignedGetObject(ctx, m.Bucket, objectKey, m.Expiry, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
