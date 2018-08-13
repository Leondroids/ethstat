package app

import (
	"github.com/Leondroids/go-ethereum-rpc/rpc"
	"github.com/Leondroids/gox"
	"github.com/Leondroids/go-ethereum-db/etherdb"
	"fmt"
	"bytes"
	"log"
)

const (
	InfuraEndpoint     = "https://mainnet.infura.io/3l5dxBOP3wPspnRDdG1u"
	DefaultPort        = "8060"
	DefaultRPCEndpoint = rpc.GCloudEndpoint
	DefaultSQLHost     = "localhost"
	DefaultSQLUser     = "postgres"
	DefaultSQLPass     = "pass"
	DefaultSQLDB       = "ethstat"

	EnvPort        = "PORT"
	EnvRPCEndpoint = "RPC_ENDPOINT"

	EnvSQLHost = "SQL_HOST"
	EnvSQLDB   = "SQL_DB"
	EnvSQLUser = "SQL_USER"
	EnvSQLPass = "SQL_PASS"
)

type Config struct {
	Port        string
	RPCEndpoint string
	SQLInfo     *etherdb.SqlDBInfo
}

type Context struct {
	Config         *Config
	Client         *rpc.Client
	DB             *gox.SQLDB
	TableBlockDate *TableBlockDate
}

func InitApp() (*Context, error) {
	config := NewConfig()

	db, blockDateTable, err := InitDB(config.SQLInfo)

	if err != nil {
		return nil, err
	}

	return &Context{
		Config:         config,
		Client:         rpc.NewRPCClient(config.RPCEndpoint),
		DB:             db,
		TableBlockDate: blockDateTable,
	}, nil
}

func InitDB(info *etherdb.SqlDBInfo) (*gox.SQLDB, *TableBlockDate, error) {
	schema := "etherstats"

	blockDateTable := NewTableBlockDate(schema, TableBlockDateDefaultName)

	log.Println(blockDateTable.CreateTableStatement())

	db, _, err := gox.NewDatabaseBuilder(info.DBName).
		WithHost(info.Host).
		WithSchema(schema).
		WithUser(info.User).
		WithPassword(info.Password).
		AddTable(blockDateTable).
		OpenAndInitializeDB(false)

	return db, blockDateTable, err
}

func NewConfig() *Config {
	return &Config{
		Port:        fmt.Sprintf(":%v", gox.EnvReadStringOr(EnvPort, DefaultPort)),
		RPCEndpoint: gox.EnvReadStringOr(EnvRPCEndpoint, DefaultRPCEndpoint),
		SQLInfo:
		&etherdb.SqlDBInfo{
			Host:     gox.EnvReadStringOr(EnvSQLHost, DefaultSQLHost),
			Password: gox.EnvReadStringOr(EnvSQLPass, DefaultSQLPass),
			DBName:   gox.EnvReadStringOr(EnvSQLDB, DefaultSQLDB),
			User:     gox.EnvReadStringOr(EnvSQLUser, DefaultSQLUser),
		},
	}
}

func (c *Config) String() string {
	bb := &bytes.Buffer{}
	bb.WriteString("\n")
	bb.WriteString(gox.ConsoleWriteLabeledValue("RPCEndpoint", c.RPCEndpoint))
	bb.WriteString(gox.ConsoleWriteLabeledValue("Port", c.Port))
	bb.WriteString(gox.ConsoleWriteLabeledValue("PG Host", c.SQLInfo.Host))
	bb.WriteString(gox.ConsoleWriteLabeledValue("PG DBName", c.SQLInfo.DBName))
	bb.WriteString(gox.ConsoleWriteLabeledValue("PG User", c.SQLInfo.User))
	bb.WriteString(gox.ConsoleWriteLabeledValue("PG Pass", c.SQLInfo.Password))

	return bb.String()
}
