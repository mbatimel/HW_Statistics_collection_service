-- initial_migration.sql
CREATE USER IF NOT EXISTS my_user IDENTIFIED BY 'my_password';
GRANT ALL ON my_database.* TO my_user;
CREATE DATABASE IF NOT EXISTS my_database;
