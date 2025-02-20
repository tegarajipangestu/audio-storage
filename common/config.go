package common

import (
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Env       string `env:"ENV,default=development"`
	ChunkSize int    `env:"CHUNK_SIZE"`
	TempDir   string `env:"TEMP_DIR"`

	Minio Minio
}

type Minio struct {
	MinioHost      string `env:"MINIO_HOST"`
	MinioPort      int    `env:"MINIO_PORT"`
	MinioAccessKey string `env:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `env:"MINIO_SECRET_KEY"`
	BucketName     string `env:"MINIO_BUCKET_NAME"`
}

func NewConfig(env string) (*Config, error) {
	_ = godotenv.Load(env)

	var config Config
	if err := envdecode.Decode(&config); err != nil {
		return nil, errors.Wrap(err, "[NewServerConfig] error decoding env")
	}

	return &config, nil
}
