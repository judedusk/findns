package tui

import (
	"fmt"
	"strings"

	"github.com/SamNet-dev/findns/internal/data"
	"github.com/SamNet-dev/findns/internal/scanner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

type inputOption struct {
	label string
	desc  string
	group string // section header (shown before first item in group)
}

const (
	inputResolvers = iota
	inputCIDRLight
	inputCIDRMedium
	inputCIDRFull
	inputCombinedLight
	inputCombinedMedium
	inputCustom
	numInputChoices
)

var inputOptions = []inputOption{
	{
		label: "Known resolvers",
		desc:  "7,854 pre-verified Iranian DNS resolvers",
		group: "Bundled lists (embedded in binary)",
	},
	{
		label: "CIDR scan — light",
		desc:  "~19K IPs  (10 random samples per CIDR)",
		group: "Iranian IP range scan (1,919 CIDR blocks)",
	},
	{
		label: "CIDR scan — medium",
		desc:  "~96K IPs  (50 random samples per CIDR)",
	},
	{
		label: "CIDR scan — full",
		desc:  "~10.8M IPs  (entire Iranian IP space — very slow)",
	},
	{
		label: "Combined — light",
		desc:  "Resolvers + CIDR light  (~27K IPs)",
		group: "Combined (resolvers + CIDR samples, deduplicated)",
	},
	{
		label: "Combined — medium",
		desc:  "Resolvers + CIDR medium  (~104K IPs)",
	},
	{
		label: "Custom file",
		desc:  "Load IPs from a text file or JSON report",
		group: "Custom",
	},
}

func updateInput(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case inputLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		if len(msg.ips) == 0 {
			m.err = fmt.Errorf("no valid IPs found in input")
			return m, nil
		}
		m.ips = msg.ips
		m.screen = screenConfig
		m.cursor = 0
		m.err = nil
		return m, focusConfigInput(&m)

	case tea.KeyMsg:
		// Block input while loading
		if m.loading {
			return m, nil
		}
		// If typing custom path
		if m.typingPath {
			switch msg.String() {
			case "enter":
				path := m.pathInput.Value()
				if path == "" {
					return m, nil
				}
				m.typingPath = false
				m.loading = true
				return m, loadInputFile(path)
			case "esc":
				m.typingPath = false
				m.err = nil
				return m, nil
			}
			var cmd tea.Cmd
			m.pathInput, cmd = m.pathInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < numInputChoices-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == inputCustom {
				m.typingPath = true
				m.pathInput = textinput.New()
				m.pathInput.Placeholder = "path/to/resolvers.txt"
				m.pathInput.Focus()
				return m, m.pathInput.Cursor.BlinkCmd()
			}
			m.loading = true
			return m, selectInput(m.cursor)
		case "backspace", "left":
			m.screen = screenWelcome
			m.cursor = 0
			m.err = nil
		}
	}
	return m, nil
}

func selectInput(choice int) tea.Cmd {
	switch choice {
	case inputResolvers:
		return loadBundledResolvers
	case inputCIDRLight:
		return loadCIDRSampled(10)
	case inputCIDRMedium:
		return loadCIDRSampled(50)
	case inputCIDRFull:
		return loadCIDRSampled(0) // 0 = all IPs
	case inputCombinedLight:
		return loadCombined(10)
	case inputCombinedMedium:
		return loadCombined(50)
	case inputCustom:
		return nil // handled inline (typing path)
	}
	return nil
}

func viewInput(m Model) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(titleStyle.Render("  Select Input"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(dimStyle.Render("  Loading resolvers..."))
		b.WriteString("\n")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(redStyle.Render(fmt.Sprintf("  Error: %v", m.err)))
		b.WriteString("\n\n")
	}

	if m.typingPath {
		b.WriteString(dimStyle.Render("  Enter file path:"))
		b.WriteString("\n\n")
		b.WriteString("  " + m.pathInput.View())
		b.WriteString("\n\n")
		b.WriteString(dimStyle.Render("  enter confirm  esc cancel"))
		b.WriteString("\n")
		return b.String()
	}

	mode := "UDP"
	if m.config.DoH {
		mode = "DoH"
	}
	b.WriteString(dimStyle.Render(fmt.Sprintf("  Mode: %s", mode)))
	b.WriteString("\n\n")

	lastGroup := ""
	for i, opt := range inputOptions {
		// Section header
		if opt.group != "" && opt.group != lastGroup {
			if lastGroup != "" {
				b.WriteString("\n")
			}
			b.WriteString(dimStyle.Render(fmt.Sprintf("  — %s", opt.group)))
			b.WriteString("\n")
			lastGroup = opt.group
		}

		cursor := "  "
		labelStyle := normalStyle
		descStyle := dimStyle
		if i == m.cursor {
			cursor = "> "
			labelStyle = selectedStyle
			descStyle = normalStyle
		}
		b.WriteString(fmt.Sprintf("  %s%-22s  %s\n",
			cursor,
			labelStyle.Render(opt.label),
			descStyle.Render(opt.desc)))
	}

	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  ↑/↓ navigate  enter select  ← back  q quit"))
	b.WriteString("\n")

	return b.String()
}

// Commands

func loadBundledResolvers() tea.Msg {
	ips, err := data.IRResolvers()
	return inputLoadedMsg{ips: ips, err: err}
}

func loadCIDRSampled(samplePer int) tea.Cmd {
	return func() tea.Msg {
		cidrs, err := data.IRCIDRs()
		if err != nil {
			return inputLoadedMsg{err: err}
		}
		ips, err := data.ExpandCIDRsSampled(cidrs, samplePer)
		return inputLoadedMsg{ips: ips, err: err}
	}
}

func loadCombined(samplePer int) tea.Cmd {
	return func() tea.Msg {
		resolvers, err := data.IRResolvers()
		if err != nil {
			return inputLoadedMsg{err: err}
		}
		cidrs, err := data.IRCIDRs()
		if err != nil {
			return inputLoadedMsg{err: err}
		}
		cidrIPs, err := data.ExpandCIDRsSampled(cidrs, samplePer)
		if err != nil {
			return inputLoadedMsg{err: err}
		}
		// Deduplicate
		seen := make(map[string]struct{}, len(resolvers))
		for _, ip := range resolvers {
			seen[ip] = struct{}{}
		}
		combined := append([]string{}, resolvers...)
		for _, ip := range cidrIPs {
			if _, ok := seen[ip]; !ok {
				combined = append(combined, ip)
				seen[ip] = struct{}{}
			}
		}
		return inputLoadedMsg{ips: combined}
	}
}

func loadInputFile(path string) tea.Cmd {
	return func() tea.Msg {
		ips, err := scanner.LoadInput(path, false)
		return inputLoadedMsg{ips: ips, err: err}
	}
}
