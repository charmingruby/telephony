package database

import (
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/jmoiron/sqlx"
)

const (
	createUser      = "create user"
	findUserByID    = "find user by id"
	findUserByEmail = "find user by email"
)

func userQueries() map[string]string {
	return map[string]string{
		createUser: `INSERT INTO users
		(first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING *`,
		findUserByID: `SELECT * FROM users 
		WHERE id = $1`,
		findUserByEmail: `SELECT * FROM users 
		WHERE email = $1`,
	}
}

func NewPostgresUserRepository(db *sqlx.DB) (*PostgresUserRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range userQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				NewPreparationErr(queryName, "user", err)
		}

		stmts[queryName] = stmt
	}

	return &PostgresUserRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

type PostgresUserRepository struct {
	db    *sqlx.DB
	stmts map[string]*sqlx.Stmt
}

func (r *PostgresUserRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			NewStatementNotPreparedErr(queryName, "user")
	}

	return stmt, nil
}

func (r *PostgresUserRepository) Store(u *entity.User) (int, error) {
	stmt, err := r.statement(createUser)
	if err != nil {
		return -1, err
	}

	var result entity.User
	if err := stmt.Get(
		&result,
		u.FirstName,
		u.LastName,
		u.Email,
		u.PasswordHash,
	); err != nil {
		return -1, err
	}

	return result.ID, nil
}

func (r *PostgresUserRepository) FindByID(id int) (*entity.User, error) {
	stmt, err := r.statement(findUserByID)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := stmt.Get(&user, id); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (*entity.User, error) {
	stmt, err := r.statement(findUserByEmail)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := stmt.Get(&user, email); err != nil {
		return nil, err
	}

	return &user, nil
}
