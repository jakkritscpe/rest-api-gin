package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	host := "dpg-cecl59g2i3mr5pdgpqe0-a.singapore-postgres.render.com"
	user := "portfolio_y3pt_user"
	password := "lIDo0AGigiP5sZapOGevRXaEybBKsXJV"
	dbname := "portfolio_y3pt"
	port := "5432"

	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=prefer",
		host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Book{})

	handler := newHandler(db)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/books", handler.listBooksHandler)
	r.POST("/books", handler.createBookHandler)
	r.DELETE("/books/:id", handler.deleteBookHandler)

	r.Run()

}

type Handler struct {
	db *gorm.DB
}

func newHandler(db *gorm.DB) *Handler {
	return &Handler{db}
}

type Book struct {
	ID     uint   `gorm:"primaryKey;autoIncrement;not null"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func (h *Handler) listBooksHandler(c *gin.Context) {
	var books []Book

	if result := h.db.Find(&books); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &books)
}

func (h *Handler) createBookHandler(c *gin.Context) {
	var book Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if result := h.db.Create(&book); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, &book)
}

func (h *Handler) deleteBookHandler(c *gin.Context) {
	id := c.Param("id")

	if result := h.db.Delete(&Book{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
