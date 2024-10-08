package main

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// Gofile API URL
const GofileAPIURL = "https://api.gofile.io"

// JWT Secret
var jwtSecret = []byte("your_jwt_secret")

// Struct for JWT Claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// UploadFileResponse Struct
type UploadFileResponse struct {
	Status string `json:"status"`
	Data   struct {
		DownloadPage string `json:"downloadPage"`
	} `json:"data"`
}

func main() {
	r := gin.Default()

	// Middleware para verificar JWT
	r.Use(AuthMiddleware())

	// Rutas
	r.POST("/upload", uploadFile)
	r.GET("/download/:fileId", downloadFile)
	r.GET("/files", listFiles)
	r.DELETE("/files/:fileId", deleteFile)
	r.POST("/share", shareFile)

	// Puerto
	r.Run(":8080")
}

// Middleware para autenticación JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

// POST /upload - Subir archivo
func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	// Abrir archivo
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileContent.Close()

	client := resty.New()
	resp, err := client.R().
		SetFileReader("file", file.Filename, fileContent).
		Post(GofileAPIURL + "/uploadFile")

	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed"})
		return
	}

	var result UploadFileResponse
	if err := resp.Unmarshal(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"downloadPage": result.Data.DownloadPage})
}

// GET /download/:fileId - Descargar archivo
func downloadFile(c *gin.Context) {
	fileId := c.Param("fileId")

	c.Redirect(http.StatusFound, GofileAPIURL+"/download/"+fileId)
}

// GET /files - Listar archivos (implementación simplificada)
func listFiles(c *gin.Context) {
	// Puedes agregar lógica aquí para listar archivos del usuario actual.
	// La API de Gofile tiene opciones para obtener archivos de la cuenta, pero requiere token de cuenta.
	c.JSON(http.StatusOK, gin.H{"message": "List of files (Not implemented)"})
}

// DELETE /files/:fileId - Eliminar archivo
func deleteFile(c *gin.Context) {
	fileId := c.Param("fileId")

	client := resty.New()
	resp, err := client.R().
		SetQueryParam("fileId", fileId).
		Post(GofileAPIURL + "/deleteFile")

	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

// POST /share - Compartir archivo por email (simulado)
func shareFile(c *gin.Context) {
	var request struct {
		FileID string `json:"fileId" binding:"required"`
		Email  string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Aquí podrías implementar el envío de correos usando un servicio como SendGrid o Mailgun.
	// Por ahora, simulamos el compartir.
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Shared file %s with %s", request.FileID, request.Email)})
}
