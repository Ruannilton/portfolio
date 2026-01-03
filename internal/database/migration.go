package database

import (
	"database/sql"
	"log"
	"portfolio"

	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")
	// Define que o sistema de arquivos base é o embed
	goose.SetBaseFS(portfolio.EmbedMigrations)

	// Configura o dialeto (postgres, mysql, etc)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// Roda o comando "up" na pasta onde estão os arquivos (dentro do embed)
	// Note que o caminho deve bater com a estrutura de pastas do embed
	if err := goose.Up(db, "sql/migrations"); err != nil {
		return err
	}

	log.Println("Migrações executadas com sucesso!")
	return nil
}
