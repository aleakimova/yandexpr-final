# Планировщик задач --- Дипломный проект

## Описание

Веб-сервер на Go для управления задачами. Поддерживает создание, редактирование, удаление задач и расписание повторений (ежедневно, еженедельно, ежемесячно, ежегодно).

## Функционал

- Создание, просмотр, редактирование и удаление задач
- Расписание повторений: ежедневно (`d <n>`), еженедельно (`w <дни>`), ежемесячно (`m <числа> [месяцы]`), ежегодно (`y`)
- Умная нормализация дат: просроченные задачи с повторением автоматически переносятся на следующий подходящий день
- Поиск задач по тексту и дате
- JWT-аутентификация (опционально, через `TODO_PASSWORD`)
- Фронтенд, встроенный в сервер (папка `web/`)

## Файловая структура

- [`cmd/main.go`](cmd/main.go) --- точка входа: инициализация БД, роуты, раздача статики
- [`internal/db/db.go`](internal/db/db.go) --- работа с бд (`Start`, `Get`)
- `pkg/api/` --- HTTP-обработчики:
  - [`api.go`](pkg/api/api.go) --- регистрация роутов (`Init`)
  - [`nextdate.go`](pkg/api/nextdate.go) --- `NextDate(now, dstart, repeat)` вычисляет следующую дату задачи; `GET /api/nextdate` --- HTTP-обёртка над этой функцией
  - [`task.go`](pkg/api/task.go) --- роутер `GET/POST/PUT/DELETE /api/task`
  - [`addtask.go`](pkg/api/addtask.go) --- `POST /api/task` --- создание задачи
  - [`gettask.go`](pkg/api/gettask.go) --- `GET /api/task?id=` --- получение задачи
  - [`edittask.go`](pkg/api/edittask.go) --- `PUT /api/task` --- редактирование
  - [`deletetask.go`](pkg/api/deletetask.go) --- `DELETE /api/task?id=` --- удаление
  - [`donetask.go`](pkg/api/donetask.go) --- `POST /api/task/done?id=` --- отметка выполнения
  - [`tasks.go`](pkg/api/tasks.go) --- `GET /api/tasks?search=` --- список задач
  - [`signin.go`](pkg/api/signin.go) --- `POST /api/signin` --- вход
- `pkg/middleware/` --- middleware:
  - [`auth.go`](pkg/middleware/auth.go) --- middleware JWT-аутентификации
  - [`log.go`](pkg/middleware/log.go) --- middleware логирования запросов (`LogRequest`)
- [`internal/db/task.go`](internal/db/task.go) --- структура `Task` + SQL-операции
- `tests/` --- интеграционные тесты (требуют запущенного сервера):
  - [`settings.go`](tests/settings.go)
  - [`app_1_test.go`](tests/app_1_test.go)
  - [`db_2_test.go`](tests/db_2_test.go)
  - [`nextdate_3_test.go`](tests/nextdate_3_test.go)
  - [`addtask_4_test.go`](tests/addtask_4_test.go)
  - [`tasks_5_test.go`](tests/tasks_5_test.go)
  - [`task_6_test.go`](tests/task_6_test.go)
  - [`task_7_test.go`](tests/task_7_test.go)
- `web/` --- фронтенд
- [`Dockerfile`](Dockerfile)
- [`go.mod`](go.mod)

## API

| Метод | Путь | Описание |
|----------|---------|---------------|
| GET | `/api/nextdate` | Вычислить следующую дату по правилу повторения |
| GET | `/api/tasks` | Список задач (с поиском по тексту/дате) |
| POST | `/api/task` | Создать задачу |
| GET | `/api/task?id=` | Получить задачу по ID |
| PUT | `/api/task` | Редактировать задачу |
| DELETE | `/api/task?id=` | Удалить задачу |
| POST | `/api/task/done?id=` | Отметить задачу выполненной |
| POST | `/api/signin` | Авторизация (возвращает JWT) |

## Выполненные задания со звёздочкой

- **Поиск задач** --- `GET /api/tasks?search=` поддерживает фильтрацию по тексту и дате
- **Полная поддержка NextDate** --- правила `w` (по дням недели) и `m` (по числам месяца, включая `-1` / `-2`)
- **Аутентификация** --- JWT-токен через переменную окружения `TODO_PASSWORD`
- **Docker** --- многоэтапный `Dockerfile` с компиляцией и запуском в `ubuntu:latest`

## Запуск локально

```bash
# Запуск сервера (фронтенд из папки web/)
go run cmd/main.go web/

# Открыть в браузере
http://localhost:7540
```

Переменные окружения (опционально):

| Переменная | Описание | По умолчанию |
|----|----|----|
| `TODO_PORT` | Порт сервера | `7540` |
| `TODO_DB` | Путь к файлу базы данных | `scheduler.db` |
| `TODO_PASSWORD` | Пароль для входа (если не задан --- аутентификация отключена) | --- |
| `LOG_LEVEL` | Уровень логирования (`DEBUG`, `INFO`, `WARN`, `ERROR`) | `INFO` |

Пример с паролем:
```bash
TODO_PASSWORD=secret go run cmd/main.go web/
```

## Запуск тестов

Тесты требуют запущенного сервера:

```bash
# Терминал 1
go run cmd/main.go web/

# Терминал 2
cd tests && go test ./... -v
```

### Настройка `tests/settings.go`

| Параметр | По умолчанию | Описание |
|----|----|----|
| `Port` | `7540` | Порт сервера |
| `DBFile` | `"../scheduler.db"` | Путь к БД |
| `FullNextDate` | `true` | Включить тесты для правил `w` и `m` |
| `Search` | `true` | Включить тесты поиска |
| `Token` | `""` | JWT-токен для тестов с аутентификацией (см. ниже) |

### Тесты с аутентификацией

Если сервер запущен с `TODO_PASSWORD`:

1. Получить токен:
   ```bash
   curl -X POST http://localhost:7540/api/signin -d '{"password":"secret"}'
   # {"token":"eyJ..."}
   ```
2. Вставить значение токена в `tests/settings.go`:
   ```go
   var Token = `eyJ...`
   ```
3. Запустить тесты: `cd tests && go test ./... -v`

### Запуск одного теста

```bash
cd tests && go test -v -run TestNextDate
```

Доступные тесты: `TestApp`, `TestDB`, `TestNextDate`, `TestAddTask`, `TestTasks`, `TestTask`, `TestEditTask`, `TestDone`, `TestDelTask`

## Docker

### Сборка образа

```bash
docker build -t scheduler .
```

### Запуск контейнера

Минимальный запуск:
```bash
docker run -p 7540:7540 scheduler
```

С монтированием БД и всеми переменными окружения:
```bash
docker run -p 8080:8080 \
  -v $(pwd)/scheduler.db:/app/scheduler.db \
  -e TODO_PORT=8080 \
  -e TODO_DB=/app/scheduler.db \
  -e TODO_PASSWORD=secret \
  -e LOG_LEVEL=DEBUG \
  scheduler
```

Переменные окружения контейнера:

| Переменная | Описание | По умолчанию |
|----|----|----|
| `TODO_PORT` | Порт сервера | `7540` |
| `TODO_DB` | Путь к файлу БД внутри контейнера | `scheduler.db` |
| `TODO_PASSWORD` | Пароль для входа (если не задан --- аутентификация отключена) | --- |
| `LOG_LEVEL` | Уровень логирования (`DEBUG`, `INFO`, `WARN`, `ERROR`) | `INFO` |

> **Примечание:** При изменении `TODO_PORT` нужно также обновить маппинг портов (`-p <host>:<container>`). База данных монтируется с хоста через `-v`, чтобы данные сохранялись между перезапусками контейнера.

Открыть в браузере: http://localhost:7540
