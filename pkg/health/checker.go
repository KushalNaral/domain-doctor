package health

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

// DomainHealth contains the results of the health checks
type DomainHealth struct {
	Domain      string
	HasARecord  bool
	HasMXRecord bool
	HasSPF      bool
	HasDMARC    bool
	Issues      []string
	Warnings    []string
}

// Check performs DNS-based health checks on a domain
func Check(domain string) *DomainHealth {
	health := &DomainHealth{
		Domain: domain,
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// Check A/AAAA records (Liveness)
	go func() {
		defer wg.Done()
		ips, err := net.LookupHost(domain)
		if err == nil && len(ips) > 0 {
			health.HasARecord = true
		} else {
			health.Issues = append(health.Issues, "No A or AAAA records found. The domain might not be live.")
		}
	}()

	// Check MX records
	go func() {
		defer wg.Done()
		mxs, err := net.LookupMX(domain)
		if err == nil && len(mxs) > 0 {
			health.HasMXRecord = true
		} else {
			health.Issues = append(health.Issues, "No MX records found. The domain cannot receive email.")
		}
	}()

	// Check TXT records for SPF and DMARC
	go func() {
		defer wg.Done()
		txts, err := net.LookupTXT(domain)
		if err == nil {
			for _, txt := range txts {
				if strings.HasPrefix(txt, "v=spf1") {
					health.HasSPF = true
					if strings.Contains(txt, "+all") {
						health.Warnings = append(health.Warnings, "SPF record uses '+all', which allows anyone to send email on your behalf. Use '~all' or '-all'.")
					}
				}
			}
		} else {
			health.Issues = append(health.Issues, "Could not fetch TXT records for SPF.")
		}

		// Check DMARC (on _dmarc.domain)
		dmarcDomain := "_dmarc." + domain
		dmarcTxts, err := net.LookupTXT(dmarcDomain)
		if err == nil {
			for _, txt := range dmarcTxts {
				if strings.HasPrefix(txt, "v=DMARC1") {
					health.HasDMARC = true
					if strings.Contains(txt, "p=none") {
						health.Warnings = append(health.Warnings, "DMARC policy is set to 'p=none'. This does not prevent spoofing. Consider upgrading to 'p=quarantine' or 'p=reject'.")
					}
				}
			}
		} else {
			health.Warnings = append(health.Warnings, "No DMARC record found or could not be fetched.")
		}
	}()

	wg.Wait()

	if !health.HasSPF {
		health.Issues = append(health.Issues, "No SPF record found. Your emails are highly likely to go to spam.")
	}

	return health
}

// PrintReport prints a colorized health report to standard output
func (h *DomainHealth) PrintReport() {
	const (
		ColorReset  = "\033[0m"
		ColorRed    = "\033[31m"
		ColorGreen  = "\033[32m"
		ColorYellow = "\033[33m"
		ColorBold   = "\033[1m"
	)

	fmt.Printf("\n%s🩺 Domain Health Report for: %s%s\n", ColorBold, h.Domain, ColorReset)
	fmt.Println(strings.Repeat("-", 50))

	printStatus("Live (A/AAAA)", h.HasARecord)
	printStatus("Email Receiver (MX)", h.HasMXRecord)
	printStatus("Sender Policy (SPF)", h.HasSPF)
	printStatus("Domain Auth (DMARC)", h.HasDMARC)

	if len(h.Warnings) > 0 {
		fmt.Printf("\n%s⚠️  Warnings:%s\n", ColorYellow, ColorReset)
		for _, w := range h.Warnings {
			fmt.Printf("  - %s\n", w)
		}
	}

	if len(h.Issues) > 0 {
		fmt.Printf("\n%s❌ Critical Issues:%s\n", ColorRed, ColorReset)
		for _, idx := range h.Issues {
			fmt.Printf("  - %s\n", idx)
		}
	}

	if len(h.Issues) == 0 && len(h.Warnings) == 0 {
		fmt.Printf("\n%s✅ Your domain health is perfect!%s\n", ColorGreen, ColorReset)
	}
	fmt.Println()
}

func printStatus(name string, passed bool) {
	const (
		ColorReset = "\033[0m"
		ColorRed   = "\033[31m"
		ColorGreen = "\033[32m"
	)
	if passed {
		fmt.Printf("[ %sPASS%s ] %s\n", ColorGreen, ColorReset, name)
	} else {
		fmt.Printf("[ %sFAIL%s ] %s\n", ColorRed, ColorReset, name)
	}
}
