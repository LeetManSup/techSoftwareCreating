package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var jwtKey = []byte("secret")

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Good struct {
	ID          string  `json:"id" example:"1"`
	Name        string  `json:"name" example:"Стол"`
	Description string  `json:"description" example:"Обычный деревянный стол"`
	Price       float32 `json:"price" example:"10000.0"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"id not found"`
}

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

var goods = []Good{
	{
		ID:          "1",
		Name:        "Стол",
		Description: "Обычный деревянный стол",
		Price:       10000.0,
	},
	{
		ID:          "2",
		Name:        "Стул",
		Description: "Обычный железный стул",
		Price:       3000.0,
	},
	{
		ID:          "3",
		Name:        "Ковёр",
		Description: "Красный совковый ковёр",
		Price:       5000.0,
	},
}

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

func getGoods(c *gin.Context) {
	c.JSON(http.StatusOK, goods)
}

func getGood(c *gin.Context) {
	id := c.Param("id")

	for _, good := range goods {
		if good.ID == id {
			c.JSON(http.StatusOK, good)
			return
		}
	}
	c.JSON(http.StatusNotFound, ErrorResponse{Error: "id not found"})
}

func createGood(c *gin.Context) {
	var newGood Good

	if err := c.ShouldBind(&newGood); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	goods = append(goods, newGood)
	c.JSON(http.StatusCreated, newGood)
}

func updateGood(c *gin.Context) {
	id := c.Param("id")
	var updatedGood Good

	if err := c.ShouldBind(&updatedGood); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	for i, good := range goods {
		if good.ID == id {
			goods[i] = updatedGood
			c.JSON(http.StatusOK, updatedGood)
			return
		}
	}
	c.JSON(http.StatusNotFound, ErrorResponse{Error: "id not found"})
}

func deleteGood(c *gin.Context) {
	id := c.Param("id")

	for i, good := range goods {
		if good.ID == id {
			goods = append(goods[:i], goods[i+1:]...)
			c.JSON(http.StatusOK, good)
			return
		}
	}
	c.JSON(http.StatusNotFound, ErrorResponse{Error: "id not found"})
}

func main() {
	router := gin.Default()

	router.POST("/register", register) // Регистрация
	router.POST("/login", login)       // Вход
	router.POST("/refresh", refresh)   // Маршрут для обновления Access Token

	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/goods", getGoods)                                  // Получение всех товаров
		protected.GET("/good/:id", getGood)                                // Получение товара по ID
		protected.POST("/goods", roleMiddleware("admin"), createGood)      // Создание нового товара
		protected.PATCH("/good/:id", roleMiddleware("admin"), updateGood)  // Обновление существующего товара
		protected.DELETE("/good/:id", roleMiddleware("admin"), deleteGood) // Удаление товара

	}

	_ = router.Run(":8888")
}
