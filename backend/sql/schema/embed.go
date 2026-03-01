package schema

import "embed"

// Migrations contains all SQL migration files embedded in the binary.
// This allows running migrations without needing the SQL files on disk.
//
//go:embed *.sql
var Migrations embed.FS
