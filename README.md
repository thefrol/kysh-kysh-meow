# Кыш-кыш-мяу

![покрытие_тестами](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/thefrol/571cee70d9542f9f39127a69ed007706/raw/main.json)

Сервис алертинга и сбора метрик

----

Основная задача этого продукта: показать какой я классный сотрудник и как классно я подхожу на должность, на которую вы меня сейчас рассматриваете. Поэтому тут вас ждет:

+ Фантастическая скорость работы достигается статической типизацией, использованием стека
+ Сладкая, как `торт наполеон`, и слоистая архитектура. Все зависимости идут снизу вверх. И нет циклических зависимостей. `DDD` во все поля
+ Защищает связь сервера и агента при помощи подписей `SHA256`
+ Сжимает запросы и ответы при помощи `gzip`
+ Превосходное оформление коммитов, ишью и пул-реквестов, прохождение код-ревью. Гитхаб использован как только возможно.
+ Все пакеты и функции документированы
+ Модульно, предметно ориентировано
+ Более ста тестов

## Стек

+ Go: chi, http, easyjson,
+ gzip, encoding, crypto,
+ runtime, sync
+ zerolog

Это учебная работе в [Яндекс-практикуме](https://practicum.yandex.ru), поэтому такие вещи как сжатие ответа и запроса, подписывание реализовано ручками при помощи стандартных библиотек.

## Вещи которыми я горжусь

### `config.Secret`, `config.ConnectionString`

Эта структура для чувствительных данных, которая защищена от попадания в консоль. Если попытаться распечатать конфиг, то мы увидим

SecretValue:******
DataBaseString: host=localhost user=user ... // хотя тут был пароль

## Открытия

### `net/http`

Когда мы используем `http.ResponseWriter` все заголовки и код ответа надо успеть записать до первого использования `Write()`. 

### Структура

Иногда хочется написать более абстрактный код, но его сложно планировать и сложно потом читать. Порой прямолинейно гораздо круче.(см. исследование /internal/store/README.md)

### Удаление из массива

Легко в массив что-то добавить, а вот с удалением какая-то куча проблем, включая там изменение исходного массива. Лучше алгоритмы стоить так, чтобы добавлять элементы

### Имена тестов

Тест называется через `нижнее почеркивание "_"`, чтобы не пересекаться с `docs`, там если тест имеет название `TestFoo()`, то он войдет в документацию как пример

### Тесты позволяют найти бесполезный код

Судя по покрытию, мои тесты никогда не заходят в одну область, если присмотрятся, то код туда никогда и не зайдет по архитектуре, просто сам себя проверяет. Это можно убрать.

### Называть программу агентом

Да удобно говорить `http-клиент`, но когда дело касается работы в какой-то момент сложно понять о каком клиенте речь - о заказчике или программе.

> В описании клиента стояла задача покрасить кнопку в красный

Это человек так сказал, или просто в тз написано?

### `http.NotFound()`

Требует после себя return или продолжится выполнение функции

### `http.NoBody`

Отличнай замена `nil` в запросах, где не требудется тело

### Проблема архитектуры

Что мне не нравится в архитектуре сейчас, что отдельные блоки сильно зависят от других. Так, например, функция отправки агента должна точно знать какие метрики существуют, чтобы их отправлять. То есть добавь новую метрику, и нужно будет бегать по всему коду проверять, чтобы все нормально работало.

## Рефакторинг может и замусорить

Отвязывал stats от хранилища, потом слил с другим пакетом, и получился модуль-мусорная корзина, пусть довольно уверенный. Работает классно, выглядит ужасно. Не всегда рефакторинг делает общее пространство чище.

## Выделяй побольше пакетов из программ

Это очень помогает, когда хочешь настроить зависимости одних классов или компонентов от других. Пока все лежит в пакете `main`, мы даже на уровне кода не можем прописать зависимости. Мы не можем ссылаться на пакет `main` из других пакетов. А когда такая возможность появляется, голова как бы сразу проясняется. Да, это порой бывает непросто. Тут же вознимают проблемы с глобальными переменными и прочим, но но того стоит

## git rm --cached

Эта команда удаляет в репозитории, но оставляет на диске. Потом надо просто добавить в `.gitignore`

## Не стоит слишком сильно оптимизировать по мелочам

Порой оптимизации просто разбиваются о планировщик.

Короче я провел исследование. Померил скорость работы и оказывается, что, если мы не будем лишний раз проверять типы данных то это дам сэкономик несколько наносекунд. Короче эта штука будет иметь смысл только, когда количество запросов перевалит за миллиард. Тогда это станет разница в две секунды.

В реальных условиях, даже 1000 запусков может то тут то там быстрее работать, видимо из-за особенностей запусков контекстов и переключений между потоками, так что **не стоит тратить свою жизнь на оптимизацию нескольких наносекунд**, которые и то игрушка ~~дьявола~~ планировщика потоков.

## PiePie - сильные стороны

Сейчас, когда понятно, что бизнес логику, возможно стоит поменять основной образ приложения на пирог, пирожное, может бутерброд. Даже добавить это в ридми.

Вообще, мне кажется, это хорошая идея выделить какие-то сильные стороны своего продукта. Например, пироговую архитектуры, где метрики собираются как пирог. Щепотка обработки ошибок, корж из джейсона, и крем из DTO.

Мне даже интересно, станет ли программа и код понятнее, если я буду использовать пироговую идиому. Мне кажется что да. Потому что идеома с котами, как-то не добавляет понятности, потому что никак не связана с кодом

## Обнулить слайс

Если мы использовали слайс как буфер, и хотим его обнулить, вместо того, чтобы выделять новое место под него, можно поспользоваться `arr=arr[:0]`, при этом мы сохраним капасити слайса и место в памяти под него выделенное

## Собрать кучу ошибок в одной

Я не знал, но с этим отлично справляется `errors.Join()`

## Временные комменты

Когда добавляешь новый код, или поле в структуру, можно написать над этим полем комментарий. Потом его можно будет стереть, но мне кажется прикольно, что в коммите навсегда останется намерение с котором было сделано изменение.

## Свежесть и движение

Залип над кодом, и не знаешь что сделать дальше? Подумай, а не сделал ли ты уже фичу, может быть пора закрыть ПР, или создать новоую ветвь.

Вообще, процесс должен быть примерно таким. Идея - новая ветка - ПР - документация более крупного ПР

Тогда не будешь застревать, маленьких инкременты заставляют думать голову. А что я делаю сейчас. А что надо? А как я это опишу?

## Читаемость

Иногда код начинает походить на черт-знает-что

```go
// сборщик мемстатс выделен в отдельный планировщик
inMs := generator(context.TODO(), fetch.MemStats, time.Second*time.Duration(config.PollingInterval))
_ = generator(context.TODO(), fetch.PollCount, time.Second*time.Duration(config.PollingInterval))
_ = generator(context.TODO(), fetch.RandomValue, time.Second*time.Duration(config.PollingInterval))
```

Go требует особого подхода к стилю. Тут нельзя указать имена аргументов для функции, но это не проблема для того, чтобы создать читаемый код

```go
// создать каналы сбора метрик
interval := time.Second * time.Duration(config.PollingInterval)
ctx := context.TODO()
inMs := generator(ctx, fetch.MemStats, interval)
_ = generator(ctx, fetch.PollCount, interval)
_ = generator(ctx, fetch.RandomValue, interval)
```

Ну совсем другое дело! И выглядит по-гошному и читается легко. Параметры в хорошо названных переменных, комментарий сверху и пара пустых строк сделают этот код прекрасным гошным великолепием.

## А ещё

+ Хорошее название для программы - `metrica`
