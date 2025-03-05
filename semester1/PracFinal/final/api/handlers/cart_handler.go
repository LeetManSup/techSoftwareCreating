package handlers

import (
    "net/http"
    "final/db"
    "final/models"

    "github.com/gin-gonic/gin"
)

// @Summary Get all carts
// @Description Get a list of all shopping carts
// @Tags Carts
// @Produce json
// @Success 200 {array} models.Cart
// @Router /carts [get]
func GetCarts(c *gin.Context) {
    var carts []models.Cart
    db.DB.Find(&carts)
    c.JSON(http.StatusOK, carts)
}

// @Summary Create new cart
// @Description Creates a new shopping cart
// @Tags Carts
// @Accept json
// @Produce json
// @Param cart body models.Cart true "Cart object to be created"
// @Success 201 {object} models.Cart
// @Failure 400 {object} string "Invalid request body"
// @Router /carts [post]
func CreateCart(c *gin.Context) {
    var cart models.Cart
    if err := c.ShouldBindJSON(&cart); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.DB.Create(&cart)
    c.JSON(http.StatusCreated, cart)
}

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

// @Summary Delete a cart
// @Description Deletes a shopping cart by ID
// @Tags Carts
// @Produce json
// @Param id path int true "Cart ID"
// @Success 200 {string} string "Cart deleted"
// @Failure 404 {string} string "Cart not found"
// @Router /carts/{id} [delete]
func DeleteCart(c *gin.Context) {
    id := c.Param("id")
    if err := db.DB.Delete(&models.Cart{}, "CartID = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Cart deleted"})
}
