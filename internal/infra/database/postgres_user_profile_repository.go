package database

import (
	"github.com/charmingruby/telephony/internal/domain/user/entity"
	"github.com/jmoiron/sqlx"
)

const (
	createUserProfile            = "create user profile"
	findUserProfileByUserID      = "find user profile by user id"
	findUserProfileByDisplayName = "find user profile by display name"
)

func userProfileQueries() map[string]string {
	return map[string]string{
		createUserProfile: `INSERT INTO users_profile
		(display_name, bio, guilds_quantity, messages_quantity, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *`,
		findUserProfileByUserID: `SELECT * FROM users_profile 
		WHERE user_id = $1`,
		findUserProfileByDisplayName: `SELECT * FROM users_profile 
		WHERE display_name = $1`,
	}
}

func NewPostgresUserProfileRepository(db *sqlx.DB) (*PostgresUserProfileRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range userProfileQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				NewPreparationErr(queryName, "user profile", err)
		}

		stmts[queryName] = stmt
	}

	return &PostgresUserProfileRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

type PostgresUserProfileRepository struct {
	db    *sqlx.DB
	stmts map[string]*sqlx.Stmt
}

func (r *PostgresUserProfileRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			NewStatementNotPreparedErr(queryName, "user profile")
	}

	return stmt, nil
}

func (r *PostgresUserProfileRepository) Store(p *entity.UserProfile) (int, error) {
	stmt, err := r.statement(createUserProfile)
	if err != nil {
		return -1, err
	}

	var result entity.UserProfile
	if err := stmt.Get(
		&result,
		p.DisplayName,
		p.Bio,
		p.GuildsQuantity,
		p.MessagesQuantity,
		p.UserID,
	); err != nil {
		return -1, err
	}

	return result.ID, nil
}

func (r *PostgresUserProfileRepository) FindByUserID(userID int) (*entity.UserProfile, error) {
	stmt, err := r.statement(findUserProfileByUserID)
	if err != nil {
		return nil, err
	}

	var profile entity.UserProfile
	if err := stmt.Get(&profile, userID); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *PostgresUserProfileRepository) FindByDisplayName(displayName string) (*entity.UserProfile, error) {
	stmt, err := r.statement(findUserProfileByDisplayName)
	if err != nil {
		return nil, err
	}

	var profile entity.UserProfile
	if err := stmt.Get(&profile, displayName); err != nil {
		return nil, err
	}

	return &profile, nil
}
