package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"MITI_ART/Vendors/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestVendorRegisterHandle(t *testing.T) {
	router := gin.Default()

	// Normally, you would mock a real DB here
	var fakeDB *gorm.DB = nil // Just pass nil for now to keep simple

	router.POST("/vendor/register", controller.VendorRegisterRoute(fakeDB))

	jsonBody := []byte(`{
        "VendorPassword": "securePassword#1",
        "VendorTin": 999991130,
        "ShopName": "Viva bzs grp"
    }`)

	req, _ := http.NewRequest("POST", "/vendor/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 but got %d", w.Code)
	}
}
