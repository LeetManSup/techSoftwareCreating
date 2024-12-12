package handlers

import (
    "net/http"
    "final/db"
    "final/models"

    "github.com/gin-gonic/gin"
)

// @Summary Get all orders
// @Description Get a list of all customer orders
// @Tags Orders
// @Produce json
// @Success 200 {array} models.CustomerOrder
// @Router /orders [get]
func GetOrders(c *gin.Context) {
    var orders []models.CustomerOrder
    db.DB.Find(&orders)
    c.JSON(http.StatusOK, orders)
}

// @Summary Create new order
// @Description Creates a new customer order
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.CustomerOrder true "Order object to be created"
// @Success 201 {object} models.CustomerOrder
// @Failure 400 {object} string "Invalid request body"
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
    var order models.CustomerOrder
    if err := c.ShouldBindJSON(&order); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.DB.Create(&order)
    c.JSON(http.StatusCreated, order)
}

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
func UpdateOrder(c *gin.Context) {
    id := c.Param("id")
    var order models.CustomerOrder
    if err := db.DB.First(&order, "OrderID = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        return
    }
    if err := c.ShouldBindJSON(&order); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.DB.Save(&order)
    c.JSON(http.StatusOK, order)
}

// @Summary Delete an order
// @Description Deletes a customer order by ID
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {string} string "Order deleted"
// @Failure 404 {string} string "Order not found"
// @Router /orders/{id} [delete]
func DeleteOrder(c *gin.Context) {
    id := c.Param("id")
    if err := db.DB.Delete(&models.CustomerOrder{}, "OrderID = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}
