package storage

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresStorageConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Port     uint32 `mapstructure:"port"`
}

type PostgresStorage struct {
	db     *sql.DB
	config PostgresStorageConfig
}

func (prs *PostgresStorage) Add(ruleType string, cidr net.IPNet) error {

	name := cidr.String()

	_, err := prs.db.Exec(
		"INSERT INTO rules (name,subnet,type) VALUES ($1, $2, $3);",
		name,
		cidr.String(),
		ruleType,
	)

	if err != nil {

		if errPg, ok := err.(*pq.Error); ok {

			if errPg.Code == "23505" {

				return RuleAlreadyExistsDBError
			}
		}

		return err
	}

	return nil
}

func (prs *PostgresStorage) Delete(cidr net.IPNet) error {

	res, err := prs.db.Exec("DELETE FROM rules WHERE subnet = $1;", cidr.String())

	if err != nil {

		return err
	}

	if c, _ := res.RowsAffected(); c == 0 {

		return RuleNotFoundDBError
	}

	return err
}

func (prs *PostgresStorage) GetList(ruleType string) ([]net.IPNet, error) {

	result := make([]net.IPNet, 0)

	rows, err := prs.db.Query("SELECT subnet FROM rules WHERE type=$1;", ruleType)

	if err != nil {

		return nil, err
	}

	buf := &net.IPNet{}

	for rows.Next() {

		if err = rows.Scan(buf); err != nil {

			continue
		}

		result = append(result, *buf)
	}

	return result, nil
}

func (prs *PostgresStorage) Connect() error {

	var err error

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		prs.config.User,
		prs.config.Password,
		prs.config.Host,
		prs.config.Port,
		prs.config.Database,
	)

	prs.db, err = sql.Open("postgres", dsn)

	if err != nil {

		return err
	}

	return prs.db.Ping()
}

func (prs *PostgresStorage) Close() error {

	return prs.db.Close()
}

func NewPostgresRuleStorage(config PostgresStorageConfig) *PostgresStorage {

	return &PostgresStorage{
		config: config,
	}
}
