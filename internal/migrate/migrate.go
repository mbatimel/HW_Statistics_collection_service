package migrate

import (
    "context"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"

    "github.com/ClickHouse/clickhouse-go/v2"
    "github.com/mbatimel/HW_Statistics_collection_service/internal/config"
    "gopkg.in/yaml.v3"
)

func loadConfig() (*config.Config, error) {
    configPath := filepath.Join("config", "config.yaml")
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config config.Config
    decoder := yaml.NewDecoder(file)
    err = decoder.Decode(&config)
    if err != nil {
        return nil, err
    }
    return &config, nil
}

func RunMigrations() error {
    config, err := loadConfig()
    if err != nil {
        return fmt.Errorf("loadConfig() failed: %v", err)
    }

    conn, err := clickhouse.Open(&clickhouse.Options{
        Addr: []string{fmt.Sprintf("%s:%s", config.ClickHouse.Host, config.ClickHouse.Port)},
        Auth: clickhouse.Auth{
            Database: config.ClickHouse.DB,
            Username: config.ClickHouse.Username,
            Password: config.ClickHouse.Password,
        },
        DialTimeout: 5 * time.Second,
        ConnOpenStrategy: clickhouse.ConnOpenRoundRobin,
    })
    if err != nil {
        return fmt.Errorf("failed to open ClickHouse connection: %v", err)
    }

    ctx := context.Background()

    if err := conn.Ping(ctx); err != nil {
        return fmt.Errorf("failed to ping ClickHouse: %v", err)
    }

    if err := conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.ClickHouse.DB)); err != nil {
        return fmt.Errorf("Exec CREATE DATABASE failed: %v", err)
    }

    sqlPath := filepath.Join("migration", "migration.sql")
    file, err := os.Open(sqlPath)
    if err != nil {
        return fmt.Errorf("Open migration file failed: %v", err)
    }
    defer file.Close()

    sqlData, err := io.ReadAll(file)
    if err != nil {
        return fmt.Errorf("ReadAll migration file failed: %v", err)
    }

    queries := string(sqlData)
    if err := conn.Exec(ctx, queries); err != nil {
        return fmt.Errorf("Exec migration queries failed: %v", err)
    }

    return nil
}
