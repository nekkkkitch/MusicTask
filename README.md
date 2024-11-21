# Installation and launching
1. Клонируйте проект куда-нибудь
   ```
   git clone https://github.com/nekkkkitch/MusicTask
   ```
2. Перейдите в папку MusicTask
   ```
   cd .../MusicTask
   ```
3. Запустите в терминале следующую команду(Должен быть запущенный Docker engine и Makefile должен быть установлен)
   ```
   make buildbuilder
   ```
4. Запустите в терминале ещё одну команду(100% сначала контейнер с gateway остановится из-за того, что БД запускается довольно долго. Подождите и потом запустите контейнер с gateway)
    ```
   make start
   ```
P.S. Если что-то не работает - напишите мне пожалуйста на телеграм @nekkkkitch
