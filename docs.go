// Copyright 2020
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// # Сервис сбора метрик и алертинга
//
// Используется, чтобы собирать метрики с удаленного
// хоста и хранить их: базе данных, файле или памяти.
// В отличие от прометеуса, использует push модель,
// сбора метрик, где метрики отправляются на сервер при помощи
// агента, а сервер лишь обрабатывает эти запросы, при этом
// сам агентом не опрашивает
//
// # Сервер
//
// Запускается командой
//
//	go run github.com/thefrol/kysh-kysh-meow/cmd/server
//
// Флаг "-р" позволяет посмотреть настройки запуска
//
//	-a string
//	      [адрес:порт] устанавливает адрес сервера  (default ":8080")
//	-d value
//	      [строка] подключения к базе данных
//	-f string
//	      [строка] путь к файлу, откуда будут читаться при запуске и куда будут сохраняться метрики полученные сервером, если файл пустой, то сохранение будет отменено (default "/tmp/metrics-db.json")
//	-i uint
//	      [время, сек] интервал сохранения показаний. При 0 запись делается почти синхронно (default 300)
//	-k value
//	      строка, секретный ключ подписи
//	-profile
//	      [флаг] создает хендлеры для профирирования
//	-r
//	      [флаг] если установлен, загружает из файла ранее записанные метрики (default true)
//
// # Агент
//
// Позволяет собирать базовые метрики с удаленной машины, такие как использование памяти
// и процессора. Отправляет метрики на сервер Кыш-мяу
//
// Запускается командой
//
//	go run github.com/thefrol/kysh-kysh-meow/cmd/agent
//
// Флаг "-a" позволяет указать адрес сервера кыш-мяу,
// где будут храниться метрики.  Полный перечень настроек
// можно посмотреть с флагом "-h".
package kyshkyshmeow
