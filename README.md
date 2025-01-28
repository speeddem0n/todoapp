# REST API для создания списков задач "TODO" на языке Go


# Технология реализации приложения
- Приложение написано следуя принципу Чистой Архитектуры.
- В качестве web феймворка используется <a href="https://github.com/gin-gonic/gin">gin-gonic/gin</a>.
- В качестве базы данных используется postgreSQL(Запуск через docker).
- Конфигурация приложения осуществляется с помощью библиотеки <a href="https://github.com/spf13/viper">spf13/viper</a>.
- Работа с БД осуществляется, используя библиотеку <a href="https://github.com/jmoiron/sqlx">sqlx</a>.
- Регистрация и аутентификация реализована с помощью JWT <a href="https://github.com/golang-jwt/jwt">golang-jwt</a>.
- Документация API с помошью <a href="https://github.com/swaggo/swag">swagger</a>.


### Для запуска приложения:

```
make build
```

```
make run
```

### Список доступных комманд: 

```
make help
```