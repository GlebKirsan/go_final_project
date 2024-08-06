# Трекер задач

Программа с веб-интерфейсом для отслеживания задач. 
Можно:
- завести задачу, добавить ей описание и поставить период повторения;
- отметить задачу выполненной;
- удалить задачу;
- отредактировать задачу;
- искать по названию или дате.

## Команды
### Для запуска сервера
```bash
# запускать в корне
go run ./...
```
Либо с помощью утилиты `make`:
```bash
make run
```
По умолчанию сервер запускается на порту 7540.
Полный адрес `http://localhost:7540`.
Пароль для авторизации — `dummy`. (меняется через переменную окружения `TODO_PASSWORD`, секрет для шифровки jwt можно переопределить через переменную окружения `TODO_TOKEN_SECRET`)

Чтобы изменить порт на другой, потребуется определить переменную окружения `TODO_PORT`.

При запуске сервера, если файл с базой данных не будет найден, то программа сама создаст БД и таблицу с индексом.
### Для запуска тестов
```bash
# предварительно необходимо запустить сервер
go run ./...
go test ./tests
```
После прогона тестов БД будет заполнена тестовыми данными.

## Использованы библиотеки
- `migrate` для миграции БД;
- `zerolog` для логирования;
- `chi` фреймворк для разработки api;
- `envconfig` для наполнения конфига из переменных окружения;
- `render` для удобного создания JSON ответов;
- `sqlite` для взаимодействия с sqlite3;
- `jwt` для работы с токеном.

## Выполненные задачи со звёздочкой
- [x] Задание порта через `TODO_PORT`;
- [x] Задание пути до файла с базой данных через `TODO_DBFILE`;
- [x] Периоды с неделями и месяцами для задач;
- [x] Поиск по названию задачи или по дате;
- [x] Авторизация;
- [ ] Создание докер образа.