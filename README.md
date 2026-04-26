# domain-doctor

🩺 A fast, zero-config CLI tool to instantly diagnose why your domain's emails are going to spam. Checks DNS, MX, SPF, DKIM, and DMARC in milliseconds.

## Why this exists

Setting up email deliverability is a notoriously painful process. If you configure your DNS records slightly wrong, your crucial transactional emails will silently drop into the spam folder.

`domain-doctor` is an open-source, terminal-based health check that doesn't just tell you what your records are—it grades them, tells you what you're doing wrong, and provides actionable advice.

## Installation

Ensure you have Go installed, then run:

```bash
go install github.com/dhanushnehru/domain-doctor@latest
```

## Usage

Run the doctor against any domain:

```bash
domain-doctor --domain github.com
```

### Example Output

```text
🩺 Domain Health Report for: example.com
--------------------------------------------------
[ PASS ] Live (A/AAAA)
[ PASS ] Email Receiver (MX)
[ PASS ] Sender Policy (SPF)
[ FAIL ] Domain Auth (DMARC)

⚠️  Warnings:
  - DMARC policy is set to 'p=none'. This does not prevent spoofing. Consider upgrading to 'p=quarantine' or 'p=reject'.

❌ Critical Issues:
  - No DMARC record found or could not be fetched.
```

## Contributing / Help Wanted! 🚀

This tool is built to be extremely extensible. If you want to contribute, we would love your help! Check out the issues tab, or tackle one of these high-priority roadmap items:

- **[Feature] Add `--format=json` flag:** Allow the tool to output structured JSON for CI/CD pipeline integration.
- **[Feature] DKIM Validation:** Add a check to validate DKIM records using a provided selector.
- **[Enhancement] Lipgloss UI:** Upgrade the terminal UI using `github.com/charmbracelet/lipgloss` for a beautiful, modern terminal experience.
- **[Feature] CI/CD Github Action:** Create a GitHub action that wraps this tool so teams can monitor their domain health weekly.

## License
MIT
