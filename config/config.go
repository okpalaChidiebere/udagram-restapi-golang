package config

import "os"

type postgres struct {
	Username string
	Password string
	Database string
	Host     string
	Dialect  string
}

type aws struct {
	Aws_region       string
	Aws_profile      string
	Aws_media_bucket string
}

type Config struct {
	postgres
	aws
}

func NewConfig() *Config {
	return &Config{
		postgres{
			Username: os.Getenv("POSTGRES_USERNAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DATABASE"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Dialect:  "postgres",
		},
		aws{
			Aws_region:       os.Getenv("AWS_REGION"),
			Aws_profile:      os.Getenv("AWS_PROFILE"),
			Aws_media_bucket: os.Getenv("AWS_MEDIA_BUCKET"),
		},
	}
}
