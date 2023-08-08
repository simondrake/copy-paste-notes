package cmd

import (
	"fmt"
	"os"

	"github.com/simondrake/copy-paste-notes/internal/notes"
	"github.com/simondrake/copy-paste-notes/internal/sqlite"
	"github.com/spf13/cobra"
)

func newUpdateCommand(client *sqlite.Client) *cobra.Command {
	var (
		id          int
		title       string
		description string
	)

	addCmd := &cobra.Command{
		Use:   "update",
		Short: "Updates a note",
		Run: func(_ *cobra.Command, _ []string) {
			_, err := client.UpdateNote(id, notes.Note{Title: title, Description: description})
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to update note: ", err)
				os.Exit(1)
			}
		},
	}

	addCmd.Flags().IntVar(&id, "id", 0, "id of the note")
	addCmd.Flags().StringVarP(&title, "title", "t", "", "title of the note")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "description of the note")

	addCmd.MarkFlagRequired("id")
	addCmd.MarkFlagsOneRequired("title", "description")

	return addCmd
}
