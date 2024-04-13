# Backend приложение для управления баннерами, написанное на GO
## Развертывание
### Общие настройки
Общие настройки приложения содержатся в [конфиге](https://github.com/deturpant/backend-trainee-banner-avito/blob/master/config/config.yaml). В зависимости от способа развертывания какие-либо параметры могут меняться.
В конфиге содержатся основные данные, необходимые для работы приложения.
Там же можно указать дефолтный админский пароль. Аккаунт администратора будет создан автоматически после первого поднятия приложения (логин admin).
Рекомендуется так же изменить jwt.secret ключ в целях безопасности.
### Docker
Для развертывания контейнера Docker создан [Dockerfile](https://github.com/deturpant/backend-trainee-banner-avito/blob/master/Dockerfile)
В [конфиге](https://github.com/deturpant/backend-trainee-banner-avito/blob/master/config/config.yaml) необходимо установить адрес 0.0.0.0:port.
Там же необходимо указать данные для подключения к БД PostgreSQL.
Для сборки image:
```
docker build -t <name> .
```
Для запуска контейнера:
```
docker run -d -p <внешний порт>:<внутренний порт> <name>
```
### Docker compose
Для развертывания в Docker Compose создан файл [docker-compose.yml](https://github.com/deturpant/backend-trainee-banner-avito/blob/master/docker-compose.yml)
В этом файле необходимо настроить порты (внешний и внутренний).
В [конфиге](https://github.com/deturpant/backend-trainee-banner-avito/blob/master/config/config.yaml) необходимо установить адрес 0.0.0.0:port
В настройках для подключения к БД хост необходимо указать: db
### Нативно
Для нативного запуска создан файл [build.sh](https://github.com/deturpant/backend-trainee-banner-avito/blob/master/build.sh). Он прокидывает путь к конфигу и запускает файл main.go
Предварительно, необходимо установить зависимости из [go.mod](https://github.com/deturpant/backend-trainee-banner-avito/blob/master/go.mod)
## Общее
Приложение представляет из себя сервис баннеров с функционалом регистрации, авторизации, добавления тэгов, фич, баннеров. 
С баннерами доступны следующие операции:
 - Создание (для админов)
 - Удаление (для админов)
 - Обновление (для админов)
 - Получение баннера для пользователя
 - Получения баннеров для админа
Для доступа к большинству функционала (кроме регистрации и авторизации) необходим доступ по токену.
Токен выдается пользователю после авторизации.
В дальнейшем токен должен передаваться вместе с заголовком запроса:
```
Authorization: Bearer <token>
```
## Документация
Для API сформирована [Swagger документация](https://github.com/deturpant/backend-trainee-banner-avito/tree/master/docs)
## Примеры запросов

Авторизация:
```
 curl -H 'Content-type: application/json' --data '{"name": "admin", "password" : "NJKfsjkdtoierf"}' http://localhost:8000/login
```
Регистрация:
```
curl -H 'Content-type: application/json' --data '{"name": "go", "password" : "123123"}' http://localhost:8000/users
```
Создание фичи:
```
curl -H 'Content-type: application/json' -H  'Authorization: Bearer <token>' --data '{"name": "go"}' http://localhost:8000/features
```
Создание тега:
```
curl -H 'Content-type: application/json' -H 'Authorization: Bearer <token>' --data '{"name": "go"}' http://localhost:8000/tags
```
Создание баннера:
```
curl -H 'Content-type: application/json' -H 'Authorization: Bearer <token>' --data '{"tag_ids": [0,1,2], "feature_id": 1, "content" : {"title":"some_text", "text" : "some_text", "url" : "some_url"}, "is_active" : true}' http://localhost:8000/banners
```
Удаление баннера:
```
curl -H 'Authorization: Bearer <token>' -X DELETE http://localhost:8000/banner/{id}
```
Получение пользовательского баннера:
```
curl -H 'Authorization: Bearer <token>' http://localhost:8000/user_banner?tag_id=2\&feature_id=12
curl -H 'Authorization: Bearer <token>' http://localhost:8000/user_banner?tag_id=2\&feature_id=12\&use_last_revision=true
```
Получение баннеров с фильтром:
```
curl -H 'Authorization: Bearer <token>' http://localhost:8000/banner?tag_id=2\&feature_id=12\&limit=1\&offset=1
```
