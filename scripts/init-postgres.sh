#!/bin/bash

# Запуск PostgreSQL сервера в фоне
echo "Starting RabbitMQ server..."

# Проверяем доступность PostgreSQL
until pg_isready -h postgres -U "$POSTGRES_USER"; do
  sleep 1
done

# Добавляем базу данных 'NoteProject'
  psql -h postgres -U "$POSTGRES_USER" -c "CREATE DATABASE \"$POSTGRES_DB\""





