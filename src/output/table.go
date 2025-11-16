package output

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"supalink/src/utils"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func PrintAsTable(matchingPathsAndDestinations map[string]string) {
	table := table.
		New().
		Headers("Source", "Destination").
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(accentColor)).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 1)

			if row == table.HeaderRow {
				style = style.Bold(true).Foreground(accentColor)
				return style
			}

			if row%2 == 0 {
				style = style.Foreground(lightGrayColor)
			} else {
				style = style.Foreground(whiteColor)
			}
			return style
		})

	allPaths := make([]string, 0)
	for source, destination := range matchingPathsAndDestinations {
		allPaths = append(allPaths, source, destination)
	}
	rootDirectory := utils.FindRootDirectoryOfAllPaths(allPaths)

	orderedMatchingPathsAndDestinations := make(map[string]string, 0)
	orderedSources := make([]string, 0)

	for source := range matchingPathsAndDestinations {
		orderedSources = append(orderedSources, source)
	}

	sort.Strings(orderedSources)
	for _, source := range orderedSources {
		orderedMatchingPathsAndDestinations[source] = matchingPathsAndDestinations[source]
	}

	for source, destination := range orderedMatchingPathsAndDestinations {
		source = strings.TrimPrefix(source, rootDirectory+string(os.PathSeparator))
		destination = strings.TrimPrefix(destination, rootDirectory+string(os.PathSeparator))

		if len(source) > 45 {
			extension := path.Ext(source)
			source = source[:40] + "(...)" + extension
		}

		if len(destination) > 45 {
			extension := path.Ext(destination)
			destination = destination[:40] + "(...)" + extension
		}

		table.Row(source, destination)
	}

	fmt.Println(table)
}
