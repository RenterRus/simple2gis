# Simple 2Gis
Упрощенная версия 2Гис

## Стек
1. Go 1.16.5
2. Cobra
3. SQLite

## Методы

| Метод       | Параметры                 | Что делает | Пример
| :---------------------: |:----------------------------:| :----:| :----:|
| /house     |  house <br> geo | Возвращает организации находящиеся в конкретном здании | http://127.0.0.1:9999/house?addr=test&geo=12.345,43.546 |
| /category   | category |   Возвращает организации по конкретной категории | http://127.0.0.1:9999/category?category=test,sub,auto |
| /organization  | id |   Возвращает организацию по ее ID | http://127.0.0.1:9999/organization?id=1 | 
| /newOrganization | houseAddr<br>houseGeo<br>numbers<br>category<br>name | Добавляет организацию в здание | http://127.0.0.1:9999/newOrganization?houseAddr=test&houseGeo=12.345,43.546&numbers=89990007766,76665554433&category=test,sub,auto&name=newTestName |
| /genFakeBase | number | Заполняет базу данными | http://127.0.0.1:9999/genFakeBase?number=4 |

## Инструкция по запуску
Не нуждается в дополнительном софте<br>
Если не получилось запустить, выполните в терминале проекта: <br>
1. go get -u github.com/spf13/cobra
2. go get github.com/mattn/go-sqlite3

<br>Скрипт по генерации БД лежит в scriptfordb.md<br>
БД лежит в guidebook <br>
### Запуск
В терминале с проектом: s2gis
Можно использовать флаги:
1. `--http` - устанавливает адрес прослушивания HTTP Сервера (по-умолчанию `127.0.0.1:9999`)
2. `--db` - устанавливает имя БД (по-умолчанию `guidebook`)

<br>Альтернатива: зеленая стрелочка в IDE (Goland) рядом с функцией main (cmd/orchestrator/main) (шутка)
