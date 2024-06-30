package database

import (
	"log/slog"

	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/jmoiron/sqlx"
)

const (
	createChannel     = "create channel"
	findChannelByName = "find channel by name"
)

func channelQueries() map[string]string {
	return map[string]string{
		createChannel: `INSERT INTO channels
		(name, description, channels_quantity, owner_id)
		VALUES ($1, $2, $3, $4)
		RETURNING *`,
		findChannelByName: `SELECT * FROM channels 
		WHERE name = $1`,
	}
}

func NewPostgresChannelRepository(db *sqlx.DB) (*PostgresChannelRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range channelQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				NewPreparationErr(queryName, "channel", err)
		}

		stmts[queryName] = stmt
	}

	return &PostgresChannelRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

type PostgresChannelRepository struct {
	db    *sqlx.DB
	stmts map[string]*sqlx.Stmt
}

func (r *PostgresChannelRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			NewStatementNotPreparedErr(queryName, "channel")
	}

	return stmt, nil
}

func (r *PostgresChannelRepository) Store(g *entity.Channel) (int, error) {
	stmt, err := r.statement(createChannel)
	if err != nil {
		return -1, err
	}

	var result entity.Channel
	if err := stmt.Get(
		&result,
	); err != nil {
		slog.Error(err.Error())
		return -1, err
	}

	return result.ID, nil
}

func (r *PostgresChannelRepository) FindByName(guildID int, name string) (*entity.Channel, error) {
	stmt, err := r.statement(findChannelByName)
	if err != nil {
		return nil, err
	}

	var channel entity.Channel
	if err := stmt.Get(&channel, name); err != nil {
		return nil, err
	}

	return &channel, nil
}
