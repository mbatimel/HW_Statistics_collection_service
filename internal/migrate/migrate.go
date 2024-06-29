package migrate

import (
    "context"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
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

func runSQLFile(conn clickhouse.Conn, ctx context.Context, filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("Open migration file failed: %v", err)
    }
    defer file.Close()

    sqlData, err := io.ReadAll(file)
    if err != nil {
        return fmt.Errorf("ReadAll migration file failed: %v", err)
    }

    queries := strings.Split(string(sqlData), ";")
    for _, query := range queries {
        query = strings.TrimSpace(query)
        if query == "" {
            continue
        }
        log.Printf("Executing query: %s", query)
        if err := conn.Exec(ctx, query); err != nil {
            return fmt.Errorf("Exec migration query failed: %v: %s", err, query)
        }
        log.Printf("Successfully executed query: %s", query)
    }

    return nil
}


func RunMigrations() error {
    config, err := loadConfig()
    if err != nil {
        return fmt.Errorf("loadConfig() failed: %v", err)
    }

    // Connect to ClickHouse with default user to create the database and user
    conn, err := clickhouse.Open(&clickhouse.Options{
        Addr: []string{fmt.Sprintf("%s:%s", config.ClickHouse.Host, config.ClickHouse.Port)},
        Auth: clickhouse.Auth{
            Username: "default",
            Password: "default_password",
        },
        DialTimeout: 5 * time.Minute,
        ConnOpenStrategy: clickhouse.ConnOpenRoundRobin,
    })
    if err != nil {
        return fmt.Errorf("failed to open ClickHouse connection: %v", err)
    }

    ctx := context.Background()

    log.Println("Pinging ClickHouse with default user...")
    if err := conn.Ping(ctx); err != nil {
        return fmt.Errorf("failed to ping ClickHouse: %v", err)
    }
    log.Println("Ping successful.")

    log.Println("Creating database if not exists...")
    if err := conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.ClickHouse.DB)); err != nil {
        return fmt.Errorf("Exec CREATE DATABASE failed: %v", err)
    }
    log.Println("Database created or already exists.")

    // Connect to ClickHouse with the new user and create tables
    conn, err = clickhouse.Open(&clickhouse.Options{
        Addr: []string{fmt.Sprintf("%s:%s", config.ClickHouse.Host, config.ClickHouse.Port)},
        Auth: clickhouse.Auth{
            Database: config.ClickHouse.DB,
            Username: config.ClickHouse.Username,
            Password: config.ClickHouse.Password,
        },
        DialTimeout: 5 * time.Minute,
        ConnOpenStrategy: clickhouse.ConnOpenRoundRobin,
    })
    if err != nil {
        return fmt.Errorf("failed to open ClickHouse connection with database: %v", err)
    }

    log.Println("Pinging ClickHouse with database...")
    if err := conn.Ping(ctx); err != nil {
        return fmt.Errorf("failed to ping ClickHouse with database: %v", err)
    }
    log.Println("Ping successful with database.")

    log.Println("Running migration SQL file...")
    if err := runSQLFile(conn, ctx, filepath.Join("migration", "create_tables.sql")); err != nil {
        return err
    }
    log.Println("Migration SQL file executed successfully.")

    return nil
}
