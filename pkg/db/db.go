package db

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/TestardR/seller-payout/internal/model"
	"github.com/golang-migrate/migrate/v4"
	migrate_pg "github.com/golang-migrate/migrate/v4/database/postgres"

	// perform migrate init.
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// ErrRecordNotFound record not found error.
	ErrRecordNotFound = gorm.ErrRecordNotFound
	// ErrAbsolutePath error raised when a migrationDir is absolute.
	ErrAbsolutePath = errors.New("path should not be absolute")
)

//go:generate mockgen -source=db.go -destination=$MOCK_FOLDER/db.go -package=mock

// Config holds our database configuration.
type Config struct {
	User     string
	Name     string
	Password string
	Host     string
}

// DB represents the interface to interact with the database.
type DB interface {
	Health() error
	Begin() (DB, error)
	Rollback() error
	Commit() error

	Insert(dest interface{}) error
	Update(dest interface{}) error

	FindByID(dest interface{}, id string) error
	FindAll(dest interface{}) error
	FindAllWhere(dest interface{}, conds map[string]interface{}) error

	FindPayoutsBySellerID(string) ([]model.Payout, error)
	FindUnpaidOutItemsBySellerID(string) ([]model.Item, error)
	FindUnpaidOutItems() ([]model.Item, error)

	RunMigrations(path string) error
}

type database struct {
	driver *gorm.DB
	config Config
}

// New creates a new postgres database.
func New(c Config) (DB, error) {
	var err error

	var driver *gorm.DB

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.User, c.Password, c.Name)
	if driver, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}); err != nil {
		return &database{}, err
	}

	return &database{driver: driver, config: c}, nil
}

func (d database) Health() error {
	if pinger, ok := d.driver.ConnPool.(interface{ Ping() error }); ok {
		return pinger.Ping()
	}

	return errors.New("failed to cast ConnPool to pingable interface")
}

// Begin begins a sql transaction.
func (d database) Begin() (DB, error) {
	tx := d.driver.Begin()
	return &database{driver: tx}, tx.Error
}

// Commit commit a transaction.
func (d database) Commit() error {
	return d.driver.Commit().Error
}

// Rollback rollback all previous operations since the Begin() call.
func (d database) Rollback() error {
	return d.driver.Rollback().Error
}

// Insert take pointer to struct and add new entry in database and generate automatically a new ID
//  var u User
//  u.Name = "Bob"
//  Insert(&u)
func (d database) Insert(dest interface{}) error {
	return d.driver.Create(dest).Error
}

// FindAll retrieving all object in database
//   var users []User
//   db.FindAll(&users)
//   for _, u := range users {
//   	fmt.Println(u.ID)
//   }
func (d database) FindAll(dest interface{}) error {
	return d.driver.Find(dest).Error
}

// FindWhere take pointer to struct and conditions value
//  var u []User
//  FindWhere(&u, map[string]interface{}{"name": "Bob"})
//  fmt.Println(u[0].Name) // Bob
func (d database) FindAllWhere(dest interface{}, conds map[string]interface{}) error {
	return d.driver.Where(conds).Find(dest).Error
}

// FindByID take pointer to struct and the ID, then it fill the struct
//  var u User
//  FindByID(&u, 42)
//  fmt.Println(u.Name) // Bob
func (d database) FindByID(dest interface{}, id string) error {
	return d.driver.Take(dest, "id = ?", id).Error
}

// Update update the provided struct with it fields, the value Struct must have a ID
//  var u User
//  FindById(&u, 42)
//  u.Name = "Marvin"
//  Update(&u)
func (d database) Update(value interface{}) error {
	return d.driver.Save(value).Error
}

// Migrate applies all migrations needed.
func (d database) RunMigrations(migrationDir string) error {
	if strings.HasPrefix(migrationDir, "/") {
		return ErrAbsolutePath
	}

	db, err := d.driver.DB()
	if err != nil {
		return err
	}

	instanceDB, err := migrate_pg.WithInstance(db, &migrate_pg.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+filepath.Clean(migrationDir), d.config.Name, instanceDB)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
