# Тестовое задание для стажёра Backend

## Решение

### API

Дополнил [API](api.yaml), новым функционалом.

### Запуск

В репозитории есть *Makefile*:
```
    make buildrun   // чтобы забилдить и поднять бд+само приложение
    make stop       // остановить бд+само приложение
    make genMock    // для генерации замоканых repository(для тестов)
    make test       // запуск тестов
    make linter     // запуск линтера
```

Подробнее про тесты [тут](#тест).

###### Подключение к PostgrsQL(psql):
```
psql postgresql://root:root@localhost:5432/taskdb
```

###### Запуск нагрузочного теста(YandexTank):
```
docker run     -v $(pwd):/var/loadtest     -v $SSH_AUTH_SOCK:/ssh-agent -e SSH_AUTH_SOCK=/ssh-agent     --net host     -it yandex/yandex-tank
```

_________________________________________________
### Проблемы

1. Организация проверки существования пары (feature_id, tag_id).

Первое что пришло в голову хранить по отдельности пары, те когда приходит запрос на создание записи с n tag-ами,
то будет запись n строк. По хорошему я бы посмотрел каких запросов больше на конкретный баннер или на баннер по 1 параметру. 
Если бы первых было бы больше, то наверное это норм решение, но было бы очень много записей. Поэтому я решил, что лучше хранить прям массив tag-ов.
И для проверки использовать intersect.

UDP: Понял что лучше надо было иметь две таблицы. Одну с баннером, а другую со связями: {tagID, featureID, bannerID}.
Тогда бы не было повторяющейся информации, а при запросах поиска по какому-то определенному featureID/bannerID 
надо было бы сделать 2 индекса для featureID, tagID, чтобы не было поиска по всей таблице.

2. Написание кэша и его устройство

Чтобы быстро обрабатывать запросы я сделал 3 кэша для 3 видов запросов: по паре{featureID, bannerID}, по featureID, по bannerID.
Устройство: есть горутина, которая каждые 5 минут идет в БД, и берет от туда всю инфу, после чего создает 3 мапы(кэша), берет mutex на запись и заменяет новыми мапами.

### Общие детали разработки
- Для общения использовал web framework - [Fiber](https://github.com/gofiber/fiber).
- Для SQL запросов использовал библиотеку [SQLX](https://github.com/jmoiron/sqlx).
- Для авторизации использовал [JWT-token](https://github.com/dgrijalva/jwt-go)(access) и просто генерировал рандомную строку(refresh). 
В JWT-token внутрь зашил роль клиента и его ID, так что при изменении роли нужно будет перелогиниваться(получить новый access). 
Чтобы изменить роль с обычного user на admin необходимо в ручную поменять в БД колонку *role* 0 на 1. Обработка access token-a в Middleware.

_________________________________________________

### Дизайн системы
![Design System](./picture/designSystem.jpg)

_________________________________________________
### Тест

Мокаем repository и делаем самый обычный интергационный тест:
- SignIn пользоваетеля
- Берем content из БД 
- Проверяем что все совпадает

_________________________________________________
### Дополнительные задания 

#### Адаптировать систему для значительного увеличения количества тегов и фичей

Не придумал ничего умнее кроме как сделать, lru или lfu cache. Идея в том, что самые частые запросы мы храним в cache,
за нечастными запросами идем в БД. Но тогда нужно будет online менять данные в бд. И тогда пропадает свойство того, что 
5 минут может быть старая информация, у нас будет всегда новая информация. Но это тяжелее содержать, могут быть проблемы с когерентностью и тд.

#### Нагрузочное тестирование

###### GeneratorData
Для чистоты нагрузочного тестирования нужно было накатить данных в БД. Чтобы это сделать необходимо:
- настроить config для генератора
- удалить все данные из таблицы *Banner*. Генератор не проверяет есть ли в БД такой {featureId, tagId} или нет, если настроить правильно конфиг, то базу можно не удалять

Подробнее про *config*:
Чтобы накатить новых данных, во-первых, необходимо в конфиге сервера сделать поле *WithGenerator: true*

Во-вторых, сам *config*:

```
  cntRow: 500        // максимальное количество записей
  cntFeature: 100    // количество разных featureID
  cntTag: 100        // количество разных tagID
  maxTagInRow: 5     // максимальное количество tagID в одной записи
```

###### YandexTank  
После генерации данных, надо настроить *load.yaml* - конфиг для теста. Посмотреть его [тут](./yandex_tank/load.yaml).

UDP: Да, в данном тесте нет запросов update, delete, create; так как вероятнее всего сервис будет иметь +95% GET-запросов 
и те несколько процентов большую роль не сыграют. 

##### Результат тестирования

Результат можно посмотреть по ссылке на [overload](https://overload.yandex.net/693232#tab=test_data&tags=&plot_groups=main&machines=&metrics=&slider_start=1712955941&slider_end=1712956541).

#### Версионирование баннеров

Решил сделать что-то на подобии журнала обновлений: как только делается какое-то изменение данных мы записываем это изменение в журнал обновлений.
Дальше можем посмотреть по banner_id его версии и выбрать его. Создал новую базу данных vbanner{version_id, banner_id, data}(храню data в байтах нужный JSON).
Добавил 3 новых ручки:
```
GET     /version/:id    - дает информацию о версия определенного banner_id
DELETE  /version/:vid   - удаляет определенную версию version_id
PATCH   /version        - берет определеную версию и выкатывает её в акутальную (в query передается banner_id и version_id)
```

#### Дополнительные тесты

Для доп тестов вместо Mock-ов написал фэйковую реализацию БД просто через map.
- Проверка всех get запросов с авторизацией.
- Проверка create с авторизацией
- проверка delete с авторизацией

#### Конфигурация линтера

Чтобы запустить линтер: ```make linter``` 

Config linter-а можно посмотреть [тут](./linter/.golangci.yml).

_________________________________________________
_________________________________________________
_________________________________________________

## Задание

### Сервис баннеров
В Авито есть большое количество неоднородного контента, для которого необходимо иметь единую систему управления.  В частности, необходимо показывать разный контент пользователям в зависимости от их принадлежности к какой-либо группе. Данный контент мы будем предоставлять с помощью баннеров.
### Описание задачи
Необходимо реализовать сервис, который позволяет показывать пользователям баннеры, в зависимости от требуемой фичи и тега пользователя, а также управлять баннерами и связанными с ними тегами и фичами.
### Общие вводные
**Баннер** — это документ, описывающий какой-либо элемент пользовательского интерфейса. Технически баннер представляет собой  JSON-документ неопределенной структуры. 
**Тег** — это сущность для обозначения группы пользователей; представляет собой число (ID тега). 
**Фича** — это домен или функциональность; представляет собой число (ID фичи).  
1. Один баннер может быть связан только с одной фичей и несколькими тегами
2. При этом один тег, как и одна фича, могут принадлежать разным баннерам одновременно
3. Фича и тег однозначно определяют баннер

Так как баннеры являются для пользователя вспомогательным функционалом, допускается, если пользователь в течение короткого срока будет получать устаревшую информацию.  При этом существует часть пользователей (порядка 10%), которым обязательно получать самую актуальную информацию. Для таких пользователей нужно предусмотреть механизм получения информации напрямую из БД.
### Условия
1. Используйте этот [API](https://github.com/avito-tech/backend-trainee-assignment-2024/blob/main/api.yaml)
2. Тегов и фичей небольшое количество (до 1000), RPS — 1k, SLI времени ответа — 50 мс, SLI успешности ответа — 99.99%
3. Для авторизации доступов должны использоваться 2 вида токенов: пользовательский и админский.  Получение баннера может происходить с помощью пользовательского или админского токена, а все остальные действия могут выполняться только с помощью админского токена.  
4. Реализуйте интеграционный или E2E-тест на сценарий получения баннера.
5. Если при получении баннера передан флаг use_last_revision, необходимо отдавать самую актуальную информацию.  В ином случае допускается передача информации, которая была актуальна 5 минут назад.
6. Баннеры могут быть временно выключены. Если баннер выключен, то обычные пользователи не должны его получать, при этом админы должны иметь к нему доступ.

### Дополнительные задания:
Эти задания не являются обязательными, но выполнение всех или части из них даст вам преимущество перед другими кандидатами. 
1. Адаптировать систему для значительного увеличения количества тегов и фичей, при котором допускается увеличение времени исполнения по редко запрашиваемым тегам и фичам
2. Провести нагрузочное тестирование полученного решения и приложить результаты тестирования к решению
3. Иногда получается так, что необходимо вернуться к одной из трех предыдущих версий баннера в связи с найденной ошибкой в логике, тексте и т.д.  Измените API таким образом, чтобы можно было просмотреть существующие версии баннера и выбрать подходящую версию
4. Добавить метод удаления баннеров по фиче или тегу, время ответа которого не должно превышать 100 мс, независимо от количества баннеров.  В связи с небольшим временем ответа метода, рекомендуется ознакомиться с механизмом выполнения отложенных действий 
5. Реализовать интеграционное или E2E-тестирование для остальных сценариев
6. Описать конфигурацию линтера

### Требования по стеку
- **Язык сервиса:** предпочтительным будет Go, при этом вы можете выбрать любой, удобный вам. 
- **База данных:** предпочтительной будет PostgreSQL, при этом вы можете выбрать любую, удобную вам. 
- Для **деплоя зависимостей и самого сервиса** рекомендуется использовать Docker и Docker Compose.
### Ход решения
Если у вас возникнут вопросы по заданию, ответы на которые вы не найдете в описанных «Условиях», то вы вольны принимать решения самостоятельно.  
В таком случае приложите к проекту README-файл, в котором будет список вопросов и пояснения о том, как вы решили проблему и почему именно выбранным вами способом.
### Оформление решения
Необходимо предоставить публичный git-репозиторий на любом публичном хосте (GitHub / GitLab / etc), содержащий в master/main ветке: 
1. Код сервиса
2. Makefile c командами сборки проекта / Описанная в README.md инструкция по запуску
3. Описанные в README.md вопросы/проблемы, с которыми столкнулись,  и ваша логика их решений (если требуется)
