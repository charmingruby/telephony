package database

import (
	"log/slog"

	"github.com/charmingruby/telephony/internal/core"
	"github.com/charmingruby/telephony/internal/domain/guild/entity"
	"github.com/jmoiron/sqlx"
)

const (
	createGuild        = "create user"
	findGuildByID      = "find user by id"
	findGuildByName    = "find user by email"
	listAvailableGuild = "list paginated available guilds"
	deleteGuild        = "delete a guild"
)

func guildQueries() map[string]string {
	return map[string]string{
		createGuild: `INSERT INTO guilds
		(name, description, channels_quantity, owner_id)
		VALUES ($1, $2, $3, $4)
		RETURNING *`,
		findGuildByID: `SELECT * FROM guilds 
		WHERE id = $1`,
		findGuildByName: `SELECT * FROM guilds 
		WHERE name = $1`,
		listAvailableGuild: `SELECT * FROM guilds 
		WHERE deleted_at IS NULL
		LIMIT $1 
		OFFSET $2`,
		deleteGuild: `UPDATE guilds 
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL`,
	}
}

func NewPostgresGuildRepository(db *sqlx.DB) (*PostgresGuildRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range guildQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				NewPreparationErr(queryName, "guild", err)
		}

		stmts[queryName] = stmt
	}

	return &PostgresGuildRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

type PostgresGuildRepository struct {
	db    *sqlx.DB
	stmts map[string]*sqlx.Stmt
}

func (r *PostgresGuildRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			NewStatementNotPreparedErr(queryName, "guild")
	}

	return stmt, nil
}

func (r *PostgresGuildRepository) Store(g *entity.Guild) (int, error) {
	stmt, err := r.statement(createGuild)
	if err != nil {
		return -1, err
	}

	var result entity.Guild
	if err := stmt.Get(
		&result,
		g.Name,
		g.Description,
		g.ChannelsQuantity,
		g.OwnerID,
	); err != nil {
		slog.Error(err.Error())
		return -1, err
	}

	return 2, nil
}

func (r *PostgresGuildRepository) FindByID(id int) (*entity.Guild, error) {
	stmt, err := r.statement(findGuildByID)
	if err != nil {
		return nil, err
	}

	var guild entity.Guild
	if err := stmt.Get(&guild, id); err != nil {
		return nil, err
	}

	return &guild, nil
}

func (r *PostgresGuildRepository) FindByName(name string) (*entity.Guild, error) {
	stmt, err := r.statement(findGuildByName)
	if err != nil {
		return nil, err
	}

	var guild entity.Guild
	if err := stmt.Get(&guild, name); err != nil {
		return nil, err
	}

	return &guild, nil
}

func (r *PostgresGuildRepository) ListAvailables(page int) ([]entity.Guild, error) {
	stmt, err := r.statement(findUserByEmail)
	if err != nil {
		return nil, err
	}

	var guilds []entity.Guild

	offset := page * core.ItemsPerPage()
	limit := core.ItemsPerPage()

	if err := stmt.Select(&guilds, limit, offset); err != nil {
		return nil, err
	}

	return guilds, nil
}

func (r *PostgresGuildRepository) Delete(g *entity.Guild) error {
	stmt, err := r.statement(deleteGuild)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(g.DeletedAt, g.ID); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
