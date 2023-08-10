package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/simondrake/copy-paste-notes/internal/notes"
	"github.com/simondrake/copy-paste-notes/internal/sqlite"
	"github.com/spf13/cobra"
)

func newGetCommand(client *sqlite.Client) *cobra.Command {
	var (
		id     int
		title  string
		format string
	)

	addCmd := &cobra.Command{
		Use:   "get",
		Short: "Gets a note in JSON format",
		Run: func(_ *cobra.Command, _ []string) {
			var n *notes.Note

			if title != "" {
				var err error
				n, err = client.GetNoteByTitle(title)
				if err != nil {
					fmt.Fprintln(os.Stderr, "unable to get note: ", err)
					os.Exit(1)
				}
			} else {
				// If title isn't defined then id must be
				var err error
				n, err = client.GetNoteByID(id)
				if err != nil {
					fmt.Fprintln(os.Stderr, "unable to get note: ", err)
					os.Exit(1)
				}
			}

			switch format {
			case "table":
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"ID", "Create Timestamp", "Title", "Description"})
				table.Append([]string{fmt.Sprint(n.ID), n.CreateTimestamp, n.Title, n.Description})

				table.Render()
			case "json":
				e := json.NewEncoder(os.Stdout)
				e.SetIndent("", "    ")
				if err := e.Encode(n); err != nil {
					fmt.Fprintln(os.Stderr, "Unable to encode json response: ", err)
					os.Exit(1)
				}
			default:
				fmt.Fprintln(os.Stderr, "unsupported format option")
				os.Exit(1)
			}
		},
	}

	addCmd.Flags().IntVar(&id, "id", 0, "id of the note")
	addCmd.Flags().StringVarP(&title, "title", "t", "", "title of the note")
	addCmd.Flags().StringVarP(&format, "format", "f", "table", "output format to use [table, json]")

	addCmd.MarkFlagsOneRequired("id", "title")
	addCmd.MarkFlagsMutuallyExclusive("id", "title")

	return addCmd
}
