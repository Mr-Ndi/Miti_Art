package controller

import (
	utils "MITI_ART/Utils"
	"MITI_ART/src/Vendors/dto"
	"MITI_ART/src/Vendors/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RegisterHandle godoc
// @Summary Register a new vendor
// @Description Vendor registration using token in query and form details in body
// @Tags Vendor
// @Accept json
// @Produce json
// @Param token query string true "Invitation token sent via email"
// @Param request body dto.VendorRegisterRequest true "Vendor registration data"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /vendor/register [post]
func RegisterHandle(c *gin.Context, db *gorm.DB) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Missing token in query"})
		return
	}

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

// UploadHandle godoc
// @Summary Upload a product
// @Description Upload a product with metadata and image (multipart form)
// @Tags Vendor
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param image formData file true "Product image"
// @Param name formData string true "Product name"
// @Param category formData string true "Category"
// @Param material formData string true "Material"
// @Param price formData number true "Price"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /vendor/upload [post]
func UploadHandle(c *gin.Context, db *gorm.DB) {
	vendorID := c.MustGet("user_id").(uuid.UUID)

	name := c.PostForm("name")
	category := c.PostForm("category")
	material := c.PostForm("material")
	priceStr := c.PostForm("price")

	price, err := parseFloat(priceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open image"})
		return
	}
	defer file.Close()

	imageURL, err := utils.UploadToCloudinary(file, fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// MyProduct godoc
// @Summary Get a single product by ID
// @Description Retrieves a product by ID for the vendor
// @Tags Vendor
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /vendor/my-product/{id} [get]
func MyProduct(c *gin.Context, db *gorm.DB) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	product, err := service.ProductByVendorID(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// MyOrders godoc
// @Summary Get vendor's product orders
// @Description Returns orders related to the vendor's products
// @Tags Vendor
// @Accept json
// @Produce json
// @Success 200 {array} dto.OrderResponse
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /vendor/required-product [get]
func MyOrders(c *gin.Context, db *gorm.DB) {
	vendorID := c.MustGet("user_id").(uuid.UUID)

	order, err := service.OrderByVendorID(db, vendorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Deletes a vendor's product by ID
// @Tags Vendor
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Security BearerAuth
// @Router /vendor/remove-product/{id} [delete]
func DeleteProduct(c *gin.Context, db *gorm.DB) {
	vendorID := c.MustGet("user_id").(uuid.UUID)
	idStr := c.Param("id")

	productID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := service.DeleteProductByID(db, productID, vendorID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// EditProduct godoc
// @Summary Edit product details
// @Description Allows vendor to edit their product details
// @Tags Vendor
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Param body body dto.EditProductRequest true "Product update fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Security BearerAuth
// @Router /vendor/edit-product/{id} [patch]
func EditProduct(c *gin.Context, db *gorm.DB) {
	vendorID := c.MustGet("user_id").(uuid.UUID)
	idStr := c.Param("id")

	productID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req dto.EditProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}

	if err := service.EditProductByID(db, productID, vendorID, updates); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// Helper: parse float from string
func parseFloat(s string) (float64, error) {
	var value float64
	_, err := fmt.Sscan(s, &value)
	return value, err
}
func VendorRegisterRoute(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		RegisterHandle(c, db)
	}
}
