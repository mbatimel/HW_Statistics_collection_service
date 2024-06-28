package migrate

import (

	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/mbatimel/HW_Statistics_collection_service/internal/config"

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

    connStr := fmt.Sprintf("tcp://%s:%s?database=%s", config.ClickHouse.Host, config.ClickHouse.Port, config.ClickHouse.DB)
    fmt.Println("Connecting to ClickHouse with:", connStr)
    db, err := sql.Open("clickhouse", connStr)
    if err != nil {
        return fmt.Errorf("Open failed: %v", err)
    }
    defer db.Close()


    if err := db.Ping(); err != nil {
        return fmt.Errorf("Ping failed: %v", err)
    }

    _, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.ClickHouse.DB))
    if err != nil {
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
    _, err = db.Exec(queries)
    if err != nil {
        return fmt.Errorf("Exec migration queries failed: %v", err)
    }

    return nil
}
