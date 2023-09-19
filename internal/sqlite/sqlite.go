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

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS [notes] ( id INTEGER NOT NULL PRIMARY KEY, create_timestamp TEXT, title TEXT NOT NULL UNIQUE, description TEXT NOT NULL);"); err != nil {
		return nil, err
	}

	return &Client{
		db: db,
	}, nil
}

func (c *Client) Ping() error {
	return c.db.Ping()
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
		if err := rows.Scan(&n.ID, &n.CreateTimestamp, &n.Title, &n.Description); err != nil {
			return nil, err
		}
		out = append(out, n)
	}

	return out, nil
}

func (c *Client) GetNoteByID(id int) (*notes.Note, error) {
	row := c.db.QueryRow("SELECT id, create_timestamp, title, description FROM notes WHERE id=?", id)

	n := &notes.Note{}

	if err := row.Scan(&n.ID, &n.CreateTimestamp, &n.Title, &n.Description); err != nil {
		return nil, err
	}

	return n, nil
}

func (c *Client) GetNoteByTitle(title string) (*notes.Note, error) {
	row := c.db.QueryRow("SELECT id, create_timestamp, title, description FROM notes WHERE title=?", title)

	n := &notes.Note{}

	if err := row.Scan(&n.ID, &n.CreateTimestamp, &n.Title, &n.Description); err != nil {
		return nil, err
	}

	return n, nil
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

func (c *Client) UpdateNote(id int, note notes.Note) (int64, error) {
	var res sql.Result

	switch {
	case note.Title != "" && note.Description != "":
		stmt, err := c.db.Prepare("UPDATE notes SET title = ?, description = ? WHERE id = ?")
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		res, err = stmt.Exec(note.Title, note.Description, id)
		if err != nil {
			return 0, err
		}
	case note.Title != "":
		stmt, err := c.db.Prepare("UPDATE notes SET title = ? WHERE id = ?")
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		res, err = stmt.Exec(note.Title, id)
		if err != nil {
			return 0, err
		}
	case note.Description != "":
		stmt, err := c.db.Prepare("UPDATE notes SET description = ? WHERE id = ?")
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		res, err = stmt.Exec(note.Description, id)
		if err != nil {
			return 0, err
		}
	default:
		return 0, errors.New("either title or description must be defined")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
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
