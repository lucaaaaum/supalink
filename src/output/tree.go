package output

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"supalink/src/utils"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/lipgloss/tree"
)

type node struct {
	value    string
	children []*node
}

func (n *node) add(path []string) {
	if len(path) == 0 {
		return
	}

	part := path[0]
	child := n.getChild(part)
	if child == nil {
		child = &node{value: part, children: make([]*node, 0)}
		n.children = append(n.children, child)
	}

	child.add(path[1:])
}

func (n *node) getChild(value string) *node {
	for _, child := range n.children {
		if child.value == value {
			return child
		}
	}
	return nil
}

func (n *node) toLipglossTree() tree.Node {
	if len(n.value) >= 45 {
		extension := path.Ext(n.value)
		n.value = n.value[:40] + "(...)" + extension
	}
	tree := tree.Root(n.value)
	for _, child := range n.children {
		tree.Child(child.toLipglossTree())
	}
	return tree
}

func PrintAsTree(matchingPathsAndDestinations map[string]string) {
	sourcePaths := make([]string, 0, len(matchingPathsAndDestinations))
	destinationPaths := make([]string, 0, len(matchingPathsAndDestinations))
	for source, destination := range matchingPathsAndDestinations {
		sourcePaths = append(sourcePaths, source)
		destinationPaths = append(destinationPaths, destination)
	}

	sourceTree := createInMemoryTree(sourcePaths).toLipglossTree()
	destinationTree := createInMemoryTree(destinationPaths).toLipglossTree()

	table := table.
		New().
		Headers("Source", "Destination").
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(accentColor)).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 1)
			if row == table.HeaderRow {
				style = style.Bold(true).Foreground(accentColor)
			}

			return style
		}).
		Rows([]string{fmt.Sprintln(sourceTree), fmt.Sprintln(destinationTree)})
	fmt.Println(table)
}

func createInMemoryTree(paths []string) *node {
	rootDirectory := utils.FindRootDirectoryOfAllPaths(paths)
	root := &node{value: rootDirectory, children: make([]*node, 0)}
	for _, path := range paths {
		relativePath, err := filepath.Rel(rootDirectory, path)
		if err != nil {
			continue
		}
		parts := strings.Split(relativePath, string(os.PathSeparator))
		root.add(parts)
	}
	return root
}
