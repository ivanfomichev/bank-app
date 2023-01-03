package database

// Config - storage connection config
type Config struct {
	URI          string `mapstructure:"uri"`
	DatabaseName string `mapstructure:"database_name"`
}
