package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DataPackageExcel(c *gin.Context, data []byte) {
	c.Header("Content-Type", "application/vnd.ms-excel")
	c.Header("Pragma", "public")
	c.Header("Cache-Control", "no-store")
	c.Header("Cache-Control", "max-age=0")
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Data(http.StatusOK, "application/vnd.ms-excel", data)
}

func DataPackageDbf(c *gin.Context, data []byte) {
	c.Header("Content-Type", "application/octet-stream;charset=UTF-8")
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Data(http.StatusOK, "application/octet-stream;charset=UTF-8", data)
}

func DataPackageZip(c *gin.Context, data []byte, fileName string) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Data(http.StatusOK, "application/octet-stream; charset=UTF-8", data)
}
