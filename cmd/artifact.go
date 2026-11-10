package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	"github.com/FacileStudio/Mycelium/internal/browser"
	"github.com/FacileStudio/Mycelium/internal/config"
	"github.com/FacileStudio/Mycelium/internal/ui"
	"github.com/spf13/cobra"
)

var artifactCmd = &cobra.Command{
	Use:     "artifact",
	Aliases: []string{"artifacts", "report", "reports"},
	Short:   "Manage markdown artifacts and reports that travel between your machines",
}

var artifactAddCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Record a markdown artifact and open it",
	Long: "Record a markdown artifact and open it.\n\n" +
		"The document is copied into " + artifacts.Dir(config.DataDir()) + " and synced across machines, " +
		"with YAML frontmatter above the body. Artifacts are markdown only: HTML files are refused.\n\n" +
		"The identifier comes from the title with a date prefix to prevent collisions:\n\n" +
		"Artifacts are permanent: nothing deletes one except mycelium artifact rm.\n\n" +
		"  mycelium artifact add spec.md\n" +
		"  mycelium artifact add plan.md --title 'Migration plan'\n" +
		"  mycelium artifact add --title 'Quick note' --body-stdin < report.md",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		noOpen, _ := cmd.Flags().GetBool("no-open")
		bodyStdin, _ := cmd.Flags().GetBool("body-stdin")
		var art artifacts.Artifact
		var err error
		fromStdin := bodyStdin || (len(args) == 1 && args[0] == "-")
		if fromStdin {
			raw, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			art, err = artifacts.AddContent(config.DataDir(), raw, "stdin.md", artifacts.Request{
				Title:   title,
				Machine: config.MachineName(),
			}, time.Now())
			if err != nil {
				return err
			}
		} else {
			if len(args) == 0 {
				return errors.New("specify a file to record or pass --body-stdin")
			}
			art, err = artifacts.Add(config.DataDir(), artifacts.Request{
				Source:  args[0],
				Title:   title,
				Machine: config.MachineName(),
			}, time.Now())
			if err != nil {
				return err
			}
			warnExternalRefs(args[0])
		}
		artifactRecorded(art, noOpen)
		return nil
	},
}

var artifactListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the artifacts on this machine, newest first",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := artifacts.List(config.DataDir())
		if err != nil {
			return err
		}
		jsonOut, _ := cmd.Flags().GetBool("json")
		if jsonOut {
			return printJSON(all)
		}
		if len(all) == 0 {
			fmt.Println("No artifacts")
			return nil
		}
		width := artifactIDWidth(all)
		for _, art := range all {
			fmt.Println(artifactLine(art, width, time.Now()))
		}
		return nil
	},
}

var artifactOpenCmd = &cobra.Command{
	Use:   "open <id>",
	Short: "Open an artifact in this machine's browser",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		art, err := artifacts.Find(config.DataDir(), args[0])
		if err != nil {
			return err
		}
		target := artifactTarget(art)
		fmt.Println(target)
		showInBrowser(target)
		return nil
	},
}

var artifactRmCmd = &cobra.Command{
	Use:   "rm <id>",
	Short: "Delete an artifact",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := artifacts.Remove(config.DataDir(), args[0]); err != nil {
			return err
		}
		ui.Success("Deleted %s", artifacts.Slug(args[0]))
		return nil
	},
}

func warnExternalRefs(source string) {
	raw, err := os.ReadFile(source)
	if err != nil {
		return
	}
	refs := artifacts.ExternalRefs(raw)
	if len(refs) == 0 {
		return
	}
	ui.Warn("%d reference(s) will not resolve from disk: %s", len(refs), strings.Join(refs, ", "))
	ui.Hint("use absolute URLs or data: URIs")
}

func artifactTarget(art artifacts.Artifact) string {
	cfg, err := config.LoadMyceliumConfig()
	if err == nil && cfg != nil && cfg.ServerURL() != "" {
		return artifacts.URL(cfg.ServerURL(), art.ID)
	}
	return art.Path
}

func artifactRecorded(art artifacts.Artifact, noOpen bool) {
	ui.Success("Recorded %s", art.ID)
	if err := pushAfterWrite(); err != nil {
		ui.Warn("Recorded here but not synced (%v)", err)
	}
	target := artifactTarget(art)
	ui.Hint("%s", target)
	if noOpen {
		return
	}
	showInBrowser(target)
}

// showInBrowser opens what was just printed, and says why it did not when this
// machine has no browser. The link is already on stdout by then, so nothing has
// failed: a headless box gets the link and one line explaining the silence,
// never an error and never an exit code.
//
// That line goes to stderr because the link is the only thing `artifact open`
// owes a pipe, and a note about this machine's hardware is not data.
func showInBrowser(target string) {
	if !browser.Available() {
		ui.ErrorHint("no browser on this machine")
		return
	}
	if err := browser.Open(target); err != nil {
		ui.Warn("Could not open a browser (%v)", err)
	}
}

func init() {
	var artifactTitle string
	var artifactExpires string
	var artifactNoOpen bool
	var artifactBodyStdin bool
	artifactAddCmd.Flags().StringVar(&artifactTitle, "title", "", "Override the document's own title")
	artifactAddCmd.Flags().StringVar(&artifactExpires, "expires", "", "Ignored: artifacts are permanent")
	_ = artifactAddCmd.Flags().MarkHidden("expires")
	artifactAddCmd.Flags().BoolVar(&artifactNoOpen, "no-open", false, "Record it without opening a browser")
	artifactAddCmd.Flags().BoolVar(&artifactBodyStdin, "body-stdin", false, "Read artifact body from stdin")
	artifactAddCmd.Flags().BoolVar(&artifactBodyStdin, "stdin", false, "Read artifact body from stdin")
	var artifactJSON bool
	artifactListCmd.Flags().BoolVar(&artifactJSON, "json", false, "Print the listing as JSON")
	artifactCmd.AddCommand(artifactAddCmd, artifactListCmd, artifactOpenCmd, artifactRmCmd)
	rootCmd.AddCommand(artifactCmd)
}
