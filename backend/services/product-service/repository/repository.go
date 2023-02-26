package repository

import (
	"context"
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"strconv"
	"sync"
	"testLmd/libs/errors"
	"testLmd/services/product-service/config"
	product_service "testLmd/services/product-service/proto/product-service"
)

type Database interface {
	ReservationProduct(ctx context.Context, id []int64, WarehouseId int64) (map[string]string, error)
	CancelReservationProduct(ctx context.Context, id []int64, WarehouseId int64) (map[string]string, error)
	GetAllProducts(ctx context.Context, id int64) ([]*product_service.Product, error)
}

type DatabaseConn struct {
	conn *sql.DB
}

func NewDatabase() (*DatabaseConn, error) {
	conf := config.ProductConfig()
	connStr := conf.PgConnString

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		errors.HandleFatalError(err, "failed to connect to postgres")
	}

	err = db.Ping()

	if err != nil {
		errors.HandleFatalError(err, "failed to connect to postgres")
	}

	return &DatabaseConn{conn: db}, err
}

func (db DatabaseConn) ReservationProduct(
	ctx context.Context,
	id []int64,
	WarehouseId int64,
) (map[string]string, error) {
	q := `
			UPDATE products 
			SET quantity = quantity - 1,
			    in_reserve = in_reserve + 1
			WHERE id = $1 AND quantity > 0 AND warehouse_id = $2
			`

	qSelect := `
				SELECT (quantity > 0) 
				FROM products
				WHERE id = $1
				`

	var wg sync.WaitGroup
	res := make(map[string]string)
	warehouseErr := func(k int, v int64, tx *sql.Tx) {
		mapK := "Товара номер " + strconv.Itoa(k+1) + " с id - " + strconv.Itoa(int(v))
		res[mapK] = " нет на складе"
		tx.Rollback()
	}

	warehouseOk := func(k int, v int64, tx *sql.Tx) {
		mapK := "Товар номер " + strconv.Itoa(k+1) + " с id - " + strconv.Itoa(int(v))
		res[mapK] = " поставлен в резерв"
	}

	for k, v := range id {
		wg.Add(1)
		go func(k int, v int64) {
			defer wg.Done()
			tx, err := db.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
			if err != nil {
				errors.LogError(err)
				return
			}

			var enough bool
			err = tx.QueryRowContext(ctx, qSelect, v).Scan(&enough)
			if err != nil {
				if err == sql.ErrNoRows {
					warehouseErr(k, v, tx)
					return
				}
				warehouseErr(k, v, tx)
				return
			}
			if !enough {
				warehouseErr(k, v, tx)
				return
			}

			_, err = tx.ExecContext(ctx, q, v, WarehouseId)
			if err != nil {
				tx.Rollback()
				return
			}

			warehouseOk(k, v, tx)

			if err = tx.Commit(); err != nil {
				tx.Rollback()
			}
		}(k, v)
	}
	wg.Wait()
	return res, nil
}

func (db DatabaseConn) CancelReservationProduct(
	ctx context.Context,
	id []int64,
	WarehouseId int64,
) (map[string]string, error) {
	q := `
			UPDATE products 
			SET quantity = quantity + 1,
			    in_reserve = in_reserve - 1
			WHERE id = $1 AND quantity > 0 AND warehouse_id = $2
			`

	var wg sync.WaitGroup
	res := make(map[string]string)

	warehouseOk := func(k int, v int64, tx *sql.Tx) {
		mapK := "Товар номер " + strconv.Itoa(k+1) + " с id - " + strconv.Itoa(int(v))
		res[mapK] = " снят с резерва"
	}

	for k, v := range id {
		wg.Add(1)
		go func(k int, v int64) {
			defer wg.Done()
			tx, err := db.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
			if err != nil {
				errors.LogError(err)
				return
			}

			_, err = tx.ExecContext(ctx, q, v, WarehouseId)
			if err != nil {
				tx.Rollback()
				return
			}

			warehouseOk(k, v, tx)

			if err = tx.Commit(); err != nil {
				tx.Rollback()
			}
		}(k, v)
	}
	wg.Wait()
	return res, nil
}

func (db DatabaseConn) GetAllProducts(
	ctx context.Context,
	id int64,
) ([]*product_service.Product, error) {
	q := `
		SELECT id, name, size, quantity, in_reserve, warehouse_id
		FROM products
		WHERE warehouse_id = $1
			`

	products := make([]*product_service.Product, 0)

	rows, err := db.conn.Query(q, id)
	if err != nil {
		return nil, errors.HandleDatabaseError(err)
	}

	for rows.Next() {
		product := &product_service.Product{}

		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Size,
			&product.Quantity,
			&product.InReserve,
			&product.WarehouseId,
		)

		if err != nil {
			return nil, errors.HandleDatabaseError(err)
		}

		products = append(products, product)
	}

	err = rows.Err()

	if err != nil {
		return nil, errors.HandleDatabaseError(err)
	}

	return products, err
}
