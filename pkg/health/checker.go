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

func (h *DomainHealth) RenderReport() string {
	var b strings.Builder

	title := headerTitleStyle.Render("🩺 Domain Health Report")
	domain := domainStyle.Render(h.Domain)
	b.WriteString(headerCardStyle.Render(title+"\n"+domain) + "\n\n")

	b.WriteString(renderStatus("Live (A/AAAA)", h.HasARecord) + "\n")
	b.WriteString(renderStatus("Email Receiver (MX)", h.HasMXRecord) + "\n")
	b.WriteString(renderStatus("Sender Policy (SPF)", h.HasSPF) + "\n")
	b.WriteString(renderStatus("Domain Auth (DMARC)", h.HasDMARC) + "\n")

	b.WriteString("\n" + dividerStyle.Render(strings.Repeat("─", 50)) + "\n")

	passed := 0
	for _, c := range []bool{h.HasARecord, h.HasMXRecord, h.HasSPF, h.HasDMARC} {
		if c {
			passed++
		}
	}
	b.WriteString(renderScoreBar(passed, 4) + "\n")

	if len(h.Warnings) > 0 {
		b.WriteString("\n  " + warnTitleStyle.Render("⚠  Warnings") + "\n")
		for _, w := range h.Warnings {
			b.WriteString("    " + warnBullet.Render("·") + " " + w + "\n")
		}
	}

	if len(h.Issues) > 0 {
		b.WriteString("\n  " + critTitleStyle.Render("✖  Critical Issues") + "\n")
		for _, issue := range h.Issues {
			b.WriteString("    " + critBullet.Render("·") + " " + issue + "\n")
		}
	}

	if len(h.Issues) == 0 && len(h.Warnings) == 0 {
		b.WriteString("\n  " + okStyle.Render("✅  All checks passed. Your domain is healthy!") + "\n")
	}

	b.WriteString("\n  " + hintStyle.Render("Press Enter or Esc to exit") + "\n")
	return b.String()
}

func renderStatus(name string, passed bool) string {
	var badge string
	if passed {
		badge = passBadge.Render("✓ PASS")
	} else {
		badge = failBadge.Render("✗ FAIL")
	}
	return "  " + badge + " " + checkLabel.Render(name)
}

func renderScoreBar(passed, total int) string {
	filled := strings.Repeat("█", passed)
	empty := strings.Repeat("░", total-passed)
	bar := scoreBarFill.Render(filled) + scoreBarEmpty.Render(empty)
	label := fmt.Sprintf(" %d/%d checks passed", passed, total)
	return "  " + bar + scoreBarLabel.Render(label)
}
