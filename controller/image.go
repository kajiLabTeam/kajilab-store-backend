package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 商品画像の取得
func GetProductImage(c *gin.Context) {
	reqPath := c.Param("path")
	fmt.Println(reqPath)

	path := "./images/products/" + reqPath
	c.Header("Cache-Control", "no-store")

	c.File(path)
}

func UploadProductImage(c *gin.Context) {
	fileName := c.Param("imgFileName")

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "file required")
		return
	}

	path := fmt.Sprintf("./images/products/%s", fileName)

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "error")
		return
	}

	c.JSON(http.StatusOK, "success")
}
