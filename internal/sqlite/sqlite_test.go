//go:build integration

package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simondrake/copy-paste-notes/internal/notes"
)

var client *Client

func removeDBFile(dbName string) error {
	// Current working directory
	p, err := os.Getwd()
	if err != nil {
		return err
	}

	p = path.Join(p, dbName)

	if _, err := os.Open(p); err == nil {
		// No errors, so file exists. Now remove it.
		if err := os.Remove(p); err != nil {
			return err
		}
	}

	return nil
}

func TestMain(m *testing.M) {
	dbName := "./testdb"

	// Because we're using sqlite, we need to make sure we don't have an existing file
	if err := removeDBFile(dbName); err != nil {
		log.Fatalf("Could not remove existing db file: %+v", err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %+v", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "keinos/sqlite3",
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %+v", err)
	}

	resource.Expire(240) // Tell docker to hard kill the container in 240 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err := pool.Retry(func() error {
		client, err = New(dbName)
		if err != nil {
			return err
		}

		return client.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		fmt.Printf("Could not purge resource: %s \n", err)
	}

	os.Exit(code)
}

func TestListNotes(t *testing.T) {
	t.Run("should return a slice of zero length", func(t *testing.T) {
		notes, err := client.ListNotes()
		assert.NoError(t, err)
		assert.Zero(t, len(notes))
	})

	note := notes.Note{
		Title:           "test-list-title",
		Description:     "test-list-description",
		CreateTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	var rid int

	t.Run("should insert note without error", func(t *testing.T) {
		var err error

		rid, err = client.InsertNote(note)
		assert.NoError(t, err)
		assert.NotZero(t, rid)
	})

	t.Run("should exist when ListNotes is called", func(t *testing.T) {
		ns, err := client.ListNotes()
		assert.NoError(t, err)
		assert.Len(t, ns, 1)

		assert.Equal(t, rid, ns[0].ID)
		assert.Equal(t, note.Title, ns[0].Title)
		assert.Equal(t, note.Description, ns[0].Description)
		assert.Equal(t, note.CreateTimestamp, ns[0].CreateTimestamp)
	})

	t.Run("should delete the note successfully", func(t *testing.T) {
		require.NoError(t, client.DeleteNote(rid))
	})
}

func TestGetNoteByID(t *testing.T) {
	t.Run("should return an error when id does not exist", func(t *testing.T) {
		n, err := client.GetNoteByID(9009)
		assert.Nil(t, n)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	note := notes.Note{
		Title:           "test-getbyID-title",
		Description:     "test-getbyID-description",
		CreateTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	var rid int

	t.Run("should insert note without error", func(t *testing.T) {
		var err error

		rid, err = client.InsertNote(note)
		assert.NoError(t, err)
		assert.NotZero(t, rid)
	})

	t.Run("should be able to get the note by id", func(t *testing.T) {
		n, err := client.GetNoteByID(rid)
		assert.NoError(t, err)

		assert.Equal(t, rid, n.ID)
		assert.Equal(t, note.Title, n.Title)
		assert.Equal(t, note.Description, n.Description)
		assert.Equal(t, note.CreateTimestamp, n.CreateTimestamp)
	})

	t.Run("should delete the note successfully", func(t *testing.T) {
		require.NoError(t, client.DeleteNote(rid))
	})
}

func TestGetNoteByTitle(t *testing.T) {
	t.Run("should return an error when title does not exist", func(t *testing.T) {
		n, err := client.GetNoteByTitle("does-not-exist")
		assert.Nil(t, n)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	note := notes.Note{
		Title:           "test-getbyTitle-title",
		Description:     "test-getbyTitle-description",
		CreateTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	var rid int

	t.Run("should insert note without error", func(t *testing.T) {
		var err error

		rid, err = client.InsertNote(note)
		assert.NoError(t, err)
		assert.NotZero(t, rid)
	})

	t.Run("should be able to get the note by title", func(t *testing.T) {
		n, err := client.GetNoteByTitle(note.Title)
		assert.NoError(t, err)

		assert.Equal(t, rid, n.ID)
		assert.Equal(t, note.Title, n.Title)
		assert.Equal(t, note.Description, n.Description)
		assert.Equal(t, note.CreateTimestamp, n.CreateTimestamp)
	})

	t.Run("should delete the note successfully", func(t *testing.T) {
		require.NoError(t, client.DeleteNote(rid))
	})
}

func TestInsertNote(t *testing.T) {
	note := notes.Note{
		Title:           "test-insert-title",
		Description:     "test-insert-description",
		CreateTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	var rid int

	t.Run("should insert note without error", func(t *testing.T) {
		var err error

		rid, err = client.InsertNote(note)
		assert.NoError(t, err)
		assert.NotZero(t, rid)

		ns, err := client.ListNotes()
		assert.NoError(t, err)
		assert.Len(t, ns, 1)
	})

	t.Run("should delete the note successfully", func(t *testing.T) {
		require.NoError(t, client.DeleteNote(rid))
	})
}

func TestUpdateNote(t *testing.T) {
	note := notes.Note{
		Title:           "test-update-title",
		Description:     "test-update-description",
		CreateTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	var rid int

	t.Run("should insert note without error", func(t *testing.T) {
		var err error

		rid, err = client.InsertNote(note)
		assert.NoError(t, err)
		assert.NotZero(t, rid)

		ns, err := client.ListNotes()
		assert.NoError(t, err)
		assert.Len(t, ns, 1)
	})

	t.Run("should update title and description", func(t *testing.T) {
		nt := "totally-different-title"
		nd := "totally-different-description"

		ra, err := client.UpdateNote(rid, notes.Note{Title: nt, Description: nd})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), ra)

		n, err := client.GetNoteByID(rid)
		assert.NoError(t, err)
		assert.Equal(t, nt, n.Title)
		assert.Equal(t, nd, n.Description)
	})

	t.Run("should update only title", func(t *testing.T) {
		nt := "something-else"

		ra, err := client.UpdateNote(rid, notes.Note{Title: nt})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), ra)

		n, err := client.GetNoteByID(rid)
		assert.NoError(t, err)
		assert.Equal(t, nt, n.Title)
		assert.Equal(t, "totally-different-description", n.Description)
	})

	t.Run("should update only description", func(t *testing.T) {
		nd := "and-another-different-thing"

		ra, err := client.UpdateNote(rid, notes.Note{Description: nd})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), ra)

		n, err := client.GetNoteByID(rid)
		assert.NoError(t, err)
		assert.Equal(t, "something-else", n.Title)
		assert.Equal(t, nd, n.Description)
	})

	t.Run("should delete the note successfully", func(t *testing.T) {
		require.NoError(t, client.DeleteNote(rid))
	})
}

func TestDeleteNote(t *testing.T) {
	t.Run("should return an error when id does not exist", func(t *testing.T) {
		err := client.DeleteNote(9009)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrDeleteFailed)
	})

	note := notes.Note{
		Title:           "test-delete-title",
		Description:     "test-delete-description",
		CreateTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	var rid int

	t.Run("should insert note without error", func(t *testing.T) {
		var err error

		rid, err = client.InsertNote(note)
		assert.NoError(t, err)
		assert.NotZero(t, rid)
	})

	t.Run("should delete the note successfully", func(t *testing.T) {
		assert.NoError(t, client.DeleteNote(rid))

		ns, err := client.ListNotes()
		assert.NoError(t, err)
		assert.Zero(t, len(ns))
	})
}

func TestAppendStatement(t *testing.T) {
	stmt := "UPDATE notes SET"

	stmt = appendStatement(stmt, "title")
	assert.Equal(t, "UPDATE notes SET title = ?", stmt)

	stmt = appendStatement(stmt, "description")
	assert.Equal(t, "UPDATE notes SET title = ?, description = ?", stmt)

	stmt = appendStatement(stmt, "some_random_field")
	assert.Equal(t, "UPDATE notes SET title = ?, description = ?, some_random_field = ?", stmt)
}
