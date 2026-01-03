package portfolio

import (
    "embed"
)

//go:embed sql/migrations/*.sql
var EmbedMigrations embed.FS