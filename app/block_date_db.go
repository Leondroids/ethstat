package app

import (
	"github.com/Leondroids/gox"
	"database/sql"
	"fmt"
)

type BlockDate struct {
	Block int64
	Date  int64
}

func (blockdate *BlockDate) GetValues() []interface{} {
	values := make([]interface{}, 2)

	values[0] = blockdate.Block
	values[1] = blockdate.Date

	return values
}

/////////////////////////////

const (
	TableBlockDateDefaultName = "blockdate"
)

type TableBlockDate struct {
	Name       string
	Scheme     string
	Definition *gox.SQLTableDefinition
}

func NewTableBlockDate(scheme string, name string) *TableBlockDate {
	definition := gox.NewSQLTableBuilder(fmt.Sprintf("%v.%v", scheme, name)).
		WithBigIntColumn("block", gox.NotNull, gox.IsPrimary).
		WithBigIntColumn("date").
		Build()

	return &TableBlockDate{
		Name:       name,
		Scheme:     scheme,
		Definition: definition,
	}
}

// Must implement as SQLTable

func (it *TableBlockDate) DataModelTag() []string {
	return []string{"block", "date"}
}

func (it *TableBlockDate) TableName() string {
	return fmt.Sprintf("%v.%v", it.Scheme, it.Name)
}

func (it *TableBlockDate) KeyTag() string {
	return "block"
}

func (it *TableBlockDate) CreateTableStatement() string {
	return fmt.Sprintf(it.Definition.CreateStatement())
}

// Methods

func (it *TableBlockDate) Max(db *gox.SQLDB) (int64, error) {
	return db.Max(it, "block")
}

func (it *TableBlockDate) Put(blockDate *BlockDate, db *gox.SQLDB) error {
	res, err := it.Get(blockDate.Block, db)

	if err != nil {
		return err
	}

	// insert
	if res == nil {
		_, err = db.Insert(it, blockDate.GetValues())
		return err
	}

	// update
	updateres, err := db.DB.Exec(gox.CreateUpdateStatement(it), blockDate.Block, blockDate.Date)
	if err != nil {
		return err
	}

	count, err := updateres.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return fmt.Errorf("didnt update %v", blockDate)
	}

	return nil
}

func (tkv *TableBlockDate) Get(block int64, db *gox.SQLDB) (*BlockDate, error) {

	statement := fmt.Sprintf("SELECT * FROM %v WHERE block=$1", tkv.TableName())

	rows, err := db.DB.Query(statement, block)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		return parseBlockDateRow(rows)
	}

	return nil, nil
}

func (it *TableBlockDate) ListAll(db *gox.SQLDB) ([]BlockDate, error) {
	statement := fmt.Sprintf("select * from %v", it.TableName())
	return it.List(statement, db, nil)
}

func (it *TableBlockDate) List(statement string, db *gox.SQLDB, params []interface{}) ([]BlockDate, error) {
	rows, err := db.DB.Query(statement, params...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]BlockDate, 0)

	for rows.Next() {
		r, err := parseBlockDateRow(rows)

		if err != nil {
			return nil, err
		}

		result = append(result, *r)
	}

	return result, nil
}

func (tkv *TableBlockDate) Delete(key string, db *gox.SQLDB) error {
	return db.Delete(key, tkv)
}

func parseBlockDateRow(rows *sql.Rows) (*BlockDate, error) {
	row := &BlockDate{}
	err := rows.Scan(&row.Block, &row.Date)
	if err != nil {
		return nil, err
	}
	return row, nil
}
