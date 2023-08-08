package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/simondrake/copy-paste-notes/internal/sqlite"
	"github.com/spf13/cobra"
)

func newListCommand(client *sqlite.Client) *cobra.Command {
	var (
		autoWrapText bool
		raw          bool
		titleOnly    bool
	)

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

			if titleOnly {
				table.SetHeader([]string{"ID", "Create Timestamp", "Title"})
			} else {
				table.SetHeader([]string{"ID", "Create Timestamp", "Title", "Description"})
			}

			table.SetAutoWrapText(autoWrapText)

			for _, n := range notes {
				if titleOnly {
					table.Append([]string{fmt.Sprint(n.ID), n.CreateTimestamp, n.Title})
					continue
				}

				if !raw {
					spl := strings.Split(n.Description, "\\n")

					out := make([]string, len(spl))
					for i, s := range spl {
						out[i] = strings.TrimSpace(s)
					}

					n.Description = strings.Join(out, "\n")
				}

				table.Append([]string{fmt.Sprint(n.ID), n.CreateTimestamp, n.Title, n.Description})
			}

			table.Render()
		},
	}

	listCmd.Flags().BoolVarP(&autoWrapText, "autowrap", "w", false, "whether to auto wrap the text output")
	listCmd.Flags().BoolVarP(&raw, "raw", "r", false, "Whether to show the raw text (e.g. don't parse the newline character as a literal newline)")
	listCmd.Flags().BoolVar(&titleOnly, "title-only", false, "Whether to only show the title")

	return listCmd
}
