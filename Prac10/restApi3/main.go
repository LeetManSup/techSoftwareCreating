package main

import (
	"restApi3/db"
	"restApi3/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()

	router := gin.Default()

	// Продукты
	router.GET("/products", handlers.GetProducts)
	router.POST("/products", handlers.CreateProduct)
	router.PUT("/products/:id", handlers.UpdateProduct)
	router.DELETE("/products/:id", handlers.DeleteProduct)

	// Клиенты
	router.GET("/customers", handlers.GetCustomers)
	router.POST("/customers", handlers.CreateCustomer)
	router.PUT("/customers/:id", handlers.UpdateCustomer)
	router.DELETE("/customers/:id", handlers.DeleteCustomer)

	// Корзины
	router.GET("/carts", handlers.GetCarts)
	router.POST("/carts", handlers.CreateCart)
	router.PUT("/carts/:id", handlers.UpdateCart)
	router.DELETE("/carts/:id", handlers.DeleteCart)

	// Заказы
	router.GET("/orders", handlers.GetOrders)
	router.POST("/orders", handlers.CreateOrder)
	router.PUT("/orders/:id", handlers.UpdateOrder)
	router.DELETE("/orders/:id", handlers.DeleteOrder)

	router.Run(":8080")
}
