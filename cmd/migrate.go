package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func newMigrateCommand(file string) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Runs the DB migration scripts",
		Run: func(_ *cobra.Command, _ []string) {
			uri := fmt.Sprintf("sqlite3://%s", file)

			var m *migrate.Migrate

			if err := retry(func() error {
				var err error

				fmt.Println("trying migration...")

				m, err = migrate.New("file://internal/migrations", uri)
				if err != nil {
					fmt.Printf("migration failed: (%+v) \n", err)
					return err
				}

				fmt.Println("migration completed successfully")

				return nil
			}); err != nil {
				fmt.Fprintln(os.Stderr, "error creating new migration: ", err)
				os.Exit(1)
			}

			if err := m.Up(); err != nil {
				fmt.Fprintln(os.Stderr, "error running migration: ", err)
				os.Exit(1)
			}
		},
	}

	return migrateCmd
}

// Retry is an exponential backoff retry helper. It is used to wait for postgres to boot up
func retry(op func() error) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 10
	bo.MaxElapsedTime = time.Minute * 5

	if err := backoff.Retry(op, bo); err != nil {
		if bo.NextBackOff() == backoff.Stop {
			return fmt.Errorf("reached retry deadline")
		}

		return fmt.Errorf("retry failed: %w", err)
	}

	return nil
}
