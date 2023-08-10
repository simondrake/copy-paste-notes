package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/simondrake/copy-paste-notes/internal/notes"
	"github.com/simondrake/copy-paste-notes/internal/sqlite"
	"github.com/spf13/cobra"
)

func newListCommand(client *sqlite.Client) *cobra.Command {
	var (
		autoWrapText bool
		raw          bool
		titleOnly    bool
		format       string
	)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all notes",
		Run: func(_ *cobra.Command, _ []string) {
			ns, err := client.ListNotes()
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to list notes: ", err)
				os.Exit(1)
			}

			switch format {
			case "table":
				outputTable(ns, titleOnly, autoWrapText, raw)
			case "json":
				e := json.NewEncoder(os.Stdout)
				e.SetIndent("", "    ")
				if err := e.Encode(ns); err != nil {
					fmt.Fprintln(os.Stderr, "Unable to encode json response: ", err)
					os.Exit(1)
				}
			default:
				fmt.Fprintln(os.Stderr, "unsupported format option")
				os.Exit(1)
			}
		},
	}

	listCmd.Flags().BoolVarP(&autoWrapText, "autowrap", "w", false, "whether to auto wrap the text output")
	listCmd.Flags().BoolVarP(&raw, "raw", "r", false, "Whether to show the raw text (e.g. don't parse the newline character as a literal newline)")
	listCmd.Flags().BoolVar(&titleOnly, "title-only", false, "Whether to only show the title")
	listCmd.Flags().StringVarP(&format, "format", "f", "table", "output format to use [table, json]")

	return listCmd
}

func outputTable(ns []notes.Note, titleOnly bool, autoWrapText bool, raw bool) {
	table := tablewriter.NewWriter(os.Stdout)

	if titleOnly {
		table.SetHeader([]string{"ID", "Create Timestamp", "Title"})
	} else {
		table.SetHeader([]string{"ID", "Create Timestamp", "Title", "Description"})
	}

	table.SetAutoWrapText(autoWrapText)

	for _, n := range ns {
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
}
