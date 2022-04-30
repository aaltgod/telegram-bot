### @SocialMediaTest_bot

### Описание
- `bot` - телеграм бот
- `api` - апи c эндпоинтами
    - /api/get_users (GET)
    - /api/get_user?id= (GET)
    - /api/get_history_by_tg?id= (GET)
    - /api/delete_record?id= &ip= (DELETE)
- `storage-service` - сервис, позволяющий взаимодействовать с базой данных посредством вызовов эндпоинтов
  - /users/:id (GET)
  - /users/:id (PUT)
  - /users/ (POST)
  - /users (GET)
  - /requests/:user_id (POST)
  - /requests (DELETE)
  - /requests/:user_id (GET)

### Подготовка для запуска
Для корректной работы требуется админ. Нужно добавить имя и ID пользователя телеграма в `newadmin.yaml`.

### Запуск

`docker-compose up --build -d`
