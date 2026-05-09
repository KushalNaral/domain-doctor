package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhanushnehru/domain-doctor/pkg/health"
	"github.com/dhanushnehru/domain-doctor/pkg/tui"
)

func main() {
	domain := flag.String("domain", "", "Domain to analyze directly, skipping the interactive prompt")
	flag.Parse()

	if *domain != "" {
		report := health.Check(*domain)
		fmt.Print(report.RenderReport(false))
		return
	}

	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
