MAPS

**Задача**  
Задача является продолжением задачи movies.
Требуется добавить команду для загрузки фильмов:  
load d <path_to_directory>

При выполнении команды фильмы загружаются из файлов в память.
Путь до дирректории указывается в path_to_directory

Условия:
* Данные загружаются асинхронно, количество go рутин на усмотрение выполняющего задание
* Данные загружаются в map, можно пользоваться всем из пакета sync