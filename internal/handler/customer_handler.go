package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"apu-ptt/internal/database"
	"apu-ptt/internal/model"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func fetchCustomersWithMultipleTransactions(db *sql.DB, page, limit int, filter, sortBy string) ([]model.CustomerTransaction, int, error) {
	// Calculate start date one year ago
	startDate := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")

	// SQL query to fetch customers who have made transactions more than twice in a year with pagination, filter, and sorting
	query := `
        SELECT investor_fund_unit_ac_no, SUM(amount_nominal) AS total_amount_nominal, COUNT(*) AS transaction_count
        FROM transactions
        WHERE transaction_date BETWEEN ? AND NOW()
        GROUP BY investor_fund_unit_ac_no
        HAVING transaction_count > 2
    `

	// Apply filter
	if filter != "" {
		query += " AND " + filter
	}

	// Apply sorting
	if sortBy != "" {
		query += " ORDER BY " + sortBy
	}

	// Apply pagination
	query += " LIMIT ? OFFSET ?"

	// Calculate offset based on pagination parameters
	offset := (page - 1) * limit

	// Execute the query
	rows, err := db.Query(query, startDate, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Iterate through rows and scan into CustomerTransaction structs
	var customerTransactions []model.CustomerTransaction
	for rows.Next() {
		var customer string
		var totalAmountNominal float64
		var transactionCount int
		err := rows.Scan(&customer, &totalAmountNominal, &transactionCount)
		if err != nil {
			return nil, 0, err
		}
		customerTransactions = append(customerTransactions, model.CustomerTransaction{
			Customer:           customer,
			TotalAmountNominal: totalAmountNominal,
			TransactionCount:   transactionCount,
		})
	}

	// Query to count total number of customers with multiple transactions
	countQuery := `
        SELECT COUNT(*) FROM (
            SELECT investor_fund_unit_ac_no
            FROM transactions
            WHERE transaction_date BETWEEN ? AND NOW()
            GROUP BY investor_fund_unit_ac_no
            HAVING COUNT(*) > 2
        ) AS subquery
    `
	var totalCustomers int
	err = db.QueryRow(countQuery, startDate).Scan(&totalCustomers)
	if err != nil {
		return nil, 0, err
	}

	// Calculate total number of pages
	totalPages := totalCustomers / limit
	if totalCustomers%limit != 0 {
		totalPages++
	}

	return customerTransactions, totalPages, nil
}

func CustomersWithMultipleTransactionsHandler(c *gin.Context) {
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

	// Call the function to fetch customers with multiple transactions
	customerTransactions, totalPages, err := fetchCustomersWithMultipleTransactions(database.DB, page, limit, filter, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set pagination information in response header
	c.Header("X-Total-Pages", strconv.Itoa(totalPages))

	// Respond with fetched customers
	c.JSON(http.StatusOK, gin.H{
		"customers": customerTransactions,
	})
}
