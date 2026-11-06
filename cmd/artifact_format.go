package cmd

import (
	"fmt"
	"math"
	"time"

	"github.com/FacileStudio/Mycelium/internal/artifacts"
	"github.com/FacileStudio/Mycelium/internal/ui"
)

// artifactIDWidth is the column width the identifiers need so the titles line up.
func artifactIDWidth(all []artifacts.Artifact) int {
	width := 0
	for _, art := range all {
		if len(art.ID) > width {
			width = len(art.ID)
		}
	}
	return width
}

// artifactLine renders one listing row: what it is, which machine wrote it, and
// how old it is.
func artifactLine(art artifacts.Artifact, width int, now time.Time) string {
	meta := artifactAge(art, now)
	if art.Machine != "" {
		meta = art.Machine + ", " + meta
	}
	return fmt.Sprintf("%-*s  %s  %s", width, art.ID, art.Title, ui.Dim("("+meta+")"))
}

// artifactAge says how long ago an artifact was last written.
func artifactAge(art artifacts.Artifact, now time.Time) string {
	if d := now.Sub(art.Updated); d >= time.Minute {
		return humanDuration(d) + " ago"
	}
	return "just now"
}

// humanDuration renders a span the way somebody scanning a list wants it: one
// unit, never "743h".
//
// It rounds to nearest rather than truncating, so a span a few seconds short
// of seven days reads "7d", not "6d".
func humanDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "0m"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(math.Round(d.Minutes())))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(math.Round(d.Hours())))
	default:
		return fmt.Sprintf("%dd", int(math.Round(d.Hours()/24)))
	}
}
