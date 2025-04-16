# Практическое задание 6: Аутентификация и авторизация в Go

Код проекта доступен [по ссылке](https://github.com/LeetManSup/mirea_finance_tracker/).

## 1. Базовая модель пользователя
```go
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	FullName     string
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	Accounts              []Account              `gorm:"foreignKey:UserID"`
	Categories            []Category             `gorm:"foreignKey:UserID"`
	Transactions          []Transaction          `gorm:"foreignKey:UserID"`
	RecurringTransactions []RecurringTransaction `gorm:"foreignKey:UserID"`
}
```


## 2. Реализация хеширования пароля
Пример метода Login структуры AuthService из `internal/service/auth_service.go`:
```go
func (s *AuthService) Login(email, password string) (string, error) {
    
    ...

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return generateJWT(user.ID)
}
```

## 3. Генерация и проверка JWT
`internal/service/auth_service.go`
Инициализация ключа JWT
```go
var jwtKey = []byte(config.Load().JWTSecret)
```

Генерация токена:
```go
func generateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
```

Мидлвар проверки токена из `internal/middleware/auth.go`
```go
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		fmt.Println("claims:", claims)

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
```

## 4. Роутинг
Защищённые маршруты (`internal/router/router.go`)
```go
func SetupRouter(db *gorm.DB) *gin.Engine {

    ...

	// Приватные маршруты с JWT middleware
	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	auth.GET("/me", userHandler.GetMe)
	auth.GET("/accounts", accountHandler.GetAccounts)
	auth.GET("/accounts/:id", accountHandler.GetAccount)
	auth.POST("/accounts", accountHandler.CreateAccount)
	auth.PATCH("/accounts/:id", accountHandler.UpdateAccount)
	auth.DELETE("/accounts/:id", accountHandler.DeleteAccount)

	return r
}
```