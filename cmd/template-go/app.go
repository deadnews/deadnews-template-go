package main

import "github.com/jackc/pgx/v5/pgxpool"

// App holds application dependencies.
type App struct {
	DB *pgxpool.Pool
}

// NewApp creates a new App with the given configuration.
func NewApp(cfg *Config) (*App, error) {
	pool, err := openDB(cfg.DSN)
	if err != nil {
		return nil, err
	}
	return &App{DB: pool}, nil
}

// Close closes application resources.
func (app *App) Close() {
	if app.DB != nil {
		app.DB.Close()
	}
}
