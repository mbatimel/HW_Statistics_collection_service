# Statistics Collection Service Documentation

## Overview

The Statistics Collection Service is a Go-based service that interfaces with ClickHouse to manage order books and order histories. This documentation provides an overview of the service, its API endpoints, and how to configure and run it.

## Table of Contents

1. [Configuration](#configuration)
2. [Running the Service](#running-the-service)
3. [API Endpoints](#api-endpoints)
   - [Get Order Book](#get-order-book)
   - [Save Order Book](#save-order-book)
   - [Get Order History](#get-order-history)
   - [Save Order History](#save-order-history)
4. [Database Migrations](#database-migrations)

## Configuration

The service configuration is specified in the `config/config.yaml` file. Below is an example configuration:

```yaml
server:
  host: "localhost"
  port: "8080"

clickhouse:
  host: "localhost"
  port: "9005"
  db: "my_database"
```

- **server**: Contains the server configuration.
  - `host`: The hostname or IP address the server listens on.
  - `port`: The port the server listens on.
- **clickhouse**: Contains the ClickHouse database configuration.
  - `host`: The hostname or IP address of the ClickHouse server.
  - `port`: The port of the ClickHouse server.
  - `db`: The name of the ClickHouse database.

## Running the Service

To run the service, follow these steps:

1. Ensure you have Go installed.
2. Clone the repository and navigate to the project directory.
3. Run the following command to start the server:

```sh
go run cmd/server/main.go
```

The server will start and listen on the address specified in the configuration file.

## API Endpoints

The service exposes the following API endpoints:

### Get Order Book

- **Endpoint**: `/get-order-book`
- **Method**: GET
- **Parameters**:
  - `exchange_name`: Name of the exchange (e.g., `Binance`).
  - `pair`: Currency pair (e.g., `BTC/USD`).
- **Description**: Retrieves the order book for the specified exchange and currency pair.

#### Example Request

```sh
curl http://localhost:8080/get-order-book?exchange_name=Binance&pair=BTC/USD
```

### Save Order Book

- **Endpoint**: `/save-order-book`
- **Method**: POST
- **Parameters**: 
  - `exchange_name`: Name of the exchange (e.g., `Binance`).
  - `pair`: Currency pair (e.g., `BTC/USD`).
- **Request Body**: JSON array of order book entries.
- **Description**: Saves the order book for the specified exchange and currency pair.

#### Example Request

```sh
curl -X POST http://localhost:8080/save-order-book?exchange_name=Binance&pair=BTC/USD -H "Content-Type: application/json" -d '[
  {
    "price": 10000.5,
    "base_qty": 0.1
  },
  {
    "price": 10001.0,
    "base_qty": 0.2
  }
]
'
```

### Get Order History

- **Endpoint**: `/get-order-history`
- **Method**: POST
- **Request Body**: JSON object containing client information.
- **Description**: Retrieves the order history for the specified client.

#### Example Request

```sh
curl -X POST http://localhost:8080/get-order-history -H "Content-Type: application/json" -d '{
  "client_name": "Alice",
  "exchange_name": "Binance",
  "label": "Client1",
  "pair": "BTC/USD"
}'
```

### Save Order History

- **Endpoint**: `/save-order-history`
- **Method**: POST
- **Request Body**: JSON object containing order history information.
- **Description**: Saves an order history entry for a client.

#### Example Request

```sh
curl -X POST http://localhost:8080/save-order-history -H "Content-Type: application/json" -d '{
  "client_name": "Alice",
  "exchange_name": "Binance",
  "label": "Order1",
  "pair": "BTC/USD",
  "side": "buy",
  "type_order": "limit",
  "base_qty": 0.1,
  "price": 10000.5,
  "algorithm_name_placed": "algo1",
  "lowest_sell_prc": 10000.0,
  "highest_buy_prc": 10001.0,
  "commission_quote_qty": 0.0001,
  "time_placed": "2024-06-28T12:00:00Z"
}

'
```

## Database Migrations

To run database migrations, follow these steps:

1. Ensure the ClickHouse server is running and accessible.
2. Run the following command to execute the migrations:

```sh
go run cmd/migrate/main.go
```

This will create the necessary database and tables specified in the `migration/migration.sql` file.

### Migration SQL Example

```sql
CREATE DATABASE IF NOT EXISTS my_database;

CREATE TABLE IF NOT EXISTS order_book (
    exchange String,
    pair String,
    price Float64,
    base_qty Float64
) ENGINE = MergeTree()
ORDER BY (exchange, pair);

CREATE TABLE IF NOT EXISTS order_history (
    client_name String,
    exchange_name String,
    label String,
    pair String,
    side String,
    type_order String,
    base_qty Float64,
    price Float64,
    algorithm_name_placed String,
    lowest_sell_prc Float64,
    highest_buy_prc Float64,
    commission_quote_qty Float64,
    time_placed DateTime
) ENGINE = MergeTree()
ORDER BY (client_name, exchange_name, label, pair);
```

This script sets up the necessary database and tables for the service. Modify the script as needed to fit your database schema and requirements.