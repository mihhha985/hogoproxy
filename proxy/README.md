# HugoProxy Address API 🚀

API для поиска адресов и геокодирования с JWT аутентификацией.

## ✨ Возможности

- 🔍 Поиск адресов по текстовому запросу
- 📍 Обратное геокодирование (координаты → адрес)
- 🔐 JWT аутентификация для защиты endpoints
- 📚 Полная Swagger документация
- 🐳 Docker поддержка
- ✅ Unit и интеграционные тесты

## 🏗️ Архитектура

- **Framework:** Chi Router
- **Authentication:** JWT (JSON Web Tokens)
- **Password Hashing:** bcrypt
- **API Provider:** DaData
- **Documentation:** Swagger/OpenAPI
- **Language:** Go 1.24

## 🚀 Быстрый старт

### Установка зависимостей

```powershell
cd c:\projects\hugoproxy\proxy
go mod download
```

### Генерация Swagger документации

```powershell
swag init -g main.go --output ./docs
```

### Запуск сервера

```powershell
go run .
```

Сервер запустится на `http://localhost:8080`

## 📚 Документация API

После запуска сервера, откройте Swagger UI:

**http://localhost:8080/swagger/**

## 🔐 Аутентификация

API использует JWT токены для аутентификации. Полное руководство: [JWT_AUTH_GUIDE.md](JWT_AUTH_GUIDE.md)

### Быстрый пример

```powershell
# 1. Регистрация
$body = @{
    username = "user@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/register" `
  -Method Post -ContentType "application/json" -Body $body

$token = $response.token

# 2. Использование токена
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

$searchBody = @{ query = "Москва, Красная площадь, 1" } | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/address/search" `
  -Method Post -Headers $headers -Body $searchBody
```

## 🔌 API Endpoints

### Публичные (не требуют аутентификации)

#### POST `/api/register`
Регистрация нового пользователя

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
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### POST `/api/login`
Вход в систему

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
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Защищенные (требуют JWT токен)

#### POST `/api/address/search`
Поиск адресов по текстовому запросу

**Заголовки:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Запрос:**
```json
{
  "query": "Санкт-Петербург, Невский проспект, 1"
}
```

**Ответ:**
```json
{
  "addresses": [
    {
      "city": "Санкт-Петербург",
      "street": "Невский",
      "house": "1",
      "lat": "59.934280",
      "lon": "30.335099"
    }
  ]
}
```

#### POST `/api/address/geocode`
Обратное геокодирование

**Заголовки:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Запрос:**
```json
{
  "lat": "59.934280",
  "lng": "30.335099"
}
```

**Ответ:**
```json
{
  "addresses": [
    {
      "city": "Санкт-Петербург",
      "street": "Невский",
      "house": "1",
      "lat": "59.934280",
      "lon": "30.335099"
    }
  ]
}
```

## 🧪 Тестирование

### Запуск тестов

```powershell
# Все тесты
go test -v

# С покрытием
go test -v -cover

# Конкретный тест
go test -v -run TestHandlerAddressSearch
```

### Примеры тестирования API

Смотрите [API_EXAMPLES.md](API_EXAMPLES.md) для подробных примеров.

## 🐳 Docker

### Сборка образа

```powershell
docker build -t hugoproxy-api .
```

### Запуск контейнера

```powershell
docker run -d -p 8080:8080 --name proxy-service hugoproxy-api
```

### Docker Compose

```yaml
version: '3.8'
services:
  proxy:
    build: ./proxy
    ports:
      - "8080:8080"
    environment:
      - JWT_SECRET_KEY=your-secret-key
```

## 📁 Структура проекта

```
proxy/
├── main.go              # Основной файл с роутами и обработчиками
├── address.go           # Логика работы с GeoService
├── geocodejson.go       # Структуры для DaData API
├── main_test.go         # Тесты
├── go.mod               # Зависимости
├── go.sum               # Контрольные суммы
├── Dockerfile           # Docker конфигурация
├── .dockerignore        # Исключения для Docker
├── docs/                # Сгенерированная Swagger документация
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── README.md            # Этот файл
├── JWT_AUTH_GUIDE.md    # Руководство по JWT
├── API_EXAMPLES.md      # Примеры запросов
└── QUICKSTART.md        # Быстрый старт
```

## 🔧 Конфигурация

### API ключи DaData

Настраиваются в `main.go`:

```go
const DADATA_API_KEY = "your-api-key"
const DADATA_SECRET_KEY = "your-secret-key"
```

### JWT секретный ключ

```go
var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("your-secret"), nil)
```

⚠️ **Важно:** В production используйте переменные окружения!

```go
secretKey := os.Getenv("JWT_SECRET_KEY")
if secretKey == "" {
    log.Fatal("JWT_SECRET_KEY environment variable is required")
}
```

## 📦 Зависимости

```go
require (
    github.com/ekomobile/dadata/v2 v2.10.0
    github.com/go-chi/chi/v5 v5.2.3
    github.com/go-chi/jwtauth/v5 v5.x.x
    github.com/swaggo/http-swagger/v2 v2.0.2
    github.com/swaggo/swag v1.16.6
    golang.org/x/crypto // for bcrypt
)
```

### Установка дополнительных инструментов

```powershell
# Swagger CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Testify (для тестов)
go get github.com/stretchr/testify
```

## 🔄 Workflow разработки

1. **Внесите изменения** в код
2. **Добавьте/обновите Swagger аннотации**
3. **Регенерируйте документацию:**
   ```powershell
   swag init -g main.go --output ./docs
   ```
4. **Запустите тесты:**
   ```powershell
   go test -v
   ```
5. **Запустите сервер и проверьте в Swagger UI:**
   ```powershell
   go run .
   ```

## 🛡️ Безопасность

### Best Practices

1. ✅ Используйте HTTPS в production
2. ✅ Храните секретные ключи в переменных окружения
3. ✅ Устанавливайте время истечения для JWT токенов
4. ✅ Используйте сильные пароли
5. ✅ Валидируйте все входные данные
6. ✅ Логируйте попытки аутентификации
7. ✅ Регулярно обновляйте зависимости

### Хеширование паролей

Используется bcrypt с default cost factor (10):

```go
hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

## 📊 Производительность

### Бенчмарки

```powershell
go test -bench=. -benchmem
```

### Профилирование

```powershell
go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=.
go tool pprof cpu.prof
```

## 🐛 Отладка

### Включение подробного логирования

```go
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
```

### Проверка JWT токена

Используйте [jwt.io](https://jwt.io/) для декодирования токенов.

### Логирование claims

```go
_, claims, _ := jwtauth.FromContext(r.Context())
log.Printf("User claims: %+v", claims)
```

## 📝 Changelog

### Version 1.0.0
- ✨ Добавлена JWT аутентификация
- ✨ Регистрация и вход пользователей
- ✨ Защита endpoints с помощью JWT
- ✨ Полная Swagger документация
- ✨ Unit и интеграционные тесты
- ✨ Docker поддержка

## 🤝 Contributing

1. Fork репозиторий
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📄 Лицензия

MIT License

## 👥 Контакты

- GitHub: [@mihhha985](https://github.com/mihhha985)
- Project: [hogoproxy](https://github.com/mihhha985/hogoproxy)

## 🔗 Полезные ссылки

- [Swagger Documentation](http://localhost:8080/swagger/)
- [JWT Authentication Guide](JWT_AUTH_GUIDE.md)
- [API Examples](API_EXAMPLES.md)
- [Quick Start Guide](QUICKSTART.md)
- [DaData API](https://dadata.ru/)
- [Chi Router](https://github.com/go-chi/chi)
- [JWT Auth](https://github.com/go-chi/jwtauth)

---

Made with ❤️ using Go
