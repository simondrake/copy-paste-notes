package cmd

import (
	"fmt"
	"os"

	"github.com/simondrake/copy-paste-notes/internal/sqlite"
	"github.com/spf13/cobra"
)

func newDeleteCommand(client *sqlite.Client) *cobra.Command {
	var id int

	addCmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a note by it's ID",
		Run: func(_ *cobra.Command, _ []string) {
			if err := client.DeleteNote(id); err != nil {
				fmt.Fprintln(os.Stderr, "unable to delete note: ", err)
				os.Exit(1)
			}
		},
	}

	addCmd.Flags().IntVar(&id, "id", 0, "id of the note")

	addCmd.MarkFlagRequired("id")

	return addCmd
}
