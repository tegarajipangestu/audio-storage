package common

import (
	"fmt"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Env       string `env:"ENV,default=development"`
	ChunkSize int    `env:"CHUNK_SIZE"`
	TempDir   string `env:"TEMP_DIR"`

	Minio    Minio
	Postgres Postgres
}

type Minio struct {
	MinioHost      string `env:"MINIO_HOST"`
	MinioPort      int    `env:"MINIO_PORT"`
	MinioAccessKey string `env:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `env:"MINIO_SECRET_KEY"`
	BucketName     string `env:"MINIO_BUCKET_NAME"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	Username string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Name     string `env:"POSTGRES_DB_NAME,required"`
}

func (db *Postgres) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db.Username, db.Password, db.Host, db.Port, db.Name)
}

func NewConfig(env string) (*Config, error) {
	_ = godotenv.Load(env)

	var config Config
	if err := envdecode.Decode(&config); err != nil {
		return nil, errors.Wrap(err, "[NewConfig] error decoding env")
	}

	return &config, nil
}
