package handler

import (
	"apu-ptt/internal/database"
	"apu-ptt/internal/model"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Function to fetch transactions by members not registered in ifuas with pagination, filter, and sortby
func fetchUnregisteredMemberTransactions(db *sql.DB, page, limit int, filter string, sortBy string) ([]model.Transaction, int, error) {
	// SQL query to fetch transactions by members not registered in ifuas with pagination, filter, and sortby
	query := `
        SELECT id, transaction_date, transaction_type, reference_no, sa_code, investor_fund_unit_ac_no,
               fund_code, amount_nominal, amount_unit, amount_all_units, fee_nominal, fee_unit,
               fee_percent, redm_payment_ac_sequential_code, redm_payment_bank_bic_code,
               redm_payment_bank_bi_member_code, redm_payment_ac_no, payment_date, transfer_type,
               sa_reference_no, created_at, updated_at
        FROM transactions
        WHERE investor_fund_unit_ac_no NOT IN (
            SELECT ifua_no FROM ifuas
        )
    `

	// Apply filter
	if filter != "" {
		query += " AND " + filter
	}

	// Apply sortby
	if sortBy != "" {
		query += " ORDER BY " + sortBy
	}

	// Apply pagination
	query += " LIMIT ? OFFSET ?"

	// Calculate offset based on pagination parameters
	offset := (page - 1) * limit

	// Execute the query
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Iterate through rows and scan into Transaction structs
	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		var createdAt, updatedAt sql.NullString
		err := rows.Scan(&t.ID, &t.TransactionDate, &t.TransactionType, &t.ReferenceNo,
			&t.SACode, &t.InvestorFundUnitACNo, &t.FundCode, &t.AmountNominal, &t.AmountUnit,
			&t.AmountAllUnits, &t.FeeNominal, &t.FeeUnit, &t.FeePercent, &t.RedmPaymentACSequential,
			&t.RedmPaymentBankBICCode, &t.RedmPaymentBankBIMember, &t.RedmPaymentACNo, &t.PaymentDate,
			&t.TransferType, &t.SAReferenceNo, &createdAt, &updatedAt)
		if err != nil {
			return nil, 0, err
		}
		if createdAt.Valid {
			t.CreatedAt = &createdAt.String
		}
		if updatedAt.Valid {
			t.UpdatedAt = &updatedAt.String
		}
		transactions = append(transactions, t)
	}

	// Query to count total number of unregistered member transactions
	countQuery := `
        SELECT COUNT(*) FROM transactions
        WHERE investor_fund_unit_ac_no NOT IN (
            SELECT ifua_no FROM ifuas
        )
    `
	var totalTransactions int
	err = db.QueryRow(countQuery).Scan(&totalTransactions)
	if err != nil {
		return nil, 0, err
	}

	// Calculate total number of pages
	totalPages := totalTransactions / limit
	if totalTransactions%limit != 0 {
		totalPages++
	}

	return transactions, totalPages, nil
}

func ImportUnregisteredMemberTransactionsHandler(c *gin.Context) {
	// Parse pagination parameters from query string
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Parse filter parameter from query string
	filter := c.DefaultQuery("filter", "")

	// Parse sortby parameter from query string
	sortBy := c.DefaultQuery("sortby", "")

	// Call the function to fetch unregistered member transactions
	transactions, totalPages, err := fetchUnregisteredMemberTransactions(database.DB, page, limit, filter, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set pagination information in response header
	c.Header("X-Total-Count", strconv.Itoa(totalPages))
	c.Header("X-Page", strconv.Itoa(page))
	c.Header("X-Limit", strconv.Itoa(limit))

	// Respond with fetched transactions
	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
	})
}

func fetchRegisteredMemberTransactions(db *sql.DB, page, limit int, filter string, sortBy string) ([]model.Transaction, int, error) {
	// SQL query to fetch transactions by members already registered in ifuas with pagination, filter, and sortby
	query := `
        SELECT id, transaction_date, transaction_type, reference_no, sa_code, investor_fund_unit_ac_no,
               fund_code, amount_nominal, amount_unit, amount_all_units, fee_nominal, fee_unit,
               fee_percent, redm_payment_ac_sequential_code, redm_payment_bank_bic_code,
               redm_payment_bank_bi_member_code, redm_payment_ac_no, payment_date, transfer_type,
               sa_reference_no, created_at, updated_at
        FROM transactions
        WHERE investor_fund_unit_ac_no IN (
            SELECT ifua_no FROM ifuas
        )
    `

	// Apply filter
	if filter != "" {
		query += " AND " + filter
	}

	// Apply sortby
	if sortBy != "" {
		query += " ORDER BY " + sortBy
	}

	// Apply pagination
	query += " LIMIT ? OFFSET ?"

	// Calculate offset based on pagination parameters
	offset := (page - 1) * limit

	// Execute the query
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Iterate through rows and scan into Transaction structs
	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		var createdAt, updatedAt sql.NullString
		err := rows.Scan(&t.ID, &t.TransactionDate, &t.TransactionType, &t.ReferenceNo,
			&t.SACode, &t.InvestorFundUnitACNo, &t.FundCode, &t.AmountNominal, &t.AmountUnit,
			&t.AmountAllUnits, &t.FeeNominal, &t.FeeUnit, &t.FeePercent, &t.RedmPaymentACSequential,
			&t.RedmPaymentBankBICCode, &t.RedmPaymentBankBIMember, &t.RedmPaymentACNo, &t.PaymentDate,
			&t.TransferType, &t.SAReferenceNo, &createdAt, &updatedAt)
		if err != nil {
			return nil, 0, err
		}
		if createdAt.Valid {
			t.CreatedAt = &createdAt.String
		}
		if updatedAt.Valid {
			t.UpdatedAt = &updatedAt.String
		}
		transactions = append(transactions, t)
	}

	// Query to count total number of registered member transactions
	countQuery := `
        SELECT COUNT(*) FROM transactions
        WHERE investor_fund_unit_ac_no IN (
            SELECT ifua_no FROM ifuas
        )
    `
	var totalTransactions int
	err = db.QueryRow(countQuery).Scan(&totalTransactions)
	if err != nil {
		return nil, 0, err
	}

	// Calculate total number of pages
	totalPages := totalTransactions / limit
	if totalTransactions%limit != 0 {
		totalPages++
	}

	return transactions, totalPages, nil
}

// Endpoint handler for fetching registered member transactions with pagination, filter, and sortby
func ImportRegisteredMemberTransactionsHandler(c *gin.Context) {
	// Parse pagination parameters from query string
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Parse filter parameter from query string
	filter := c.DefaultQuery("filter", "")

	// Parse sortby parameter from query string
	sortBy := c.DefaultQuery("sortby", "")

	// Call the function to fetch registered member transactions
	transactions, totalPages, err := fetchRegisteredMemberTransactions(database.DB, page, limit, filter, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set pagination information in response header
	c.Header("X-Total-Pages", strconv.Itoa(totalPages))

	// Respond with fetched transactions
	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
	})
}

func fetchAllTransactions(db *sql.DB, page, limit int, filter string, sortBy string) ([]model.Transaction, int, error) {
	// SQL query to fetch transactions by members not registered in ifuas with pagination, filter, and sortby
	query := `
        SELECT id, transaction_date, transaction_type, reference_no, sa_code, investor_fund_unit_ac_no,
               fund_code, amount_nominal, amount_unit, amount_all_units, fee_nominal, fee_unit,
               fee_percent, redm_payment_ac_sequential_code, redm_payment_bank_bic_code,
               redm_payment_bank_bi_member_code, redm_payment_ac_no, payment_date, transfer_type,
               sa_reference_no, created_at, updated_at
        FROM transactions
    `

	// Apply filter
	if filter != "" {
		query += " AND " + filter
	}

	// Apply sortby
	if sortBy != "" {
		query += " ORDER BY " + sortBy
	}

	// Apply pagination
	query += " LIMIT ? OFFSET ?"

	// Calculate offset based on pagination parameters
	offset := (page - 1) * limit

	// Execute the query
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Iterate through rows and scan into Transaction structs
	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		var createdAt, updatedAt sql.NullString
		err := rows.Scan(&t.ID, &t.TransactionDate, &t.TransactionType, &t.ReferenceNo,
			&t.SACode, &t.InvestorFundUnitACNo, &t.FundCode, &t.AmountNominal, &t.AmountUnit,
			&t.AmountAllUnits, &t.FeeNominal, &t.FeeUnit, &t.FeePercent, &t.RedmPaymentACSequential,
			&t.RedmPaymentBankBICCode, &t.RedmPaymentBankBIMember, &t.RedmPaymentACNo, &t.PaymentDate,
			&t.TransferType, &t.SAReferenceNo, &createdAt, &updatedAt)
		if err != nil {
			return nil, 0, err
		}
		if createdAt.Valid {
			t.CreatedAt = &createdAt.String
		}
		if updatedAt.Valid {
			t.UpdatedAt = &updatedAt.String
		}
		transactions = append(transactions, t)
	}

	// Query to count total number of unregistered member transactions
	countQuery := `
        SELECT COUNT(*) FROM transactions
    `
	var totalTransactions int
	err = db.QueryRow(countQuery).Scan(&totalTransactions)
	if err != nil {
		return nil, 0, err
	}

	// Calculate total number of pages
	totalPages := totalTransactions / limit
	if totalTransactions%limit != 0 {
		totalPages++
	}

	return transactions, totalPages, nil
}

func ListAllTransactions(c *gin.Context) {
	// Parse pagination parameters from query string
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Parse filter parameter from query string
	filter := c.DefaultQuery("filter", "")

	// Parse sortby parameter from query string
	sortBy := c.DefaultQuery("sortby", "")

	// Call the function to fetch unregistered member transactions
	transactions, totalPages, err := fetchAllTransactions(database.DB, page, limit, filter, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set pagination information in response header
	c.Header("X-Total-Count", strconv.Itoa(totalPages))
	c.Header("X-Page", strconv.Itoa(page))
	c.Header("X-Limit", strconv.Itoa(limit))

	// Respond with fetched transactions
	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
	})
}
