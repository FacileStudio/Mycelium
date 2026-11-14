package cmd

import (
	"fmt"
	"strings"

	"github.com/FacileStudio/Mycelium/internal/config"
	"github.com/FacileStudio/Mycelium/internal/memory"
	"github.com/FacileStudio/Mycelium/internal/ui"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var memorySearchLimit int

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Manage memory",
}

var memorySearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search memory, best match first",
	Long: "Search memory, best match first.\n\n" +
		"Ranks the findings of the synced markdown folder by keyword (BM25, prefix match, " +
		"accents folded). It reads local files only, so it works offline.",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		results, err := memory.SearchChunks(config.MemoryDir(), strings.Join(args, " "))
		if err != nil {
			return err
		}
		printSearchResults(results)
		return nil
	},
}

func printSearchResults(results []memory.SearchResult) {
	if len(results) == 0 {
		fmt.Println("No results.")
		return
	}
	shown := results
	if memorySearchLimit > 0 && len(shown) > memorySearchLimit {
		shown = shown[:memorySearchLimit]
	}
	changed, err := memory.ChangedPages(config.DataDir())
	if err != nil {
		ui.Hint("cannot tell which normative pages changed (%v); results are unmarked", err)
	}
	for _, r := range shown {
		color.New(color.FgCyan).Printf("%s:%d ", r.Path, r.Line)
		if r.Date != "" {
			fmt.Print(ui.Dim(r.Date + " "))
		}
		if changed[r.Path] {
			color.New(color.FgYellow).Print("[changed since ratified] ")
		}
		fmt.Println(r.Content)
	}
	if len(shown) < len(results) {
		fmt.Printf("\n%d more match(es); raise --limit to see them.\n", len(results)-len(shown))
	}
}

var memoryIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Show memory index",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := memory.ReadIndex(config.MemoryDir())
		if err != nil {
			return err
		}
		fmt.Print(content)
		return nil
	},
}

func newMemoryCmd() *cobra.Command {
	memorySearchCmd.Flags().IntVar(&memorySearchLimit, "limit", 20, "Maximum results to print")
	memoryCmd.AddCommand(memorySearchCmd)
	memoryCmd.AddCommand(memoryIndexCmd)
	rootCmd.AddCommand(memoryCmd)
	return memoryCmd
}
