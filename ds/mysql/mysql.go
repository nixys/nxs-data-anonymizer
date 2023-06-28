package mysql

import (
	"fmt"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	client *gorm.DB
}

type Settings struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func Connect(s Settings) (MySQL, error) {

	client, err := gorm.Open(gmysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		s.User,
		s.Password,
		s.Host,
		s.Port,
		s.Database)), &gorm.Config{})
	if err != nil {
		return MySQL{}, err
	}

	return MySQL{
		client: client,
	}, nil
}

func (m *MySQL) Close() error {
	db, _ := m.client.DB()
	return db.Close()
}

func (m *MySQL) DBCleanup() error {

	tables, err := m.client.Migrator().GetTables()
	if err != nil {
		return fmt.Errorf("drop tables, get tables: %w", err)
	}

	for _, t := range tables {
		if err := m.client.Migrator().DropTable(t); err != nil {
			return fmt.Errorf("drop tables, table `%s`: %w", t, err)
		}
	}

	return nil
}
