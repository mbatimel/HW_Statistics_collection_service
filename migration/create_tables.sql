-- create_tables.sql
CREATE TABLE IF NOT EXISTS OrderBook (
    id Int64,
    exchange String,
    pair String,
    asks Array(Tuple(Float64, Float64)),
    bids Array(Tuple(Float64, Float64))
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS HistoryOrder (
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
ORDER BY (client_name, exchange_name, time_placed);

CREATE TABLE IF NOT EXISTS Client (
    client_name String,
    exchange_name String,
    label String,
    pair String
) ENGINE = MergeTree()
ORDER BY (client_name, exchange_name);
