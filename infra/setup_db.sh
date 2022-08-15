#!/bin/bash
echo wait for servers to be up
sleep 20

HOSTPARAMS="--host node_1 --insecure"
SQL="/cockroach/cockroach.sh sql $HOSTPARAMS"
$SQL -e "DROP DATABASE sub_account_db;"
$SQL -e "DROP DATABASE exchange_orders_db;"
$SQL -e "DROP DATABASE transactions_db;"

$SQL -e "CREATE DATABASE IF NOT EXISTS sub_account_db;"
$SQL -e "CREATE DATABASE IF NOT EXISTS exchange_orders_db;"
$SQL -e "CREATE DATABASE IF NOT EXISTS transactions_db;"
$SQL -e "CREATE USER postgres;"
$SQL -e "ALTER USER postgres WITH PASSWORD postgres;"
$SQL -d sub_account_db -e "CREATE TABLE IF NOT EXISTS sub_accounts
(
    sub_account_id  varchar(40) primary key,
    account_number  varchar(40),
    balance         decimal,
    linked_accounts text,
    status          varchar(10),
    credentials      text
);"

$SQL -d exchange_orders_db -e "CREATE TABLE IF NOT EXISTS exchange_orders
(
    transaction_id varchar(40) primary key,
    fro varchar(100),
    fund           varchar(100),
    amt         decimal,
    re        varchar(200),
    received_at    TIMESTAMPTZ,
    validated_at   TIMESTAMPTZ
);"

$SQL -d transactions_db -e "CREATE TABLE IF NOT EXISTS transactions
(
    transaction_id varchar(40) primary key,
    fro varchar(100),
    fund                 varchar(100),
    amt                  decimal,
    re                   varchar(200),
    received_at          TIMESTAMPTZ,
    validated_at         TIMESTAMPTZ,
    completed_at         TIMESTAMPTZ,
    fund_request_made_at TIMESTAMPTZ,
    status         varchar(10),
    system_notes   varchar(200)
);"
$SQL -e "GRANT ALL ON TABLE sub_account_db.sub_accounts TO postgres;"
$SQL -e "GRANT ALL ON TABLE exchange_orders_db.exchange_orders TO postgres;"
$SQL -e "GRANT ALL ON TABLE transactions_db.transactions TO postgres;"
