# 🔐 JWT Authentication Guide

Руководство по использованию JWT аутентификации в HugoProxy API.

## 📋 Обзор

API использует JWT (JSON Web Tokens) для аутентификации пользователей. Защищенные endpoints требуют валидный JWT токен в заголовке Authorization.

## 🚀 Быстрый старт

### 1. Регистрация нового пользователя

**Endpoint:** `POST /api/register`

**Запрос:**
```json
{
  "username": "user@example.com",
  "password": "password123"
}
```

**Ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20ifQ..."
}
```

**PowerShell:**
```powershell
$body = @{
    username = "user@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/register" `
  -Method Post `
  -ContentType "application/json" `
  -Body $body

$token = $response.token
Write-Host "Token: $token"
```

**cURL:**
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user@example.com","password":"password123"}'
```

### 2. Вход в систему

**Endpoint:** `POST /api/login`

**Запрос:**
```json
{
  "username": "user@example.com",
  "password": "password123"
}
```

**Ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20ifQ..."
}
```

**PowerShell:**
```powershell
$body = @{
    username = "user@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/login" `
  -Method Post `
  -ContentType "application/json" `
  -Body $body

$token = $response.token
```

### 3. Использование защищенных endpoints

После получения токена, добавьте его в заголовок `Authorization` со значением `Bearer <token>`.

**Поиск адреса:**

**PowerShell:**
```powershell
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

$body = @{
    query = "Москва, Красная площадь, 1"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/address/search" `
  -Method Post `
  -Headers $headers `
  -Body $body
```

**cURL:**
```bash
curl -X POST http://localhost:8080/api/address/search \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{"query":"Москва, Красная площадь, 1"}'
```

**Геокодирование:**

**PowerShell:**
```powershell
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

$body = @{
    lat = "55.753215"
    lng = "37.622504"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/address/geocode" `
  -Method Post `
  -Headers $headers `
  -Body $body
```

## 🔑 Тестирование в Swagger UI

1. **Откройте Swagger UI:** http://localhost:8080/swagger/

2. **Зарегистрируйтесь или войдите:**
   - Раскройте endpoint `/api/register` или `/api/login`
   - Нажмите "Try it out"
   - Введите username и password
   - Нажмите "Execute"
   - Скопируйте полученный токен из ответа

3. **Авторизуйтесь в Swagger:**
   - Нажмите кнопку **"Authorize"** (замок) вверху страницы
   - В поле "Value" введите: `Bearer <ваш_токен>`
   - Нажмите "Authorize"
   - Закройте диалоговое окно

4. **Используйте защищенные endpoints:**
   - Теперь вы можете тестировать `/api/address/search` и `/api/address/geocode`
   - Токен будет автоматически добавляться к запросам

## 📝 Структура токена

JWT токен содержит следующие claims:

```json
{
  "email": "user@example.com"
}
```

В коде можно получить данные пользователя:

```go
_, claims, _ := jwtauth.FromContext(r.Context())
email := claims["email"]
```

## 🔒 Endpoints

### Публичные (не требуют аутентификации):
- `POST /api/register` - Регистрация
- `POST /api/login` - Вход
- `GET /swagger/*` - Swagger UI

### Защищенные (требуют JWT токен):
- `POST /api/address/search` - Поиск адресов
- `POST /api/address/geocode` - Геокодирование

## ⚠️ Коды ошибок

### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```

### 401 Unauthorized
Возвращается когда:
- Токен отсутствует
- Токен невалидный
- Токен истек
- Неверные учетные данные при входе

### 500 Internal Server Error
```json
{
  "error": "Internal server error message"
}
```

## 🔄 Полный workflow

```powershell
# 1. Регистрация
$registerBody = @{
    username = "newuser@example.com"
    password = "securepassword"
} | ConvertTo-Json

$registerResponse = Invoke-RestMethod `
    -Uri "http://localhost:8080/api/register" `
    -Method Post `
    -ContentType "application/json" `
    -Body $registerBody

$token = $registerResponse.token
Write-Host "✓ Зарегистрирован. Token получен."

# 2. Поиск адреса
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

$searchBody = @{
    query = "Санкт-Петербург, Невский проспект, 1"
} | ConvertTo-Json

$searchResult = Invoke-RestMethod `
    -Uri "http://localhost:8080/api/address/search" `
    -Method Post `
    -Headers $headers `
    -Body $searchBody

Write-Host "✓ Найдено адресов:" $searchResult.addresses.Count
$searchResult.addresses | Format-Table City, Street, House, Lat, Lon

# 3. Геокодирование
$geocodeBody = @{
    lat = $searchResult.addresses[0].lat
    lng = $searchResult.addresses[0].lon
} | ConvertTo-Json

$geocodeResult = Invoke-RestMethod `
    -Uri "http://localhost:8080/api/address/geocode" `
    -Method Post `
    -Headers $headers `
    -Body $geocodeBody

Write-Host "✓ Обратное геокодирование выполнено"
$geocodeResult.addresses | Format-Table City, Street, House
```

## 🛠️ Настройка

### Изменение секретного ключа

В production окружении **обязательно** измените секретный ключ:

```go
var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("your-secret-key"), nil)
```

Рекомендуется использовать переменные окружения:

```go
secretKey := os.Getenv("JWT_SECRET_KEY")
if secretKey == "" {
    secretKey = "default-secret-key" // только для разработки
}
var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(secretKey), nil)
```

### Добавление времени истечения токена

```go
claims := map[string]interface{}{
    "email": data.Username,
    "exp": time.Now().Add(time.Hour * 24).Unix(), // истекает через 24 часа
}
_, tokenString, err := tokenAuth.Encode(claims)
```

## 🧪 Тестирование без регистрации

Для быстрого тестирования можно использовать заранее созданный токен (если он не истек):

```powershell
$token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

$headers = @{
    "Authorization" = "Bearer $token"
}

# Используйте headers в запросах
```

## 📚 Дополнительные ресурсы

- [JWT.io](https://jwt.io/) - Декодирование и верификация токенов
- [go-chi/jwtauth](https://github.com/go-chi/jwtauth) - Документация библиотеки
- [RFC 7519](https://tools.ietf.org/html/rfc7519) - JWT спецификация

## 💡 Best Practices

1. **Храните токены безопасно** - не коммитьте токены в git
2. **Используйте HTTPS** в production
3. **Устанавливайте время истечения** для токенов
4. **Регулярно меняйте секретный ключ**
5. **Валидируйте все входные данные**
6. **Логируйте попытки аутентификации**
7. **Используйте сильные пароли**

## 🔍 Отладка

### Декодирование токена

Используйте [jwt.io](https://jwt.io/) для просмотра содержимого токена:

1. Скопируйте токен
2. Вставьте на jwt.io
3. Посмотрите payload и claims

### Проверка токена в коде

```go
token, err := tokenAuth.Decode(tokenString)
if err != nil {
    log.Printf("Invalid token: %v", err)
}
```

## ❓ FAQ

**Q: Как долго действует токен?**  
A: По умолчанию токен не имеет срока истечения. Добавьте claim `exp` для ограничения.

**Q: Можно ли отозвать токен?**  
A: Текущая реализация не поддерживает отзыв. Для этого нужна база данных с черным списком токенов.

**Q: Безопасно ли хранить пароли?**  
A: Да, используется bcrypt для хеширования паролей с cost factor по умолчанию (10).

**Q: Что делать при ошибке 401?**  
A: Проверьте, что токен правильно передан в заголовке `Authorization: Bearer <token>`.
