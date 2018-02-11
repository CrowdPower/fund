package storage

import (
	"fmt"
	"log"

	"github.com/crowdpower/fund/models"
	"github.com/crowdpower/fund/utils"
)

type payment interface {
	CreatePayment(payment *models.Payment) error
	GetPayment(username, id string) (*models.Payment, error)
	GetPayments(username string, paymentArgs *models.PaymentArgs) ([]models.Payment, error)
	GetPaymentsSum(username string, paymentArgs *models.PaymentArgs) (int, error)
}

func (d *sqlDb) CreatePayment(payment *models.Payment) error {
	_, err := d.db.Exec(`
        INSERT INTO Payments (id, username, amount, time, url) VALUES (?, ?, ?, ?, ?)
    `, payment.Id, payment.Username, payment.Amount, payment.Time, payment.Url)
	if err != nil {
		if err.Error() == InsufficientFundsMessage {
			err = &InsufficientFunds{}
		} else {
			log.Printf("error inserting payment %v into the database\n %v", payment, err)
		}
	}
	return err
}

func (d *sqlDb) GetPayment(username, id string) (*models.Payment, error) {
	rows, err := d.db.Query(`
        SELECT id, username, amount, time FROM Payments WHERE id = ? AND username = ?
    `, id, username)
	if err != nil {
		log.Printf("error reading payment %v from database for user %v\n%v", id, username, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		p := &models.Payment{}
		err := rows.Scan(&p.Id, &p.Username, &p.Amount, &p.Time)
		if err != nil {
			log.Printf("error parsing database rows\n%v", err)
			return nil, err
		}
		return p, nil
	}

	return nil, &NotFound{fmt.Sprintf("payment %v", id)}
}

func (d *sqlDb) GetPayments(username string, paymentArgs *models.PaymentArgs) ([]models.Payment, error) {
	whereStatement, args := utils.SqlWhere([]utils.SqlCondition{
		utils.SqlCondition{"time", ">=", paymentArgs.Oldest},
		utils.SqlCondition{"time", "<=", paymentArgs.Newest},
		utils.SqlCondition{"amount", ">=", paymentArgs.MinAmount},
		utils.SqlCondition{"amount", "<=", paymentArgs.MaxAmount},
		utils.SqlCondition{"url", "LIKE", "%" + paymentArgs.Url + "%"},
		utils.SqlCondition{"username", "=", username},
	})

	var pagination string
	if paymentArgs.Count != 0 {
		pagination += fmt.Sprintf("LIMIT %v ", paymentArgs.Count)
	}
	if paymentArgs.Offset != 0 {
		pagination += fmt.Sprintf("OFFSET %v", paymentArgs.Offset)
	}

	rows, err := d.db.Query(`
        SELECT id, username, amount, time, url FROM Payments `+whereStatement+`
    `+pagination, args...)
	if err != nil {
		log.Printf("error reading payments from database for user %v\n%v", username, err)
		return nil, err
	}
	defer rows.Close()

	payments := []models.Payment{}
	for rows.Next() {
		p := models.Payment{}
		err := rows.Scan(&p.Id, &p.Username, &p.Amount, &p.Time, &p.Url)
		if err != nil {
			log.Printf("error parsing database rows\n%v", err)
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (d *sqlDb) GetPaymentsSum(username string, paymentArgs *models.PaymentArgs) (int, error) {
	var sum int

	whereStatement, args := utils.SqlWhere([]utils.SqlCondition{
		utils.SqlCondition{"time", ">=", paymentArgs.Oldest},
		utils.SqlCondition{"time", "<=", paymentArgs.Newest},
		utils.SqlCondition{"amount", ">=", paymentArgs.MinAmount},
		utils.SqlCondition{"amount", "<=", paymentArgs.MaxAmount},
		utils.SqlCondition{"url", "LIKE", "%" + paymentArgs.Url + "%"},
		utils.SqlCondition{"username", "=", username},
	})

	rows, err := d.db.Query(`
        SELECT SUM(amount) FROM Payments `+whereStatement+`
    `, args...)
	if err != nil {
		log.Printf("error summing payments from database for user %v\n%v", username, err)
		return sum, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&sum)
		if err != nil {
			log.Printf("error parsing database rows\n%v", err)
			return sum, err
		}
	}

	return sum, nil
}
