#!/bin/bash
set -e

# Функция для ожидания старта RabbitMQ
function wait_for_rabbitmq() {
    until curl -sI http://localhost:15672 | grep "200 OK" > /dev/null; do
        echo "Waiting for RabbitMQ to be ready..."
        sleep 5
    done
}

# Запуск RabbitMQ сервера в фоне
echo "Starting RabbitMQ server..."
docker-entrypoint.sh rabbitmq-server &

# Ожидание старта RabbitMQ
echo "Waiting for RabbitMQ to start..."
wait_for_rabbitmq

# Настройка RabbitMQ
echo "Setting up RabbitMQ..."

# Создание Exchange
rabbitmqadmin declare exchange name=aggregation_logs type=fanout durable=true

# Создание очереди
rabbitmqadmin declare queue name=log_queue durable=true

# Привязка очереди к обменнику
rabbitmqadmin declare binding source=aggregation_logs destination=log_queue

echo "RabbitMQ setup completed!"

# Ожидание завершения RabbitMQ сервера
wait
