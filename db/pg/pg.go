package pg

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/dAdAbird/searchy"
	"github.com/dAdAbird/searchy/db"
)

type Conn struct {
	db *sql.DB

	stmt *sql.Stmt
}

func New(source string) (ret db.DB, err error) {
	c := new(Conn)
	for i := 0; i < 5; i++ {
		c.db, err = sql.Open("postgres", source) //"postgres://dsp:dsp@localhost/dsp?sslmode=disable")
		if err != nil {
			return nil, err
		}
		if err == nil && c.db.Ping() == nil {
			break
		}
	}

	c.stmt, err = c.db.Prepare(`SELECT site_id, site_url, site_text FROM sites WHERE site_tokens @@ to_tsquery($1) limit $2 offset $3`)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Conn) Search(query string, limit, offset int) ([]*searchy.Site, error) {
	sites := make([]*searchy.Site, 0, limit)

	rows, err := c.stmt.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		site := new(searchy.Site)
		err = rows.Scan(&site.ID, &site.URL, &site.Text)
		if err != nil {
			continue
		}

		sites = append(sites, site)
	}

	return sites, err
}
