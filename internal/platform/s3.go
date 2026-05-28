package platform

import (
	"fmt"

	"github.com/gkarman/demo/internal/config"
	s3storage "github.com/gkarman/demo/internal/infrastructure/storage/s3"
)

func NewS3Client(cfg *config.Config) (*s3storage.Client, error) {
	client, err := s3storage.NewClient(
		cfg.S3.Endpoint,
		cfg.S3.PublicURL,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.Bucket,
		cfg.S3.Region,
	)
	if err != nil {
		return nil, fmt.Errorf("init s3 client: %w", err)
	}
	return client, nil
}
