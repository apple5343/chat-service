# User Service API Documentation

## Base URL
`http://localhost:8090/api`

---

### 1. Регистрация пользователя

**Endpoint:**  
`POST /register`

**Описание:**  
Регистрирует нового пользователя в системе.

**Тело запроса:**
```json
{
    "email": "example@domain.com",
    "password": "password123",
    "grade": "Senior"
}
```

**Ответ:**
```json
{
    "id": "a33c66f5-96d2-46cf-a8cc-461ed4ab5f15"
}
```

---

### 2. Авторизация пользователя

**Endpoint:**  
`POST /login`

**Описание:**  
Авторизует пользователя в системе и возвращает refresh_token.

**Тело запроса:**
```json
{
  "email": "example@domain.com",
  "password": "password123"
}
```

**Ответ:**
```json
{
  "refresh_token": "eyJhbGci0iJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### 3. Получение нового refresh_token

**Endpoint:**  
`POST /get_refresh_token`

**Описание:**  
Обновляет refresh_token.

**Тело запроса:**
```json
{
  "refresh_token": "eyJhbGci0iJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ответ:**
```json
{
  "refresh_token": "eyJhbGci0iJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### 4. Получение access_token

**Endpoint:**  
`POST /get_access_token`

**Описание:**
Генерирует новый access_token для авторизованного пользователя.

**Тело запроса:**
```json
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ответ:**
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---


### 5. Получение информации о пользователе

**Endpoint:**  
`GET /users/{user_id}`

**Описание:**
Возвращает информацию о пользователе по его id. Только для авторизованных пользователей.

**Ответ:**
```json
{
  "id": "user123",
  "email": "example@domain.com",
  "grade": "Senior",
  "created_at": "2024-06-01T12:00:00Z"
}
```

---

### 6. Обновление пользователя

**Endpoint:**  
`PATCH /users/{user_id}`

**Описание:**
Обновляет данные пользователя.

**Тело запроса:**
```json
{
  "grade": "Updated Grade"
}
```

**Ответ:**
```json
{
  "message": "User updated successfully"
}
```

---

### 7. Удаление пользователя

**Endpoint:**  
`DELETE /users/{user_id}`

**Описание:**
Удаляет пользователя из системы.

**Ответ:**
```json
{
  "message": "User deleted successfully"
}
```

---