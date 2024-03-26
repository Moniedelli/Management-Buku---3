package model

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

type Book struct {
	Model
	ISBN    string `gorm:"not null" json:"isbn"`
	Penulis string `gorm:"not null" json:"penulis"`
	Tahun   uint   `gorm:"not null" json:"tahun"`
	Judul   string `gorm:"not null" json:"judul"`
	Gambar  string `gorm:"not null" json:"gambar"`
	Stok    uint   `gorm:"not null" json:"stok"`
}

func (cr *Book) Create(db *gorm.DB) error {
	err := db.
		Model(Book{}).
		Create(&cr).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (cr *Book) GetByID(db *gorm.DB) (Book, error) {
	res := Book{}

	err := db.
		Model(Book{}).
		Where("id = ?", cr.Model.ID).
		Take(&res).
		Error

	if err != nil {
		return Book{}, err
	}

	return res, nil
}

func (cr *Book) GetAll(db *gorm.DB) ([]Book, error) {
	res := []Book{}

	err := db.
		Model(Book{}).
		Find(&res).
		Error

	if err != nil {
		return []Book{}, err
	}

	return res, nil
}

func (cr *Book) UpdateOneByID(db *gorm.DB) error {
	err := db.
		Model(Book{}).
		Select("insb", "penulis", "tahun", "judul", "gambar", "stok").
		Where("id = ?", cr.Model.ID).
		Updates(map[string]any{
			"insb":    cr.ISBN,
			"penulis": cr.Penulis,
			"tahun":   cr.Tahun,
			"judul":   cr.Judul,
			"gambar":  cr.Gambar,
			"stok":    cr.Stok,
		}).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (cr *Book) DeleteByID(db *gorm.DB) error {
	err := db.
		Model(Book{}).
		Where("id = ?", cr.Model.ID).
		Delete(&cr).
		Error

	if err != nil {
		return err
	}

	return nil
}

func ImportCSVData(db *gorm.DB, filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	var importedCount int

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return importedCount, err
		}

		if len(record) >= 7 {
			isbn := strings.TrimSpace(record[0])
			penulis := strings.TrimSpace(record[1])
			tahun, err := strconv.Atoi(strings.TrimSpace(record[2]))
			if err != nil {
				return importedCount, err
			}
			judul := strings.TrimSpace(record[3])
			gambarURL := strings.TrimSpace(record[4])
			stok, err := strconv.Atoi(strings.TrimSpace(record[5]))
			if err != nil {
				return importedCount, err
			}

			var existingBook Book
			result := db.Where("isbn = ?", isbn).First(&existingBook)
			if result.RowsAffected > 0 {
				existingBook.Penulis = penulis
				existingBook.Tahun = uint(tahun)
				existingBook.Judul = judul
				existingBook.Stok = uint(stok)
				existingBook.Gambar = gambarURL

				err = db.Save(&existingBook).Error
				if err != nil {
					return importedCount, err
				}
			} else {
				book := &Book{
					ISBN:    isbn,
					Penulis: penulis,
					Tahun:   uint(tahun),
					Judul:   judul,
					Stok:    uint(stok),
					Gambar:  gambarURL,
				}

				if err := db.Create(book).Error; err != nil {
					return importedCount, err
				}
			}

			importedCount++
		}
	}

	return importedCount, nil
}

func (cr *Book) PrintPDF(db *gorm.DB, filePath string) error {
	books, err := cr.GetAll(db)
	if err != nil {
		return err
	}

	pdfFilePath := filepath.Join("pdf", filePath)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	for _, book := range books {
		pdf.Cell(0, 10, "ID: "+fmt.Sprintf("%d", book.ID))
		pdf.Ln(10)
		pdf.Cell(0, 10, "ISBN: "+book.ISBN)
		pdf.Ln(10)
		pdf.Cell(0, 10, "Penulis: "+book.Penulis)
		pdf.Ln(10)
		pdf.Cell(0, 10, "Tahun: "+fmt.Sprintf("%d", book.Tahun))
		pdf.Ln(10)
		pdf.Cell(0, 10, "Judul: "+book.Judul)
		pdf.Ln(10)
		pdf.Cell(0, 10, "Gambar: "+book.Gambar)
		pdf.Ln(10)
		pdf.Cell(0, 10, "Stok: "+fmt.Sprintf("%d", book.Stok))
		pdf.Ln(10)
		pdf.Ln(10)
	}

	err = pdf.OutputFileAndClose(pdfFilePath)
	if err != nil {
		return err
	}

	return nil
}
