package main

import (
	"github.com/abdullah0iq/gator/internal/config"
	"github.com/abdullah0iq/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}
