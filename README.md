# SSO (Single Sign-On) Service

Микросервис аутентификации и авторизации с поддержкой JWT-токенов и gRPC API.

## Функционал

- Регистрация новых пользователей
- Аутентификация по email/password
- Генерация JWT-токенов
- Проверка прав администратора
- Поддержка нескольких приложений (multi-tenant)
- Готовые метрики для мониторинга

## Технологии

- **Язык**: Go 1.21+
- **gRPC**: высокопроизводительный RPC-фреймворк
- **PostgreSQL**: хранение пользователей и приложений
- **JWT**: безопасные токены доступа
- **Slog**: логирование

## Быстрый старт

### Требования
- Go 1.21+
- PostgreSQL 15+
- protoc (для генерации кода)

## **Функциональность**
- Регистрация пользователей с хешированием пароля  
- Аутентификация по email и паролю  
- Проверка роли администратора  
- Работа с приложениями (Apps)  
- Подключение к PostgreSQL  
- Graceful shutdown  

## **Запуск проекта**  

### **Требования**  
- **Go 1.21+**  
- **PostgreSQL 14+**  

### **Запуск сервиса**
```sh
go run cmd/sso/main.go
```

## **API gRPC**

### **Примеры запросов**
#### **Регистрация пользователя**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

## **Graceful Shutdown**
При завершении работы сервера:
- **Закрывается gRPC-сервер**  
- **Отключается соединение с PostgreSQL**  
