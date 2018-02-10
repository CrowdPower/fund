package storage

import (
	"fmt"
	"log"
	"strings"

	"github.com/crowdpower/fund/models"
)

type deposit interface {
	CreateDeposit(deposit *models.Deposit) error
	GetDeposit(username, id string) (*models.Deposit, error)
	GetDeposits(username string, depositArgs *models.DepositArgs) ([]models.Deposit, error)
	GetDepositsSum(username string, depositArgs *models.DepositArgs) (int, error)
}

func (d *sqlDb) CreateDeposit(deposit *models.Deposit) error {
	_, err := d.db.Exec(`
        INSERT INTO Deposits (id, username, amount, time) VALUES (?, ?, ?, ?)
    `, deposit.Id, deposit.Username, deposit.Amount, deposit.Time)
	if err != nil {
		log.Printf("error inserting deposit %v into the database\n %v", deposit, err)
	}
	return err
}

func (d *sqlDb) GetDeposit(username, id string) (*models.Deposit, error) {
	rows, err := d.db.Query(`
        SELECT id, username, amount, time FROM Deposits WHERE id = ? AND username = ?
    `, id, username)
	if err != nil {
		log.Printf("error reading deposit %v from database for user %v\n%v", id, username, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		d := &models.Deposit{}
		err := rows.Scan(&d.Id, &d.Username, &d.Amount, &d.Time)
		if err != nil {
			log.Printf("error parsing database rows\n%v", err)
			return nil, err
		}
		return d, nil
	}

	return nil, &NotFound{fmt.Sprintf("deposit %v", id)}
}

func getDepositsConditions(depositArgs *models.DepositArgs) ([]string, []interface{}) {
	conditions := []string{}
	args := make([]interface{}, 0)

	if !depositArgs.Oldest.IsZero() {
		conditions = append(conditions, "time >= ?")
		args = append(args, depositArgs.Oldest.Format(TimeFormat))
	}

	if !depositArgs.Newest.IsZero() {
		conditions = append(conditions, "time <= ?")
		args = append(args, depositArgs.Newest.Format(TimeFormat))
	}

	if depositArgs.MaxAmount != 0 {
		conditions = append(conditions, "amount <= ?")
		args = append(args, depositArgs.MaxAmount)
	}

	if depositArgs.MinAmount != 0 {
		conditions = append(conditions, "amount >= ?")
		args = append(args, depositArgs.MinAmount)
	}

	return conditions, args
}

func (d *sqlDb) GetDeposits(username string, depositArgs *models.DepositArgs) ([]models.Deposit, error) {
	conditions, args := getDepositsConditions(depositArgs)
	conditions = append(conditions, "username = ?")
	args = append(args, username)

	var pagination string
	if depositArgs.Count != 0 {
		pagination += fmt.Sprintf("LIMIT %v ", depositArgs.Count)
	}
	if depositArgs.Offset != 0 {
		pagination += fmt.Sprintf("OFFSET %v", depositArgs.Offset)
	}

	rows, err := d.db.Query(`
        SELECT id, username, amount, time FROM Deposits
		WHERE `+strings.Join(conditions, " AND ")+`
    `+pagination, args...)
	if err != nil {
		log.Printf("error reading deposits from database for user %v\n%v", username, err)
		return nil, err
	}
	defer rows.Close()

	deposits := []models.Deposit{}
	for rows.Next() {
		d := models.Deposit{}
		err := rows.Scan(&d.Id, &d.Username, &d.Amount, &d.Time)
		if err != nil {
			log.Printf("error parsing database rows\n%v", err)
			return nil, err
		}
		deposits = append(deposits, d)
	}

	return deposits, nil
}

func (d *sqlDb) GetDepositsSum(username string, depositArgs *models.DepositArgs) (int, error) {
	var sum int

	conditions, args := getDepositsConditions(depositArgs)
	conditions = append(conditions, "username = ?")
	args = append(args, username)

	rows, err := d.db.Query(`
        SELECT SUM(amount) FROM Deposits
		WHERE `+strings.Join(conditions, " AND ")+`
    `, args...)
	if err != nil {
		log.Printf("error summing deposits from database for user %v\n%v", username, err)
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
