package handler

import (
	"encoding/csv"
	"net/http"

	"apu-ptt/internal/database"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func ImportTeroristHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Read data from Excel files
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	xlsx, err := excelize.OpenReader(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Extract data from an Excel file and insert it into a database
	rows := xlsx.GetRows("Sheet1")
	for i, row := range rows {
		if i == 0 {
			// Skip header row
			continue
		}

		densus_code := row[0]
		name := row[1]
		description := row[2]
		investor_type := row[3]
		place_of_birth := row[4]
		date_of_birth := row[5]
		nationality := row[6]
		address := row[7]

		// Save data to the database
		_, err := database.DB.Exec("INSERT INTO terorism_lists (densus_code, name, description, investor_type, place_of_birth, date_of_birth, nationality, address) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", densus_code, name, description, investor_type, place_of_birth, date_of_birth, nationality, address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	fileContent.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully"})
}

func ImportPpspmtHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Read data from Excel files
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	xlsx, err := excelize.OpenReader(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Extract data from an Excel file and insert it into a database
	rows := xlsx.GetRows("Sheet1")
	for i, row := range rows {
		if i == 0 {
			// Skip header row
			continue
		}

		reference := row[0]
		investor_type := row[1]
		name := row[2]
		degree := row[3]
		employment := row[4]
		date_of_birth := row[5]
		place_of_birth := row[6]
		_, err := database.DB.Exec("INSERT INTO pspm_lists (reference, investor_type, name, degree, employment, date_of_birth, place_of_birth) VALUES (?, ?, ?, ?, ?, ?, ?)", reference, investor_type, name, degree, employment, date_of_birth, place_of_birth)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	fileContent.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully"})
}

func ImportTransactionHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open the uploaded file
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileContent.Close()

	// Parse the CSV file
	reader := csv.NewReader(fileContent)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	// Prepare the INSERT statement
	stmt, err := tx.Prepare("INSERT INTO transactions (transaction_date, transaction_type, reference_no, sa_code, investor_fund_unit_ac_no, fund_code, amount_nominal, amount_unit, amount_all_units, fee_nominal, fee_unit, fee_percent, redm_payment_ac_sequential_code, redm_payment_bank_bic_code, redm_payment_bank_bi_member_code, redm_payment_ac_no, payment_date, transfer_type, sa_reference_no) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	// Insert each record into the database
	for _, row := range records {
		_, err := stmt.Exec(row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9], row[10], row[11], row[12], row[13], row[14], row[15], row[16], row[17], row[18])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully"})
}
