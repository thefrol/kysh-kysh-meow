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

## Имена тестов

Тест называется через `нижнее почеркивание "_"`, чтобы не пересекаться с `docs`, там если тест имеет название `TestFoo()`, то он войдет в документацию как пример

## Тесты позволяют найти бесполезный код

Судя по покрытию, мои тесты никогда не заходят в одну область, если присмотрятся, то код туда никогда и не зайдет по архитектуре, просто сам себя проверяет. Это можно убрать. 

## Называть программу агентом

Да удобно говорить `http-клиент`, но когда дело касается работы в какой-то момент сложно понять о каком клиенте речь - о заказчике или программе.

> В описании клиента стояла задача покрасить кнопку в красный

Это человек так сказал, или просто в тз написано?

# А ещё
+ Хорошее название для программы - `metrica`