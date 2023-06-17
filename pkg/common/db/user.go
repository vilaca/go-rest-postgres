package db

// TODO hash and salt password
import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/pauljamescleary/gomin/pkg/common/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUser(id string) (*models.User, error)
	UserExists(name string) (bool, error)
	Login(name string, password string) (bool, error)
	CreateSession(*models.Session) error
}

type PostgresUserRepository struct {
	db *Database
}

func (repo PostgresUserRepository) CreateUser(user *models.User) (*models.User, error) {
	sql := `
	INSERT INTO users (id, name, password, enabled)
	VALUES ($1, $2, $3, $4)
	`
	_, err := repo.db.Conn.Exec(context.Background(), sql, user.ID, user.Name, user.Password, user.Enabled)
	if err != nil {
		panic(err)
	}
	return user, nil
}

func (repo PostgresUserRepository) CreateSession(session *models.Session) error {
	sql := `
	INSERT INTO sessions (id, username, started, ends)
	VALUES ($1, $2, $3, $4)
	`
	_, err := repo.db.Conn.Exec(context.Background(), sql, session.Id, session.UserName, session.Started, session.Ends)
	if err != nil {
		panic(err)
	}
	return nil
}

func (repo PostgresUserRepository) GetUser(id string) (*models.User, error) {
	sql := `
	SELECT id, name
	FROM users
	WHERE id = $1
	`
	var user models.User
	rows, err := repo.db.Conn.Query(context.Background(), sql, id)
	if err != nil {
		panic(err)
	}

	if err := pgxscan.ScanOne(&user, rows); err != nil {
		panic(err)
	}

	return &user, nil
}

func (repo PostgresUserRepository) UserExists(name string) (bool, error) {
	sql := `
	SELECT name
	FROM users
	WHERE name = $1
	`
	rows, err := repo.db.Conn.Query(context.Background(), sql, name)
	return rows.Next(), err
}

func (repo PostgresUserRepository) Login(name string, password string) (bool, error) {
	sql := `
	SELECT name, password
	FROM users
	WHERE name = $1 AND password = $2 AND enabled = true
	`
	rows, err := repo.db.Conn.Query(context.Background(), sql, name, password)
	return rows.Next(), err
}

func NewUserRepository(db *Database) (*PostgresUserRepository, error) {
	return &PostgresUserRepository{db: db}, nil
}
