package model_test

import (
	"fmt"
	"os"
	"perpustakaan/miniproject/config"
	"perpustakaan/miniproject/model"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func Init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("env not found, using global env")
	}
	config.OpenDB()
}

func TestCreateBook(t *testing.T) {
	Init()

	bookData := model.Book{
		ISBN:    "978",
		Penulis: "Del",
		Tahun:   1998,
		Judul:   "Toyota",
		Gambar:  "crown",
		Stok:    456,
	}
	err := bookData.Create(config.Mysql.DB)
	assert.Nil(t, err)

	fmt.Println(bookData.ID)
}

func TestGetByID(t *testing.T) {
	Init()

	bookData := model.Book{
		Model: model.Model{
			ID: 1,
		},
	}

	data, err := bookData.GetByID(config.Mysql.DB)
	assert.Nil(t, err)

	fmt.Println(data)
}

func TestGetAll(t *testing.T) {
	Init()

	bookData := model.Book{
		ISBN:    "978",
		Penulis: "Del",
		Tahun:   1998,
		Judul:   "Toyota",
		Gambar:  "crown",
		Stok:    456,
	}

	err := bookData.Create(config.Mysql.DB)
	assert.Nil(t, err)

	res, err := bookData.GetAll(config.Mysql.DB)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(res), 1)

	fmt.Println(res)
}

func TestDeleteByID(t *testing.T) {
	Init()

	bookData := model.Book{
		Model: model.Model{
			ID: 1,
		},
	}

	err := bookData.DeleteByID(config.Mysql.DB)
	assert.Nil(t, err)
}

func TestImportCSVData(t *testing.T) {
	Init()

	// Path file CSV yang akan diimpor
	filePath := "../sample_books.csv" // Ganti dengan path yang benar

	// Import data dari file CSV
	_, err := model.ImportCSVData(config.Mysql.DB, filePath)
	assert.Nil(t, err)

	// Cek apakah data berhasil diimpor
	var books []model.Book
	err = config.Mysql.DB.Find(&books).Error
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(books), 1) // Memastikan setidaknya ada satu buku yang diimpor
}

func TestPrintPDF(t *testing.T) {
	Init()

	// Insert dummy data
	dummyBook := model.Book{
		ISBN:    "1234567890",
		Penulis: "John Doe",
		Tahun:   2022,
		Judul:   "Sample Book",
		Gambar:  "sample.jpg",
		Stok:    10,
	}
	err := dummyBook.Create(config.Mysql.DB)
	if err != nil {
		t.Fatalf("failed to create dummy data: %v", err)
	}

	// Path untuk menyimpan file PDF
	filePath := "sample_books.pdf"

	// Print PDF
	err = dummyBook.PrintPDF(config.Mysql.DB, filePath)
	if err != nil {
		t.Fatalf("failed to print PDF: %v", err)
	}

	// Check if file exists
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Fatalf("PDF file not found")
	}

	// Clean up
	os.Remove(filePath)
}
