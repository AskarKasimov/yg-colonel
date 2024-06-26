# Распределенный вычислитель арифметических выражений
Забегая вперёд, это репозиторий оркестратора, а репозиторий с агентом вот: https://github.com/AskarKasimov/yg-recruit

Связь со мной: [Telegram](https://t.me/a_s_k_a_rr), [прочее тут](https://askar.su)

## Архитектура
Система работает по системе поллинга.

Агенты (коих я назвал **Recruit**) с определенным интервалом посылают GET-запросы Оркестратору (коего я назвал **Colonel**).
Colonel имеет связь с базой данных на PostgreSQL (поля продемонстрированы на картинке ниже).

Названия полей базы прямо отражают их суть.

![Иллюстрация](https://github.com/AskarKasimov/yg-colonel/blob/master/scheme.drawio.png)

P.S. поле *progress* в таблице *Expressions* имеет три возможных состояния: done, processing и waiting, соотвественно, когда выражение обработано, когда в процессе счета, когда ожидает своей очереди на исполнение.

С определенным временным интервалом Colonel также проходится по Recruit'ам, обновляя их состояние согласно последнему *Heartbeat'у*. Если выясняется, что какой-то из агентов не выходил на связь в последнее время, его выражения снова переходят в статус ожидающих решения и в ближайшее время снова найдут своего решателя)).

При включении Recruit *регистрируется* у оркестратора, получая свой ID в обмен на *Имя*. Имя – своеобразный идентификатор, который берется из файла ./conf/conf.json (намеренно включён в .gitignore; при запуске всегда создаёт новое *Имя*, если не находит старого). Это обеспечивает **бесперебойную** работу агентов.

Таким образом обеспечена отказоустойчивость. 

При отправке запроса на получение выражения для его счёта сервер обновляет состояние *isAlive* Recruit'а и время последнего heartbeat'а, чтобы списать неактуальных агентов даже после длительного периода отключения оркестратора.

## Счёт

Получая выражение, Recruit использует библиотеку для вычисления, а затем *спит* определенное время в соответствии с заданной сложностью арифметических операций (**задаются в переменных окружения в секундах**).

## Авторизация

Чтобы создать юзера, нужно отправить запрос на /api/v1/auth/register с JSON'ом:
```
{
  "login":"",
  "password":""
}
```
После получения ответа **"OK"** можно отправляться на эндпоинт /api/v1/auth/login с таким же JSON'ом, как и на этапе регистрации. В ответ будет выдан токен в следующем формате:
```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImUzN2E0NDY1LTJiMmItNDQxMy05MmZlLTQwNmRkZmUxMWMxZiIsImV4cCI6MTcxMzgxMzUxNX0.BTff50lGAuYq-r8jdbpHigpV8iviMiBt0skwnCjQHRw"
}
```
Его необходимо далее использовать при **добавлении** и **просмотре** выражений. У юзеров есть доступ только к своим выражениям.

Чтобы два этих метода работали корректно, нужно вставить в header HTTP-запроса следующее:

```
Authorization: Bearer ТОКЕН
```

*P.S. Оставлен работающим открыто метод просмотра всех существующих в базе выражений, чтобы можно было смотреть их через фронтенд. Ограничение на авторизацию работает при просмотре конкретного выражения по ID. Это допущенная мною условность для вашего же удобства)*

## Добавление и просмотр выражений

Для добавления отправляем POST-запрос на /api/v1/expression/add со следующим JSON'ом:
```
{
    "expression": "2+2"
}
```
В ответ получаем UUID созданного выражения.

Далее используем его для получения актуальной информации. Для этого шлём GET-запрос на /api/v1/expression/**выданный UUID**.
Если **авторизация проходит успешно**, и выражение было добавлено **из-под текущей учётной записи**, увидим следующее:
```
{
    "id": "2c2eab10-5fd6-47bd-9596-b8e9a02958ce",
    "incomingDate": 1713726698,
    "vanilla": "2+2",
    "answer": "4",
    "progress": "done"
}
```

Параллельно с этим в корне порта 8080 мы можем увидеть **очень минималистичную** сводку работы системы со всей актуальной информацией (без авторизации, как ранее и было замечено))


## Запуск

Необходимо запустить один единственный Colonel, а также неограниченное количество Recruit.

Разработанная мною система **(согласно легенде)** предполагает запуск этих частей проекта на разных устройствах, НО будет работать и на одном.

(!) Единственный нюанс – Recruit'ов сделайте несколько копий и запускайте docker-compose из разных директорий, чтобы сервер различал их как разных. ИНАЧЕ СЕРВЕР ОПРЕДЕЛИТ ИХ ВСЕХ ЗАПУЩЕННЫМИ ПОД ОДНИМ ИМЕНЕМ.

Используем Docker Compose из корневой директории этого репозитория:
```
docker-compose up --build
```
Далее скачиваем [репозиторий Recruit](https://github.com/AskarKasimov/yg-recruit) (каждый желаемый экземпляр нужно склонить с гитхаба отдельно, о чем я уже написал выше))) и тоже запускаем его с помощью Docker Compose:
```
docker-compose up --build
```
P.S. При желании можно не уделять отдельное внимание настройке переменных окружения, ведь они изначально неплохо подготовлены)

P.P.S. описаны они в docker-compose.yaml, и их названия полностью отражают их суть

Подготовлен Frontend на http://localhost:8080 (с авторизацией немного потерял свой смысл, рекомендую использовать Postman), а также **есть и Swagger на http://localhost:8080/swagger/index.html**

## Тестовый сценарий:

1. Регистрация
2. Вход
3. Добавление выражения
4. Ожидание (по умолчанию оно выставлено небольшим)
5. Проверка

Всё нужное для этого подробно описано выше

Подведу небольшой итог:
- [x] Авторизация по JWT и настройка приватности просмотра и создания выражений
- [x] Всё хранится в базе PostgreSQL, обеспечена полная отказоустойчивость
- [ ] gRPC
- [ ] Тесты

В будущем планирую разделить логику оркестратора на микросервисы с gRPC и экранировать внутренние эндпоинты от пользователя, а также выделить frontend на React.
