package database

import (
	"log/slog"

	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/jmoiron/sqlx"
)

const (
	createChannel         = "create channel"
	findChannelByName     = "find channel by name"
	findChannelByID       = "find channel by id"
	listChannelsByGuildID = "list channels by guild_id"
	listAllChannels       = "list all channels"
)

func channelQueries() map[string]string {
	return map[string]string{
		createChannel: `INSERT INTO channels
		(name, messages_quantity, guild_id, owner_id)
		VALUES ($1, $2, $3, $4)
		RETURNING *`,
		findChannelByName: `SELECT * FROM channels
		WHERE guild_id = $1 AND name = $2 AND deleted_at IS NULL`,
		findChannelByID: `SELECT * FROM channels
		WHERE guild_id = $1 AND id = $2 AND deleted_at IS NULL`,
		listChannelsByGuildID: `SELECT * FROM channels 
		WHERE guild_id = $1 AND deleted_at IS NULL`,
		listAllChannels: `SELECT * FROM channels 
		WHERE deleted_at IS NULL`,
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
		g.Name,
		g.MessagesQuantity,
		g.GuildID,
		g.OwnerID,
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
	if err := stmt.Get(&channel, guildID, name); err != nil {
		return nil, err
	}

	return &channel, nil
}

func (r *PostgresChannelRepository) FindByID(channelID, guildID int) (*entity.Channel, error) {
	stmt, err := r.statement(findChannelByID)
	if err != nil {
		return nil, err
	}

	var channel entity.Channel
	if err := stmt.Get(&channel, guildID, channelID); err != nil {
		return nil, err
	}

	return &channel, nil
}

func (r *PostgresChannelRepository) ListChannelsByGuildID(guildID int) ([]entity.Channel, error) {
	stmt, err := r.statement(listChannelsByGuildID)
	if err != nil {
		return nil, err
	}

	var channels []entity.Channel

	if err := stmt.Select(&channels, guildID); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return channels, nil
}

func (r *PostgresChannelRepository) ListAllChannels() ([]entity.Channel, error) {
	stmt, err := r.statement(listAllChannels)
	if err != nil {
		return nil, err
	}

	var channels []entity.Channel

	if err := stmt.Select(&channels); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return channels, nil
}
