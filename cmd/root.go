package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/FacileStudio/Mycelium/internal/ui"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "mycelium",
	Short: "Shared agent memory across AI coding agents and machines",
	Long:  "Mycelium manages a canonical source of truth for agent memory, rules, and skills. It generates per-agent configs via thin adapters and syncs across machines.",
}

func init() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable colored output")
	var flagSpace string
	rootCmd.PersistentFlags().StringVar(&flagSpace, "space", "", "Target memory space ID or name")
	cobra.OnInitialize(func() {
		if v, _ := rootCmd.PersistentFlags().GetBool("no-color"); v {
			ui.DisableColor()
		}
		if flagSpace != "" {
			_ = os.Setenv("MYCELIUM_SPACE", flagSpace)
		}
	})
	_ = artifactCmd
	_ = claimCmd
	_ = daemonCmd
	_ = diffCmd
	_ = doctorCmd
	_ = flowCmd
	_ = initCmd
	newInstallCmd()
	_ = loginCmd
	_ = logoutCmd
	_ = mcpCmd
	_ = memoryCmd
	_ = memoryAddCmd
	_ = memoryLogCmd
	_ = memoryRatifyCmd
	_ = recapCmd
	_ = rulesCmd
	newServeCmd()
	_ = sessionsCmd
	_ = statusCmd
	_ = syncCmd
	_ = usageCmd
}

// Execute runs the root command and exits non-zero on failure.
// Trigger rebuild for Dokploy.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
