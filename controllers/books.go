package controllers

import (
	"echo-auth-crud/config"
	"echo-auth-crud/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AddBook(c echo.Context) error {
	user := c.Get("user").(*models.User) // Mendapatkan informasi pengguna dari middleware GetUser

	book := models.Book{}
	if err := c.Bind(&book); err != nil {
		return err
	}

	book.UserID = user.Id // Set ID pengguna sebagai pemilik buku
	config.DB.Create(&book)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Book added successfully",
		"book":    book,
	})
}

func GetBooks(c echo.Context) error {
	books := []models.Book{}
	config.DB.Find(&books)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"books": books,
	})
}

func GetBook(c echo.Context) error {
	id := c.Param("id")
	book := models.Book{}

	if err := config.DB.Where("id = ?", id).First(&book).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"book": book,
	})
}

func UpdateBook(c echo.Context) error {

	id := c.Param("id")
	book := models.Book{}

	if err := config.DB.Where("id = ?", id).First(&book).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
	}

	if err := c.Bind(&book); err != nil {
		return err
	}

	config.DB.Save(&book)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Book updated successfully",
		"book":    book,
	})
}

func DeleteBook(c echo.Context) error {
	id := c.Param("id")
	book := models.Book{}

	if err := config.DB.Where("id = ?", id).First(&book).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
	}

	config.DB.Delete(&book)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Book deleted successfully",
	})
}
