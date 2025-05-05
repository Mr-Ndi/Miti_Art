package controller

import (
	utils "MITI_ART/Utils"
	"MITI_ART/Vendors/service"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func VendorRegisterRoute(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		RegisterHandle(c, db)
	}
}

func RegisterHandle(c *gin.Context, db *gorm.DB) {
	vendorToken := c.GetHeader("Authorization")
	tokenParts := strings.Split(vendorToken, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid token format"})
		return
	}

	token := tokenParts[1]
	payload, err := utils.ValidateToken(token)
	// fmt.Printf("Decoded payload: %+v\n", payload)
	if err != nil {
		fmt.Println("Token validation error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid or expired token"})
		return
	}

	VendorEmail, emailOk := payload["VendorEmail"].(string)
	VendorFirstName, firstNameOk := payload["VendorFirstName"].(string)
	VendorOtherName, otherNameOk := payload["VendorOtherName"].(string)
	role, roleOk := payload["role"].(string)

	if !emailOk || !firstNameOk || !otherNameOk || !roleOk {
		fmt.Println("Token payload missing required fields")
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid token payload"})
		return
	}

	exp, expOk := payload["exp"].(float64)
	if !expOk {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid token payload cz it is missing expiration time"})
		return
	}

	currentTimestamp := float64(time.Now().Unix())
	if exp < currentTimestamp {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Token has expired"})
		return
	}

	var req struct {
		VendorPassword string `json:"vendorPassword" binding:"required"`
		VendorTin      int    `json:"vendorTin" binding:"required"`
		ShopName       string `json:"ShopName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	message, err := service.RegisterVendor(db, VendorEmail, VendorFirstName, VendorOtherName, req.VendorPassword, role, req.VendorTin, req.ShopName)
	// fmt.Println("Calling RegisterVendor with:", VendorEmail, VendorFirstName, VendorOtherName, role, req.VendorPassword, req.VendorTin, req.ShopName)

	if err != nil {
		fmt.Println("Error from RegisterVendor Services:", err)
		if err.Error() == "user with that email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      message,
		"Vendor email": VendorEmail,
	})
}

func UploadHandle(c *gin.Context, db *gorm.DB) {
	// Get email set by middleware
	vendorEmail, exists := c.Get("vendor_email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	vendorID, err := utils.GetVendorIDByEmail(db, vendorEmail.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle file upload
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Println("Error receiving file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image not found"})
		return
	}
	defer file.Close()

	imagePath, err := utils.UploadImage(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse form fields
	name := c.PostForm("name")
	category := c.PostForm("category")
	material := c.PostForm("material")
	priceStr := c.PostForm("price")

	if name == "" || category == "" || material == "" || priceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
		return
	}

	message, err := service.RegisterProduct(db, vendorID, name, price, category, material, imagePath)
	if err != nil {
		fmt.Println("Error from RegisterProduct Services:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      message,
		"vendor_email": vendorEmail,
	})
}

func MyProduct(c *gin.Context, db *gorm.DB) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	products, err := service.ProductByVendorID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func MyOrders(c *gin.Context, db *gorm.DB) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	products, err := service.OrderByVendorID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func DeleteById(c *gin.Context, db *gorm.DB) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
	}
	err = service.DeleteById(db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})

}

func EditProduct(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var body struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updateData := map[string]interface{}{}
	if body.Name != "" {
		updateData["name"] = body.Name
	}
	if body.Price > 0 {
		updateData["price"] = body.Price
	}

	err = service.EditProductByID(db, productID, updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}
