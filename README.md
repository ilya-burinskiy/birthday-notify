Запуск сервера:
```
SMTP_AUTH_USERNAME='email@example.com' SMTP_AUTH_PASSWORD='password' SMTP_HOST='smtp.gmail.com' SMTP_PORT='587' go run cmd/notifier/main.go
```

Регистрация:
```
curl -v -X POST 'http://localhost:8000/api/users/register' \
     -H "Content-Type: application/json" \
     -d '{"email": "email@example.com", "password": "pwd"}'
```

Аутентификация:
```
curl -v -X POST 'http://localhost:8000/api/users/login' \
     -H "Content-Type: application/json" \
     -d '{"email": "email@example.com", "password": "pwd"}'
```

Список пользователей:
```
curl -v -X GET 'http://localhost:8000/api/users' \
     -H "Content-Type: application/json"
```

Подписаться на пользователя c id равным {id}:
```
curl -v -X POST 'http://localhost:8000/api/users/{id}/subscribe' \
     -H "Content-Type: application/json" \
     --cookie jwt={your-jwt}
```

Отписаться от пользователя с id равным {id}:
```
curl -v -X DELETE 'http://localhost:8000/api/users/{id}/unsubscribe' \
     -H "Content-Type: application/json" \
     --cookie jwt={your-jwt}
```

Создать настройки для уведомлений:
```
curl -v -X POST 'http://localhost:8000/api/notify_settings' \
     -H "Content-Type: application/json" \
     --cookie jwt={your-jwt} \
     -d '{"days_before_notify": 1}'
```

Обновить настройки для уведомлений:
```
curl -v -X PATCH 'http://localhost:8000/api/notify_settings/{id}' \
     -H "Content-Type: application/json" \
     --cookie jwt={your-jwt} \
     -d '{"days_before_notify": 2}'
```
