package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/simondrake/copy-paste-notes/internal/sqlite"
	"github.com/spf13/cobra"
)

func newListCommand(client *sqlite.Client) *cobra.Command {
	var autoWrapText bool

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all notes",
		Run: func(_ *cobra.Command, _ []string) {
			notes, err := client.ListNotes()
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to list notes: ", err)
				os.Exit(1)
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Title", "Description"})

			table.SetAutoWrapText(autoWrapText)

			for _, n := range notes {
				table.Append([]string{fmt.Sprint(n.ID), n.Title, n.Description})
			}
			table.Render()
		},
	}

	listCmd.Flags().BoolVarP(&autoWrapText, "autowrap", "w", false, "whether to auto wrap the text output")

	return listCmd
}
