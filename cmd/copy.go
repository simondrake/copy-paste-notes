package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.design/x/clipboard"

	"github.com/simondrake/copy-paste-notes/internal/sqlite"
	"github.com/spf13/cobra"
)

func newCopyCommand(client *sqlite.Client) *cobra.Command {
	var (
		id           int
		parseNewline bool
	)

	addCmd := &cobra.Command{
		Use:   "copy",
		Short: "Copies a note into the system clipboard",
		Run: func(_ *cobra.Command, _ []string) {
			if id == 0 {
				fmt.Fprintln(os.Stderr, errors.New("id must be specified"))
				os.Exit(1)
			}

			note, err := client.GetNote(id)
			if err != nil {
				fmt.Fprintln(os.Stderr, "unable to get note: ", err)
				os.Exit(1)
			}

			if err := clipboard.Init(); err != nil {
				fmt.Fprintln(os.Stderr, "unable to initialise clipboard: ", err)
				os.Exit(1)
			}

			if parseNewline {
				note.Description = strings.ReplaceAll(note.Description, "\\n ", "\n")
			}

			if os.Getenv("WAYLAND_DISPLAY") != "" {
				if err := exec.Command("wl-copy", note.Description).Run(); err != nil {
					fmt.Fprintln(os.Stderr, "unable to copy with wl-clipboard: ", err)
					os.Exit(1)
				}

				return
			}

			clipboard.Write(clipboard.FmtText, []byte(note.Description))
		},
	}

	addCmd.Flags().IntVar(&id, "id", 0, "id of the note")
	addCmd.Flags().BoolVarP(&parseNewline, "parsenewline", "p", false, "Whether to parse the newline character as a literal newline")

	return addCmd
}
