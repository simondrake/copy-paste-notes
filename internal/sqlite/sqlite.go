package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/simondrake/copy-paste-notes/internal/notes"
)

var ErrDeleteFailed = errors.New("delete failed")

type Client struct {
	db *sql.DB
}

func New(file string) (*Client, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS [notes] ( id INTEGER NOT NULL PRIMARY KEY, create_timestamp TEXT, title TEXT, description TEXT);"); err != nil {
		return nil, err
	}

	return &Client{
		db: db,
	}, nil
}

func (c *Client) ListNotes() ([]notes.Note, error) {
	rows, err := c.db.Query("SELECT id, create_timestamp, title, description FROM notes")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := make([]notes.Note, 0)
	for rows.Next() {
		n := notes.Note{}
		if err := rows.Scan(&n.ID, &n.Createtimestamp, &n.Title, &n.Description); err != nil {
			return nil, err
		}
		out = append(out, n)
	}

	return out, nil
}

func (c *Client) GetNote(id int) (*notes.Note, error) {
	row := c.db.QueryRow("SELECT id, create_timestamp, title, description FROM notes WHERE id=?", id)

	n := &notes.Note{}

	err := row.Scan(&n.ID, &n.Createtimestamp, &n.Title, &n.Description)

	return n, err
}

func (c *Client) InsertNote(n notes.Note) (int, error) {
	res, err := c.db.Exec("INSERT INTO notes VALUES(NULL,?,?,?);", n.CreateTimestamp, n.Title, n.Description)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (c *Client) UpdateNote(string, notes.Note) (*notes.Note, error) {
	return nil, errors.New("not implemented")
}

func (c *Client) DeleteNote(id int) error {
	res, err := c.db.Exec("DELETE FROM notes WHERE id=?", id)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return ErrDeleteFailed
	}

	return nil
}
