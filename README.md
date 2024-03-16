# Загаловок будет добавлен позже
Описание приложения будет добавлено позже

```bash
# Сборка и запуск
cd ...

# API Endpoint : <null>
```

## Структура
```
...
```

## API

#### /...
* `GET` : ...
* `POST` : ...

#### /projects/:title
* `GET` : ...
* `PUT` : ...
* `DELETE` : ...

## Todo
- [x] добавление информации об актёре (имя, пол, дата рождения)
- [x] изменение информации об актёре (Возможно изменить любую информацию об актёре, как частично, так и полностью)
- [x] ﻿﻿удаление информации об актёре
- [x] добавление информации о фильме (При добавлении фильма указываются его название (не менее 1 и не более 150 символов), описание (не более 1000 символов), дата выпуска, рейтинг (от 0 до 10) и список актёров)
- [x] изменение информации о фильме (Возможно изменить любую информацию о фильме, как частично, так и полностью)
- [x] ﻿﻿удаление информации о фильме
- [x] получение списка фильмов с возможностью сортировки по названию, по рейтингу, по дате выпуска (По умолчанию используется сортировка по рейтингу (по убыванию))
- [x] поиск фильма по фрагменту названия, по фрагменту имени актёра
- [x] получение списка актёров, для каждого актёра выдаётся также список фильмов с его участием
- [x] АРІ должен быть закрыт авторизацией (﻿﻿Поддерживаются две роли пользователей - обычный пользователь и администратор. Обычный пользователь имеет доступ только на получение данных и поиск, администратор - на все действия. Для упрощения можно считать, что соответствие пользователей и ролей задается вручную (например, напрямую через БД))
- [ ] разворот бэкенда на Docker'е
- [x] разворот базы данных Postgresql на Docker'е
- [ ] создание скрипта SQL для наполнения таблицами базу данных
- [ ] спецификация на АРІ (в формате Swagger 2.0 или OpenAPI 3.0)
      
## ДОПОЛНИТЕЛЬНО:
      
- [ ] использовать подход api-first (генерация кода из спецификации) или code-first (генерация спецификации из кода)
- [ ] добавить логирование - в лог должна попадать базовая информация об обрабатываемых запросах, ошибки
- [ ] покрыть код приложения юнит-тестами не менее чем на 70%
- [ ] создать Dockerfile для сборки образа
- [ ] создать docker-compose файл для запуска окружения с работающим приложением и СУБД
