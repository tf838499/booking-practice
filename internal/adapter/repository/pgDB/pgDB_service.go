package PostgresDB

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type pgDB interface {
	Querier
	CreateOrderAndSetSeatIsSoldTx(ctx context.Context, arg CreateOrderAndSetSeatIsSoldParams) ([]Ticket, error)
	SelectSeatsListAndUpdateSeatStatusTx(ctx context.Context, arg SelectSeatsListByEventAndSectionParams) ([]Seat, error)
	DeletTicketAndUpdateSeatIsAvailableStatusTx(ctx context.Context, orderTradeNo string, seatStaus bool) ([]Seat, error)
	// TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	// CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	// VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}
type PostgresRepository struct {
	*Queries
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) pgDB {
	return &PostgresRepository{
		Queries: New(db),
		db:      db,
	}
}
func (store *PostgresRepository) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// func (store *PostgresRepository) execTx(ctx context.Context, fn func(*Queries) error) error {
// 	tx, err := store.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}

// 	q := New(tx)
// 	err = fn(q)
// 	if err != nil {
// 		fmt.Println("Error occur, Start Rollback")
// 		if rbErr := tx.Rollback(); rbErr != nil {
// 			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
// 		}
// 		return err
// 	}

// 	return tx.Commit()
// }

// func (store *PostgresRepository) InsertGoodsWithTx(ctx context.Context, goodsParams []InsertGoodsParams) error {
// 	err := store.execTx(ctx, func(q *Queries) error {
// 		var err error
// 		for i := range goodsParams {
// 			err = store.InsertGoods(ctx, InsertGoodsParams{
// 				ImageName: sql.NullString{String: goodsParams[i].ImageName.String, Valid: true},
// 				Descript:  sql.NullString{String: goodsParams[i].Descript.String, Valid: true},
// 				Price:     sql.NullInt64{Int64: goodsParams[i].Price.Int64, Valid: true},
// 				Class:     sql.NullString{String: goodsParams[i].Class.String, Valid: true},
// 			})
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// 	return err
// }
