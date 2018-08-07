package db

import (
	"github.com/dAdAbird/searchy"
)

type DB interface {
	Search(query string, limit, offset int) ([]*searchy.Site, error)
}
