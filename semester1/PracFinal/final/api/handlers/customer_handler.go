package handlers

import (
    "net/http"
    "final/db"
    "final/models"

    "github.com/gin-gonic/gin"
)

// @Summary Get all customers
// @Description Get a list of all customers
// @Tags Customers
// @Produce json
// @Success 200 {array} models.Customer
// @Router /customers [get]
func GetCustomers(c *gin.Context) {
    var customers []models.Customer
    db.DB.Find(&customers)
    c.JSON(http.StatusOK, customers)
}

// @Summary Create new customer
// @Description Creates a new customer
// @Tags Customers
// @Accept json
// @Produce json
// @Param customer body models.Customer true "Customer object to be created"
// @Success 201 {object} models.Customer
// @Failure 400 {object} string "Invalid request body"
// @Router /customers [post]
func CreateCustomer(c *gin.Context) {
    var customer models.Customer
    if err := c.ShouldBindJSON(&customer); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.DB.Create(&customer)
    c.JSON(http.StatusCreated, customer)
}

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
func UpdateCustomer(c *gin.Context) {
    id := c.Param("id")
    var customer models.Customer
    if err := db.DB.First(&customer, "CustomerID = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
        return
    }
    if err := c.ShouldBindJSON(&customer); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.DB.Save(&customer)
    c.JSON(http.StatusOK, customer)
}

// @Summary Delete a customer
// @Description Deletes a customer by ID
// @Tags Customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {string} string "Customer deleted"
// @Failure 404 {string} string "Customer not found"
// @Router /customers/{id} [delete]
func DeleteCustomer(c *gin.Context) {
    id := c.Param("id")
    if err := db.DB.Delete(&models.Customer{}, "CustomerID = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Customer deleted"})
}
