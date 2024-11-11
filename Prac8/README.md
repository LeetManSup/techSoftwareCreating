# Практическая работа №8. Аутентификация и авторизация в REST API

## 0. Установка зависимостей (JWT-Go)
Установим библиотеку для работы с JWT:
```bash
    go get -u github.com/dgrijalva/jwt-go
```

## 1. Подготовка переменных и структур
Создадим 3 структуры:
- Пользователь
    ```go
    type User struct {
        Username string `json:"username"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }
    ```
- Требования для Access JWT
    ```go
    type Claims struct {
        Username string `json:"username"`
        Role     string `json:"role"`
        jwt.StandardClaims
    }
    ```
- Требования для Refresh JWT
    ```go
    type RefreshClaims struct {
        Username string `json:"username"`
        jwt.StandardClaims
    }
    ```

Создадим переменную для хранения "секрета" генерации токенов:
```bash
var jwtKey = []byte("secret")
```

Сымитируем базу данных, создав срез тестовых пользователей:
```bash
var users = []User{
	{
		Username: "admin",
		Password: "adminPassword",
		Role:     "admin",
	},
	{
		Username: "user1",
		Password: "password1",
		Role:     "user",
	},
	{
		Username: "user2",
		Password: "password2",
		Role:     "user",
	},
}
```

## 2. Реализация функций
Напишем функцию, которая генерирует Access и Refresh токены:
```go
func generateTokens(username, role string) (accessToken, refreshToken string, err error) {
	// Генерация Access Token (короткоживущего токена)
	accessTokenExpiration := time.Now().Add(5 * time.Minute)
	accessClaims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpiration.Unix(),
		},
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Генерация Refresh Token (долгоживущего токена)
	refreshTokenExpiration := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &RefreshClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiration.Unix(),
		},
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString(jwtKey)
	return
}
```

Напишем функцию, реализующую логику регистрации в системе:
```go
func register(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	for _, user := range users {
		if user.Username == newUser.Username {
			c.JSON(http.StatusConflict, gin.H{"message": "user already exists"})
			return
		}
	}

	users = append(users, newUser)
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}
```

Напишем функцию, реализующую логику входа в систему: 
```go
func login(c *gin.Context) {
	var credentials User
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	var role string
	for _, user := range users {
		if user.Username == credentials.Username && user.Password == credentials.Password {
			role = user.Role
			break
		}
	}

	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	accessToken, refreshToken, err := generateTokens(credentials.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
```

Напишем функцию, реализующую логику обновления Access токена:
```go
func refresh(c *gin.Context) {
	var requestBody struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	refreshToken := requestBody.RefreshToken
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	// Генерация нового Access Token
	var role string
	for _, user := range users {
		if user.Username == claims.Username {
			role = user.Role
			break
		}
	}

	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	accessTokenExpiration := time.Now().Add(5 * time.Minute)
	accessClaims := &Claims{
		Username: claims.Username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpiration.Unix(),
		},
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	newAccessToken, err := accessTokenObj.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
```

Напишем промежуточную функцию (мидлвар), реализующую логику авторизации по JWT:
```go
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
```

Напишем промежуточную функцию (мидлвар), реализующую логику проверки соответствия прав пользователя на совершение запроса:
```go
func roleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			c.Abort()
			return
		}

		userClaims := claims.(*Claims)
		if userClaims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
```

Обновим функцию main:
```go
func main() {
	router := gin.Default()

	router.POST("/register", register) // Регистрация
	router.POST("/login", login)       // Вход
	router.POST("/refresh", refresh)   // Обновление Access Token

	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/goods", getGoods)
		protected.GET("/good/:id", getGood)
		protected.POST("/goods", roleMiddleware("admin"), createGood)
		protected.PATCH("/good/:id", roleMiddleware("admin"), updateGood)
		protected.DELETE("/good/:id", roleMiddleware("admin"), deleteGood)

	}

	_ = router.Run(":8888")
}
```

## 3. Результаты
### Регистрация пользователя в системе
Отправляем `POST` запрос на `http://localhost:8888/register` c `body` со следующим содержимым (логин, пароль и роль регистрируемого пользователя):
```json
{
    "username": "newuser",
    "password": "newpassword",
    "role": "user"
}
```

Получаем ответ (сообщение об успешной регистрации):
```json
{
    "message": "user registered successfully"
}
```
Скриншот:
![alt text](pictures/image1.png)

### Вход администратора в систему
Отправляем `POST` запрос на `http://localhost:8888/login` c `body` со следующим содержимым (логин и пароль администратора):
```json
{
  "username": "admin",
  "password": "adminPassword"
}
```

Получаем ответ (access и refresh токены пользователя `admin`):
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzMxMzE4NDc5fQ.X_3mI4DGda4nonKiB5jEI0__jP-Q0CD-B6PUwh76avo",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzMxOTIyOTc5fQ.W6yaA80-SBL-LekYF-wuoIRMUVKtnwU9sA42AGh1lvs"
}
```

Скриншот:
![alt text](pictures/image2.png)

### Вход рядового пользователя в систему
Отправляем `POST` запрос на `http://localhost:8888/login` c `body` со следующим содержимым (логин и пароль рядового пользователя):
```json
{
  "username": "user1",
  "password": "password1"
}
```

Получаем ответ (access и refresh токены пользователя `user1`):
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxIiwicm9sZSI6InVzZXIiLCJleHAiOjE3MzEzMTg1OTh9.Sno9GU6NGWSIvjr1UWTwOeH8siJU6Y_JQw8daVpuCgI",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxIiwiZXhwIjoxNzMxOTIzMDk4fQ.nw_BpTS9cjV6kKiQIb1ByTXRIQm7J7LVlweYnq6neP8"
}
```

Скриншот:
![alt text](pictures/image3.png)

### Обновление Access Token (на примере токена рядового пользователя `user1`)
Отправляем `POST` запрос на `http://localhost:8888/refresh` c `body` со следующим содержимым (refresh token пользователя `user1`):
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxIiwiZXhwIjoxNzMxOTIzMDk4fQ.nw_BpTS9cjV6kKiQIb1ByTXRIQm7J7LVlweYnq6neP8"
}
```

Получаем ответ (новый access token пользователя `user1`):
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxIiwicm9sZSI6InVzZXIiLCJleHAiOjE3MzEzMTg4MDh9.Ra-EzAVP5lpkFpSP3nXMzJJgaHrwrEbnidTacjSUACc"
}
```

Скриншот:
![alt text](pictures/image4.png)

### Проверка доступности маршрутов для рядового пользователя (на примере получения списка товаров и добавления нового товара)
#### Получение списка товаров
Данный маршрут должен быть доступен рядовому пользователю.

Отправляем `GET` запрос на `http://localhost:8888/goods` cо следующим `header` (JWT рядового пользователя `user1`):
```plaintext
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxIiwicm9sZSI6InVzZXIiLCJleHAiOjE3MzEzMTg4MDh9.Ra-EzAVP5lpkFpSP3nXMzJJgaHrwrEbnidTacjSUACc
```

Получаем ответ (маршрут доступен, получили список товаров):
```json
[
    {
        "id": "1",
        "name": "Стол",
        "description": "Обычный деревянный стол",
        "price": 10000
    },
    {
        "id": "2",
        "name": "Стул",
        "description": "Обычный железный стул",
        "price": 3000
    },
    {
        "id": "3",
        "name": "Ковёр",
        "description": "Красный совковый ковёр",
        "price": 5000
    }
]
```

Скриншот:
![alt text](pictures/image5.png)

#### Добавление нового товара
Данный маршрут должен быть запрещён рядовому пользователю.

Отправляем `POST` запрос на `http://localhost:8888/goods` cо следующим `header` (JWT рядового пользователя `user1`):
```plaintext
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxIiwicm9sZSI6InVzZXIiLCJleHAiOjE3MzEzMTg4MDh9.Ra-EzAVP5lpkFpSP3nXMzJJgaHrwrEbnidTacjSUACc
```

и со следующим `body` (данные добавляемого товара):
```json
{
    "id": "4",
    "name": "ПОЛ",
    "description": "ПАРКЕТ",
    "price": 10
}
```

Получаем ответ (доступ запрещён):
```json
{
    "message": "forbidden"
}
```

Скриншот:
![alt text](pictures/image6.png)

### Проверка доступности маршрутов для администратора (на примере получения списка товаров и добавления нового товара)
#### Получение списка товаров
Данный маршрут должен быть доступен рядовому пользователю.

Отправляем `GET` запрос на `http://localhost:8888/goods` cо следующим `header` (JWT администратора `admin`):
```plaintext
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzMxMzE5Nzk5fQ.QVO5hf1RXr1l35rzBFeSZNq66CghmX4H_YO2cHKH6BQ
```

Получаем ответ (маршрут доступен, получили список товаров):
```json
[
    {
        "id": "1",
        "name": "Стол",
        "description": "Обычный деревянный стол",
        "price": 10000
    },
    {
        "id": "2",
        "name": "Стул",
        "description": "Обычный железный стул",
        "price": 3000
    },
    {
        "id": "3",
        "name": "Ковёр",
        "description": "Красный совковый ковёр",
        "price": 5000
    }
]
```

Скриншот:
![alt text](pictures/image7.png)

#### Добавление нового товара
Данный маршрут должен быть доступен администратору.

Отправляем `POST` запрос на `http://localhost:8888/goods` cо следующим `header` (JWT администратора `admin`):
```plaintext
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzMxMzE5Nzk5fQ.QVO5hf1RXr1l35rzBFeSZNq66CghmX4H_YO2cHKH6BQ
```

и со следующим `body` (данные добавляемого товара):
```json
{
    "id": "5",
    "name": "НОВЫЙ_ПОЛ",
    "description": "СУПЕР_ПАРКЕТ",
    "price": 20
}
```

Получаем ответ (маршрут доступен, товар был добавлен):
```json
{
    "id": "5",
    "name": "НОВЫЙ_ПОЛ",
    "description": "СУПЕР_ПАРКЕТ",
    "price": 20
}
```

Скриншот:
![alt text](pictures/image8.png)

Проверка добавления товара (товар успешно добавлен в систему):
![alt text](pictures/image9.png)