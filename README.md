# REST API для создания списков задач "TODO" на языке Go


# Технология реализации приложения
- Приложение написано следуя принципу Чистой Архитектуры и техники внедрения зависимостей.
- В качестве web феймворка используется <a href="https://github.com/gin-gonic/gin">gin-gonic/gin</a>.
- В качестве базы данных используется postgreSQL(Запуск через docker).
- Конфигурация приложения осуществляется с помощью библиотеки <a href="https://github.com/spf13/viper">spf13/viper</a>.
- Работа с БД осуществляется, используя библиотеку <a href="https://github.com/jmoiron/sqlx">sqlx</a>.
- Регистрация и аутентификация реализована с помощью JWT <a href="https://github.com/golang-jwt/jwt">golang-jwt</a>.


### Для запуска приложения:

```
docker-compose up --build todoapp && docker-compose up todoapp
```


Если приложение запускается впервые, необходимо применить миграции к базе данных:

```
migrate -path ./schema -database 'postgres://yourdbusername:yourpassword@host:port/todoapp?sslmode=disable' up
```
Замените yourdbusername, yourpassword, host, port нужными вам параметрами

Параметры подключения к базе данных находятся в .\configs\config.yml