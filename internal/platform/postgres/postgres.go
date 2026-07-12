package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(db *gorm.DB) *Database {
	return &Database{DB: db}
}

func (d *Database) beginTxn(ctx context.Context) (*Database, error) {
	tx := d.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("begin transaction: %w", tx.Error)
	}
	return &Database{DB: tx}, nil
}

func (d *Database) commitTx() error {
	if err := d.DB.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (d *Database) rollBackTx() error {
	if err := d.DB.Rollback().Error; err != nil {
		return fmt.Errorf("rollback transaction: %w", err)
	}
	return nil
}

func (d *Database) WithTx(ctx context.Context, fn func(tx *Database) error) error {
	tx, err := d.beginTxn(ctx)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.rollBackTx(); rbErr != nil {
			return fmt.Errorf("%w (rollback error: %v)", err, rbErr)
		}
		return err
	}

	return tx.commitTx()
}
