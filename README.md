# Мини-сервис "Цитатник"

Это REST API-сервис, написанный на языке Go для управления и хранения цитат. Сервис реализует функционал добавления, получения, фильтрации и удаления цитат, соответствующий заданным требованиям. Данные хранятся в SQLite, как указано в структуре проекта. Для маршрутизации используется стандартный пакет net/http из библиотеки Go.
## Возможности
Сервис поддерживает следующие эндпоинты:

POST /quotes: Добавление новой цитаты.
GET /quotes: Получение всех цитат.
GET /quotes/random: Получение случайной цитаты.
GET /quotes?author={author}: Фильтрация цитат по автору.
DELETE /quotes/{id}: Удаление цитаты по идентификатору.

## Требования
Go: Версия 1.24.3 или выше.
SQLite: Для постоянного хранения данных.
Настроенная среда Go с файлом go.mod.

## Установка
Склонируйте репозиторий:
`git clone https://github.com/Grino777/quotes.git`
`cd quotes`

Убедитесь, что у вас установлен Go (версия 1.24.3).
Установите зависимости:go mod tidy

## Запуск
Перейдите в директорию с основным файлом:cd cmd/quotes

Запустите приложение:
`go run main.go`
Сервис будет доступен по адресу http://localhost:8080.

## Конфигурация
Конфигурация сервиса находится в файле configs/local.yml. Основные параметры:

sqlite:
  local_path: "storage/quotes.sqlite"
api:
  addr: "127.0.0.1"
  port: "8080"

## Использование
Проверочные команды для тестирования API с помощью curl:

Добавление цитаты:
`curl -X POST http://localhost:8080/quotes \
-H "Content-Type: application/json" \
-d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'`

Получение всех цитат:
`curl http://localhost:8080/quotes`

Получение случайной цитаты:
`curl http://localhost:8080/quotes/random`

Фильтрация цитат по автору:
`curl http://localhost:8080/quotes?author=Confucius`

Удаление цитаты по ID:
`curl -X DELETE http://localhost:8080/quotes/1`

## Описание директорий

cmd/quotes: Точка входа приложения (main.go).
configs: Конфигурационные файлы.
internal: Основная логика приложения:
- api: Обработчики HTTP-запросов и middleware с использованием пакета net/http.
- app: Инициализация приложения и сервера.
- config: Логика загрузки конфигурации.
- domain/models: Модель данных для цитат.
- interfaces: Интерфейсы для сервисов и хранилища.
- lib/logger: Логирование.
- services: Бизнес-логика API.
- storage: Реализация хранилища (in-memory или SQLite).
- utils: Вспомогательные утилиты.

storage: Файл базы данных SQLite.

## Зависимости
github.com/mattn/go-sqlite3 - драйвер для sqlite3
github.com/ilyakaznacheev/cleanenv - парсинг конфиг файла
