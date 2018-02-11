package storage

import (
	"fmt"
	"log"

	"github.com/crowdpower/fund/models"
	"github.com/crowdpower/fund/utils"
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

func (d *sqlDb) GetDeposits(username string, depositArgs *models.DepositArgs) ([]models.Deposit, error) {
	whereStatement, args := utils.SqlWhere([]utils.SqlCondition{
		utils.SqlCondition{"time", ">=", depositArgs.Oldest},
		utils.SqlCondition{"time", "<=", depositArgs.Newest},
		utils.SqlCondition{"amount", ">=", depositArgs.MinAmount},
		utils.SqlCondition{"amount", "<=", depositArgs.MaxAmount},
		utils.SqlCondition{"username", "=", username},
	})

	var pagination string
	if depositArgs.Count != 0 {
		pagination += fmt.Sprintf("LIMIT %v ", depositArgs.Count)
	}
	if depositArgs.Offset != 0 {
		pagination += fmt.Sprintf("OFFSET %v ", depositArgs.Offset)
	}

	rows, err := d.db.Query(`
        SELECT id, username, amount, time FROM Deposits `+whereStatement+`
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

	whereStatement, args := utils.SqlWhere([]utils.SqlCondition{
		utils.SqlCondition{"time", ">=", depositArgs.Oldest},
		utils.SqlCondition{"time", "<=", depositArgs.Newest},
		utils.SqlCondition{"amount", ">=", depositArgs.MinAmount},
		utils.SqlCondition{"amount", "<=", depositArgs.MaxAmount},
		utils.SqlCondition{"username", "=", username},
	})

	rows, err := d.db.Query(`
        SELECT SUM(amount) FROM Deposits `+whereStatement+`
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
