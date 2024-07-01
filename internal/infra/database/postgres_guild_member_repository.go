package database

import (
	"log/slog"

	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/jmoiron/sqlx"
)

const (
	createGuildMember = "create guild member"
	isAGuildMember    = "find guild member by guild_id, user_id and profile_id"
)

func guildMemberQueries() map[string]string {
	return map[string]string{
		createGuildMember: `INSERT INTO guild_members
		(profile_id, user_id, guild_id, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING *`,
		isAGuildMember: `SELECT * FROM guild_members
		WHERE guild_id = $1 AND user_id = $2 AND profile_id = $3 AND is_active = true`,
	}
}

func NewPostgresGuildMemberRepository(db *sqlx.DB) (*PostgresGuildMemberRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range guildMemberQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				NewPreparationErr(queryName, "guild member", err)
		}

		stmts[queryName] = stmt
	}

	return &PostgresGuildMemberRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

type PostgresGuildMemberRepository struct {
	db    *sqlx.DB
	stmts map[string]*sqlx.Stmt
}

func (r *PostgresGuildMemberRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			NewStatementNotPreparedErr(queryName, "guild member")
	}

	return stmt, nil
}

func (r *PostgresGuildMemberRepository) Store(m *entity.GuildMember) (int, error) {
	stmt, err := r.statement(createGuildMember)
	if err != nil {
		return -1, err
	}

	var result entity.GuildMember
	if err := stmt.Get(
		&result,
		m.ProfileID,
		m.UserID,
		m.GuildID,
		m.IsActive,
	); err != nil {
		slog.Error(err.Error())
		return -1, err
	}

	return result.ID, nil
}

func (r *PostgresGuildMemberRepository) IsAGuildMember(profileID, userID, guildID int) (*entity.GuildMember, error) {
	stmt, err := r.statement(isAGuildMember)
	if err != nil {
		return nil, err
	}

	var guildMember entity.GuildMember
	if err := stmt.Get(&guildMember, guildID, userID, profileID); err != nil {
		return nil, err
	}

	return &guildMember, nil
}
