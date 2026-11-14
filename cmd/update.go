package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/FacileStudio/Mycelium/internal/daemon"

	"github.com/FacileStudio/Mycelium/internal/selfupdate"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade"},
	Short:   "Update mycelium to the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Current version: %s\n", version)

		updateCheckOnly, _ := cmd.Flags().GetBool("check")
		if updateCheckOnly {
			latest, available, err := selfupdate.CheckLatest(version)
			if err != nil {
				return err
			}
			if available {
				color.Yellow("Update available: %s -> %s", version, latest)
			} else {
				color.Green("mycelium is up to date (%s).", latest)
			}
			return nil
		}

		latest, available, err := selfupdate.CheckLatest(version)
		if err != nil {
			return err
		}
		if !available {
			color.Green("mycelium is already up to date (%s).", latest)
			return nil
		}

		newVersion, err := selfupdate.Apply(version)
		if err != nil {
			if errors.Is(err, selfupdate.ErrHomebrew) {
				color.Yellow("mycelium is managed by Homebrew. Run: brew upgrade mycelium")
				return nil
			}
			return err
		}
		color.Green("Updated mycelium %s -> %s", version, newVersion)
		return refreshDaemon()
	},
}

func init() {
	var updateCheckOnly bool
	updateCmd.Flags().BoolVar(&updateCheckOnly, "check", false, "Report whether an update is available without installing")
	rootCmd.AddCommand(updateCmd)
}

// refreshDaemon rewrites the background service units with the new binary,
// so a fix to the units reaches every machine that updates. The running
// process is the old version, hence the new binary does the rewrite.
func refreshDaemon() error {
	if !daemon.Installed() {
		return nil
	}
	self, err := os.Executable()
	if err != nil {
		return err
	}
	out, err := exec.Command(self, "daemon", "install").CombinedOutput()
	if err != nil {
		return fmt.Errorf("updated, but reinstalling the background sync failed: %v: %s", err, out)
	}
	fmt.Println("Background sync units refreshed.")
	return nil
}
