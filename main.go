package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dhanushnehru/domain-doctor/pkg/health"
)

func main() {
	domainPtr := flag.String("domain", "", "The domain name to analyze (e.g., example.com)")
	flag.Parse()

	if *domainPtr == "" {
		fmt.Println("Error: --domain flag is required.")
		fmt.Println("Usage: domain-doctor --domain <your-domain.com>")
		os.Exit(1)
	}

	// In the future, we could add a loading spinner here
	fmt.Printf("Analyzing %s...\n", *domainPtr)

	report := health.Check(*domainPtr)
	report.PrintReport()
}
