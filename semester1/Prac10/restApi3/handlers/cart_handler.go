package handlers

import (
	"net/http"
	"restApi3/db"
	"restApi3/models"

	"github.com/gin-gonic/gin"
)

func GetCarts(c *gin.Context) {
	var carts []models.Cart
	db.DB.Find(&carts)
	c.JSON(http.StatusOK, carts)
}

func CreateCart(c *gin.Context) {
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&cart)
	c.JSON(http.StatusCreated, cart)
}

func UpdateCart(c *gin.Context) {
	id := c.Param("id")
	var cart models.Cart
	if err := db.DB.First(&cart, "CartID = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Save(&cart)
	c.JSON(http.StatusOK, cart)
}

func DeleteCart(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Cart{}, "CartID = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cart deleted"})
}
