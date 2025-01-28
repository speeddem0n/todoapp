# REST API для создания списков задач "TODO" на языке Go

# Основные возможности
- Регистрация и авторизация пользователей (JWT токены)
- Управление списками задач:
    - Создание, чтение, обновление и удаление (CRUD)
- Управление задачами внутри списков:
    - Создание, чтение, обновление и удаление (CRUD)
- Поддержка Swagger для документации API

# Технология реализации приложения
- В качестве web феймворка используется <a href="https://github.com/gin-gonic/gin">gin-gonic/gin</a>.
- В качестве базы данных используется postgreSQL(Запуск через docker).
- Конфигурация приложения осуществляется с помощью библиотеки <a href="https://github.com/spf13/viper">spf13/viper</a>.
- Работа с БД осуществляется, используя библиотеку <a href="https://github.com/jmoiron/sqlx">sqlx</a>.
- Регистрация и аутентификация реализована с помощью JWT <a href="https://github.com/golang-jwt/jwt">golang-jwt</a>.
- Документация API с помошью <a href="https://github.com/swaggo/swag">swagger</a>.


### Для запуска приложения:

1. Создайте файл .env в корневой директории проекта и укажите следующие параметры:
```dotenv
DB_PASSWORD=postgres
```
2. При необходимости изменить настройки приложения в файле internal/configs/config.yml

3. Компиляция
```
make build
```

4. Приминение миграций
```
make migrateup
```

5. Запуск
```
make run
```

### Список доступных комманд: 

```
make help
```

### Документация Swagger:

http://localhost:8000/swagger/index.html

