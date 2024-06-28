package config

type Config struct {
	Server		Server		`yaml:"server"`
	ClickHouse	ClickHouse	`yaml:"clickhouse"`
}