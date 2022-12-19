package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/viniosilva/socialassistanceapi/internal/configuration"
	"github.com/viniosilva/socialassistanceapi/internal/exception"
)

//go:generate mockgen -destination ../../mock/donate_resource_repository_mock.go -package mock . DonateResourceRepository
type DonateResourceRepository interface {
	Donate(ctx context.Context, resourceID, addressID int, quantity float64) error
	Return(ctx context.Context, resourceID int) error
}

type DonateResourceRepositoryImpl struct {
	DB configuration.MySQL
}

func (impl *DonateResourceRepositoryImpl) Donate(ctx context.Context, resourceID, addressID int, quantity float64) error {
	tx, err := impl.DB.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	nowMysql := time.Now().Format("2006-01-02T15:04:05")

	res, err := tx.QueryContext(ctx, "SELECT quantity FROM resources WHERE id = ?", resourceID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	var dbQuantity float64
	found := false
	for res.Next() {
		found = true
		if err = res.Scan(&dbQuantity); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	if !found {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return &exception.NotFoundException{Err: fmt.Errorf("resource %d not found", resourceID)}
	}
	if dbQuantity-quantity < 0 {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return &exception.NegativeException{Err: fmt.Errorf("resource %d quantity is %.1f", resourceID, dbQuantity)}
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO resources_to_addresses (created_at, resource_id, address_id, quantity)
		VALUES (?, ?, ?, ?)
	`, nowMysql, resourceID, addressID, quantity)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1452 {
			return &exception.NotFoundException{Err: fmt.Errorf("address %d not found", addressID)}
		}
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE resources
		SET updated_at = ?,
			quantity = quantity - ?
		WHERE id = ?
	`, nowMysql, quantity, resourceID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (impl *DonateResourceRepositoryImpl) Return(ctx context.Context, resourceID int) error {
	tx, err := impl.DB.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	nowMysql := time.Now().Format("2006-01-02T15:04:05")

	res, err := tx.QueryContext(ctx, "SELECT quantity FROM resources_to_addresses WHERE resource_id = ?", resourceID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	var dbQuantity float64
	found := false
	for res.Next() {
		found = true
		if err = res.Scan(&dbQuantity); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	if !found {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return &exception.NotFoundException{Err: fmt.Errorf("resource %d not found", resourceID)}
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM resources_to_addresses WHERE resource_id = ?", resourceID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE resources
		SET updated_at = ?,
			quantity = quantity + ?
		WHERE id = ?
	`, nowMysql, dbQuantity, resourceID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
