package main

import (
    "final/db"
    "final/handlers"

    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    //"github.com/swaggo/http-swagger"
    //"log"
    //"net/http"
    _ "final/docs"
)

// @title Final Project API
// @version 1.0
// @description This is a sample server for final project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
    db.ConnectDB()

    router := gin.Default()

    // Routes for Products
    // @Summary List all products
    // @Description Get a list of all products
    // @Tags Products
    // @Produce json
    // @Success 200 {array} models.Product
    // @Router /products [get]
    router.GET("/products", handlers.GetProducts)

    // @Summary Create a new product
    // @Description Creates a new product
    // @Tags Products
    // @Accept json
    // @Produce json
    // @Param product body models.Product true "Product object to be created"
    // @Success 201 {object} models.Product
    // @Failure 400 {object} string "Invalid request body"
    // @Router /products [post]
    router.POST("/products", handlers.CreateProduct)

    // @Summary Update an existing product
    // @Description Updates the product by ID
    // @Tags Products
    // @Accept json
    // @Produce json
    // @Param id path int true "Product ID"
    // @Param product body models.Product true "Updated product object"
    // @Success 200 {object} models.Product
    // @Failure 404 {string} string "Product not found"
    // @Failure 400 {string} string "Invalid request body"
    // @Router /products/{id} [put]
    router.PUT("/products/:id", handlers.UpdateProduct)

    // @Summary Delete a product
    // @Description Deletes a product by ID
    // @Tags Products
    // @Produce json
    // @Param id path int true "Product ID"
    // @Success 200 {string} string "Product deleted"
    // @Failure 404 {string} string "Product not found"
    // @Router /products/{id} [delete]
    router.DELETE("/products/:id", handlers.DeleteProduct)

    // Routes for Customers
    // @Summary List all customers
    // @Description Get a list of all customers
    // @Tags Customers
    // @Produce json
    // @Success 200 {array} models.Customer
    // @Router /customers [get]
    router.GET("/customers", handlers.GetCustomers)

    // @Summary Create a new customer
    // @Description Creates a new customer
    // @Tags Customers
    // @Accept json
    // @Produce json
    // @Param customer body models.Customer true "Customer object to be created"
    // @Success 201 {object} models.Customer
    // @Failure 400 {object} string "Invalid request body"
    // @Router /customers [post]
    router.POST("/customers", handlers.CreateCustomer)

    // @Summary Update an existing customer
    // @Description Updates the customer by ID
    // @Tags Customers
    // @Accept json
    // @Produce json
    // @Param id path int true "Customer ID"
    // @Param customer body models.Customer true "Updated customer object"
    // @Success 200 {object} models.Customer
    // @Failure 404 {string} string "Customer not found"
    // @Failure 400 {string} string "Invalid request body"
    // @Router /customers/{id} [put]
    router.PUT("/customers/:id", handlers.UpdateCustomer)

    // @Summary Delete a customer
    // @Description Deletes a customer by ID
    // @Tags Customers
    // @Produce json
    // @Param id path int true "Customer ID"
    // @Success 200 {string} string "Customer deleted"
    // @Failure 404 {string} string "Customer not found"
    // @Router /customers/{id} [delete]
    router.DELETE("/customers/:id", handlers.DeleteCustomer)

    // Routes for Carts
    // @Summary List all carts
    // @Description Get a list of all shopping carts
    // @Tags Carts
    // @Produce json
    // @Success 200 {array} models.Cart
    // @Router /carts [get]
    router.GET("/carts", handlers.GetCarts)

    // @Summary Create a new cart
    // @Description Creates a new shopping cart
    // @Tags Carts
    // @Accept json
    // @Produce json
    // @Param cart body models.Cart true "Cart object to be created"
    // @Success 201 {object} models.Cart
    // @Failure 400 {object} string "Invalid request body"
    // @Router /carts [post]
    router.POST("/carts", handlers.CreateCart)

    // @Summary Update an existing cart
    // @Description Updates the shopping cart by ID
    // @Tags Carts
    // @Accept json
    // @Produce json
    // @Param id path int true "Cart ID"
    // @Param cart body models.Cart true "Updated cart object"
    // @Success 200 {object} models.Cart
    // @Failure 404 {string} string "Cart not found"
    // @Failure 400 {string} string "Invalid request body"
    // @Router /carts/{id} [put]
    router.PUT("/carts/:id", handlers.UpdateCart)

    // @Summary Delete a cart
    // @Description Deletes a shopping cart by ID
    // @Tags Carts
    // @Produce json
    // @Param id path int true "Cart ID"
    // @Success 200 {string} string "Cart deleted"
    // @Failure 404 {string} string "Cart not found"
    // @Router /carts/{id} [delete]
    router.DELETE("/carts/:id", handlers.DeleteCart)

    // Routes for Orders
    // @Summary List all orders
    // @Description Get a list of all customer orders
    // @Tags Orders
    // @Produce json
    // @Success 200 {array} models.CustomerOrder
    // @Router /orders [get]
    router.GET("/orders", handlers.GetOrders)

    // @Summary Create a new order
    // @Description Creates a new customer order
    // @Tags Orders
    // @Accept json
    // @Produce json
    // @Param order body models.CustomerOrder true "Order object to be created"
    // @Success 201 {object} models.CustomerOrder
    // @Failure 400 {object} string "Invalid request body"
    // @Router /orders [post]
    router.POST("/orders", handlers.CreateOrder)

    // @Summary Update an existing order
    // @Description Updates the customer order by ID
    // @Tags Orders
    // @Accept json
    // @Produce json
    // @Param id path int true "Order ID"
    // @Param order body models.CustomerOrder true "Updated order object"
    // @Success 200 {object} models.CustomerOrder
    // @Failure 404 {string} string "Order not found"
    // @Failure 400 {string} string "Invalid request body"
    // @Router /orders/{id} [put]
    router.PUT("/orders/:id", handlers.UpdateOrder)

    // @Summary Delete an order
    // @Description Deletes a customer order by ID
    // @Tags Orders
    // @Produce json
    // @Param id path int true "Order ID"
    // @Success 200 {string} string "Order deleted"
    // @Failure 404 {string} string "Order not found"
    // @Router /orders/{id} [delete]
    router.DELETE("/orders/:id", handlers.DeleteOrder)

    router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    router.Run(":8888")
}

