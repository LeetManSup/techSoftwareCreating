package handlers

import (
	"net/http"
	"restApi3/db"
	"restApi3/models"

	"github.com/gin-gonic/gin"
)

func GetCustomers(c *gin.Context) {
	var customers []models.Customer
	db.DB.Find(&customers)
	c.JSON(http.StatusOK, customers)
}

func CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&customer)
	c.JSON(http.StatusCreated, customer)
}

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

func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Customer{}, "CustomerID = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted"})
}
