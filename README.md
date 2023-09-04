# Кыш-кыш-мяу

Сервис алертинга и сбора метрик
----

Это учебная работе в [Яндекс-практикуме](https://practicum.yandex.ru), в рамках которой будет реализован сервер `мяу`, и агент `тыгыдык`, а может наоборот...


# Открытия

## `net/http`

Когда мы используем `ResponseWriter` при обработке машрутов, то после использования `Write()`, уже нельзя будет переписать заголовки с статус ответа. Хотя это в курсе не написано, там написано лишь о `WriteHeader()`, что после нее не имеет смысл менять заголовки, или использовать повторно `WriteHeader()`, чтобы записать изменения

`Write()` Действует жестче, после повторного применения, ни заголовки не статус не перезапишустся. Они останутся именно такими, как до первого применения. 

## Структура

Иногда хочется написать более абстрактный код, но его сложно планировать и сложно потом читать. Порой прямолинейно гораздо круче.(см. исследование /internal/store/README.md)

## Удаление из массива

Легко в массив что-то добавить, а вот с удалением какая-то куча проблем, включая там изменение исходного массива. Лучше алгоритмы стоить так, чтобы добавлять элементы