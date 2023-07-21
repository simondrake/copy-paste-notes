package cmd

import (
	"fmt"
	"os"

	"github.com/simondrake/copy-paste-notes/internal/notes"
	"github.com/simondrake/copy-paste-notes/internal/sqlite"
	"github.com/spf13/cobra"
)

func newAddCommand(client *sqlite.Client) *cobra.Command {
	var (
		title       string
		description string
	)

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Adds a note",
		Run: func(_ *cobra.Command, _ []string) {
			_, err := client.InsertNote(notes.Note{
				Title:       title,
				Description: description,
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to insert note: ", err)
				os.Exit(1)
			}
		},
	}

	addCmd.Flags().StringVarP(&title, "title", "t", "", "title of the note")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "description of the note")

	return addCmd
}
