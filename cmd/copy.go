package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.design/x/clipboard"

	"github.com/spf13/cobra"

	"github.com/simondrake/copy-paste-notes/internal/notes"
	"github.com/simondrake/copy-paste-notes/internal/sqlite"
)

func newCopyCommand(client *sqlite.Client) *cobra.Command {
	var (
		id    int
		title string
		raw   bool
	)

	addCmd := &cobra.Command{
		Use:   "copy",
		Short: "Copies a note into the system clipboard",
		Run: func(_ *cobra.Command, _ []string) {
			var note *notes.Note

			if id != 0 {
				var err error
				note, err = client.GetNoteByID(id)
				if err != nil {
					fmt.Fprintln(os.Stderr, "unable to get note: ", err)
					os.Exit(1)
				}
			} else {
				var err error
				note, err = client.GetNoteByTitle(title)
				if err != nil {
					fmt.Fprintln(os.Stderr, "unable to get note: ", err)
					os.Exit(1)
				}
			}

			if err := clipboard.Init(); err != nil {
				fmt.Fprintln(os.Stderr, "unable to initialise clipboard: ", err)
				os.Exit(1)
			}

			if !raw {
				spl := strings.Split(note.Description, "\\n")

				out := make([]string, len(spl))
				for i, s := range spl {
					out[i] = strings.TrimSpace(s)
				}

				note.Description = strings.Join(out, "\n")
			}

			if os.Getenv("WAYLAND_DISPLAY") != "" {
				if err := exec.Command("wl-copy", note.Description).Run(); err != nil {
					fmt.Fprintln(os.Stderr, "unable to copy with wl-clipboard: ", err)
					os.Exit(1)
				}

				return
			}

			clipboard.Write(clipboard.FmtText, []byte(note.Description))
			// TODO - for some reason this is needed on linux. Find a way to avoid this cruft.
			time.Sleep(500 * time.Millisecond)
		},
	}

	addCmd.Flags().IntVar(&id, "id", 0, "id of the note")
	addCmd.Flags().StringVar(&title, "title", "", "title of the note")
	addCmd.Flags().BoolVarP(&raw, "raw", "r", false, "Whether to copy the raw text (e.g. don't parse the newline character as a literal newline)")

	addCmd.MarkFlagsOneRequired("id", "title")
	addCmd.MarkFlagsMutuallyExclusive("id", "title")

	return addCmd
}
