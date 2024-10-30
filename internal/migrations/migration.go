package migration

import "embed"

//go:embed *
var FS embed.FS

const (

	//MongoPath      = "mongo"
	PostgresPath = "postgres"
)
