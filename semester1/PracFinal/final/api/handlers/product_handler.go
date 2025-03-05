package handlers

import (
    "net/http"
    "final/db"
    "final/models"

    "github.com/gin-gonic/gin"
)

// @Summary Get all products
// @Description Get a list of all products
// @Tags Products
// @Produce json
// @Success 200 {array} models.Product
// @Router /products [get]
func GetProducts(c *gin.Context) {
    var products []models.Product
    db.DB.Find(&products)
    c.JSON(http.StatusOK, products)
}

// @Summary Create new product
// @Description Creates a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product object to be created"
// @Success 201 {object} models.Product
// @Failure 400 {object} string "Invalid request body"
// @Router /products [post]
func CreateProduct(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.DB.Create(&product)
    c.JSON(http.StatusCreated, product)
}

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
func UpdateProduct(c *gin.Context) {
    id := c.Param("id")
    var product models.Product
    if err := db.DB.First(&product, "ProductID = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.DB.Save(&product)
    c.JSON(http.StatusOK, product)
}

// @Summary Delete a product
// @Description Deletes a product by ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {string} string "Product deleted"
// @Failure 404 {string} string "Product not found"
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
    id := c.Param("id")
    if err := db.DB.Delete(&models.Product{}, "ProductID = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
