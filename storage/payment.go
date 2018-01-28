package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/crowdpower/fund/models"
)

type PaymentArgs struct {
	Oldest    time.Time
	Newest    time.Time
	MinAmount int
	MaxAmount int
	Url       string
}

type payment interface {
	CreatePayment(payment *models.Payment) error
	GetPayment(username, id string) (*models.Payment, error)
	GetPayments(username string, paymentArgs *PaymentArgs) ([]models.Payment, error)
	GetPaymentsSum(username string, paymentArgs *PaymentArgs) (int, error)
}

func (d *sqlDb) CreatePayment(payment *models.Payment) error {
	_, err := d.db.Exec(`
        INSERT INTO Payments (id, username, amount, time, url) VALUES (?, ?, ?, ?, ?)
    `, payment.Id, payment.Username, payment.Amount, payment.Time, payment.Url)
	if err != nil {
		log.Printf("error inserting payment %v into the database\n %v", payment, err)
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

func getPaymentsConditions(paymentArgs *PaymentArgs) ([]string, []interface{}) {
	conditions := []string{}
	args := make([]interface{}, 0)

	if !paymentArgs.Oldest.IsZero() {
		conditions = append(conditions, "time >= ?")
		args = append(args, paymentArgs.Oldest.Format(TimeFormat))
	}

	if !paymentArgs.Newest.IsZero() {
		conditions = append(conditions, "time <= ?")
		args = append(args, paymentArgs.Newest.Format(TimeFormat))
	}

	if paymentArgs.MaxAmount != 0 {
		conditions = append(conditions, "amount <= ?")
		args = append(args, paymentArgs.MaxAmount)
	}

	if paymentArgs.MinAmount != 0 {
		conditions = append(conditions, "amount >= ?")
		args = append(args, paymentArgs.MinAmount)
	}

	if paymentArgs.Url != "" {
		conditions = append(conditions, "url LIKE ?")
		args = append(args, "%" + paymentArgs.Url + "%")
	}

	return conditions, args
}

func (d *sqlDb) GetPayments(username string, paymentArgs *PaymentArgs) ([]models.Payment, error) {
	conditions, args := getPaymentsConditions(paymentArgs)
	conditions = append(conditions, "username = ?")
	args = append(args, username)

	rows, err := d.db.Query(`
        SELECT id, username, amount, time, url FROM Payments
		WHERE `+strings.Join(conditions, " AND ")+`
    `, args...)
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

func (d *sqlDb) GetPaymentsSum(username string, paymentArgs *PaymentArgs) (int, error) {
	var sum int

	conditions, args := getPaymentsConditions(paymentArgs)
	conditions = append(conditions, "username = ?")
	args = append(args, username)

	rows, err := d.db.Query(`
        SELECT SUM(amount) FROM Payments
		WHERE `+strings.Join(conditions, " AND ")+`
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
