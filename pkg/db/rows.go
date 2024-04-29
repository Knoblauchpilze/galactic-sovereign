package db

import (
	"github.com/KnoblauchPilze/user-service/pkg/common"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/jackc/pgx/v5"
)

type Scannable interface {
	Scan(dest ...interface{}) error
}

type RowParser func(row Scannable) error

type rowsImpl struct {
	rows pgx.Rows
	next bool
	err  error
}

func newRows(rows pgx.Rows, err error) Rows {
	r := rowsImpl{
		rows: rows,
		err:  err,
	}

	if !common.IsInterfaceNil(r.rows) && r.err == nil {
		r.next = r.rows.Next()
	}

	return &r
}

func (r *rowsImpl) empty() bool {
	return r.rows == nil || !r.next
}

func (r *rowsImpl) Err() error {
	return r.err
}

func (r *rowsImpl) Close() {
	if r.rows != nil {
		r.rows.Close()
	}
}

func (r *rowsImpl) GetSingleValue(parser RowParser) error {
	if err := r.Err(); err != nil {
		return err
	}
	if r.empty() {
		return errors.NewCode(NoMatchingSqlRows)
	}

	defer r.Close()

	if err := parser(r.rows); err != nil {
		return err
	}

	r.next = r.rows.Next()
	if r.next {
		return errors.NewCode(MoreThanOneMatchingSqlRows)
	}

	return nil
}

func (r *rowsImpl) GetAll(parser RowParser) error {
	if err := r.Err(); err != nil {
		return err
	}

	defer r.Close()

	for r.next {
		if err := parser(r.rows); err != nil {
			return err
		}

		r.next = r.rows.Next()
	}

	return nil
}
