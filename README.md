# NoteProject

NoteProject - это веб-приложение для управления заметками. Оно предоставляет API для создания, просмотра, редактирования и удаления заметок, а также для отслеживания их статуса.

## Функциональность

- Создание новых заметок с указанием заголовка, текста и статуса выполнения.
- Получение списка всех заметок с возможностью фильтрации по статусу выполнения.
- Получение и редактирование отдельных заметок.
- Удаление заметок.

## Использованные технологии

NoteProject разработан с использованием следующих технологий:

1. **Swagger**: Используется для документирования и взаимодействия с API.
2. **Redis**: Используется для хранения сессий пользователей.
3. **SQL Postgres**: Используется для хранения данных о заметках и пользователей.
4. **Авторизация через JWT токен**: Используется для аутентификации пользователей и защиты ресурсов.
5. **Гексагональная архитектура**: Используется для организации кода и разделения бизнес-логики от инфраструктуры.
6. **Тестирование с использованием моков**: Используется для тестирования различных компонентов приложения с помощью моков.
7. **Миграции**: Используются для управления изменениями в структуре базы данных.
8. **Логирование**: Используется для отслеживания действий и ошибок приложения.
9. **Docker**: Используется для контейнеризации и развертывания приложения.