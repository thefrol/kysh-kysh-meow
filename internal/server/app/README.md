# Прикладной слой

Тут лежит то, что в моем понимании является юзкейсом

+ `manager` - это слой логики по метрикам, сохранение и удаление
+ `metricas` - тоже работа по метрикам, но специализированное для работы со струкрурой `Metrica`. Этакая обертка над `manager`. Возможно я бы это подложил бы в папку менеджер, но пока хочется выделить отдельно, чтобы было видно.
+ `dbping` - пингуем БД, отдельный юзкейс
+ `dashboard` - все для хтмл странички со всеми метриками