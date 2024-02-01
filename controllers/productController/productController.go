package productController

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/muhless/go-restapi-gin/models"
	"gorm.io/gorm"
	"net/http"
)

func Index(pc *gin.Context) {
	var Products []models.Product

	models.DB.Find(&Products)
	pc.JSON(http.StatusOK, gin.H{
		"products": Products,
	})
}

func Show(pc *gin.Context) {
	var Products models.Product
	id := pc.Param("id")

	if err := models.DB.First(&Products, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			pc.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		default:
			pc.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	pc.JSON(http.StatusOK, gin.H{
		"product": Products,
	})
}

func Create(pc *gin.Context) {
	var product models.Product

	//cara mengirim data ke json
	if err := pc.ShouldBindJSON(&product); err != nil {
		pc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return // return ini akan mengirim data ke var product
	}

	models.DB.Create(&product)
	pc.JSON(http.StatusOK, gin.H{"product": product})
}

func Update(pc *gin.Context) {
	var product models.Product

	id := pc.Param("id")

	if err := pc.ShouldBindJSON(&product); err != nil {
		pc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if models.DB.Model(&product).Where("id = ?", id).Updates(&product).RowsAffected == 0 {
		pc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "tidak dapat mengupdate product"})
		return
	}

	pc.JSON(http.StatusOK, gin.H{"product": "Data berhasil diperbarui"})
}

func Delete(pc *gin.Context) {
	var product models.Product

	var input struct {
		ID json.Number
	}

	if err := pc.ShouldBindJSON(&input); err != nil {
		pc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := input.ID.Int64() // membuat data string&int dapat dihapus
	if models.DB.Delete(&product, id).RowsAffected == 0 { // err apabila data kosong 
		pc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "tidak dapat menghapus product"})
		return
	}

	pc.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
