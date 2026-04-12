package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/derKosi/Ohm/internal/generator"
	"github.com/derKosi/Ohm/internal/model"
	"github.com/derKosi/Ohm/internal/scanner"
)

var version = "0.1.0"

func isTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

const banner = `  ───┤   ⚡  O H M     ├───`

func printBanner() {
	fmt.Print(lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Render(banner))
}

func main() {
	if len(os.Args) < 2 {
		printBanner()
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "scan":
		cmdScan()
	case "generate":
		cmdGenerate()
	case "history":
		cmdHistory()
	case "stragglers":
		cmdStragglers()
	case "version", "--version", "-v":
		printVersion()
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printVersion() {
	printBanner()
	fmt.Printf("  v%s\n", version)
	fmt.Println("  Resistance against AGI bloat.")
	fmt.Println("  Designed with help of Pi Harness and GLM-5.1")
	fmt.Println("  MIT License — © 2026 Mathias Kosinski")
	fmt.Println()
	fmt.Println("  Check for new versions: https://github.com/derKosi/Ohm/releases")
}

func printHelp() {
	fmt.Println(titleStyle.Render("⚡ Ohm — AI Software Scanner"))
	fmt.Println()
	fmt.Println("Resistance against AGI bloat.")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  ohm scan              Scan system for AI software")
	fmt.Println("  ohm scan --path       Also check PATH for AI tool entries")
	fmt.Println("  ohm scan --env        Also check environment variables")
	fmt.Println("  ohm scan --shell      Also check shell profiles")
	fmt.Println("  ohm scan --deep       Thorough filesystem crawl")
	fmt.Println("  ohm scan --all        Enable all opt-in scans (--path, --env, --shell, --deep)")
	fmt.Println("  ohm scan --no-tui     Text output (no TUI)")
	fmt.Println("  ohm generate          Generate cleanup script from last scan")
	fmt.Println("  ohm stragglers        Scan for leftover files only")
	fmt.Println("  ohm history           Show removal history")
	fmt.Println("  ohm version           Show version")
	fmt.Println()
	fmt.Println("🔒 Privacy: All scanning is local. No data leaves your machine.")
	fmt.Println()
	fmt.Printf("Ohm v%s — https://github.com/derKosi/Ohm\n", version)
	fmt.Println("Check for updates: https://github.com/derKosi/Ohm/releases")
	fmt.Println()
	fmt.Println("MIT License — © 2026 Mathias Kosinski")
}

func cmdScan() {
	allOptIn := hasFlag("--all")
	opts := scanner.Options{
		ScanPATH:  hasFlag("--path") || allOptIn,
		ScanENV:   hasFlag("--env") || allOptIn,
		ScanShell: hasFlag("--shell") || allOptIn,
		ScanDeep:  hasFlag("--deep") || allOptIn,
	}

	// JSON or no-TUI mode: scan with simple dot animation
	if hasFlag("--json") || hasFlag("--no-tui") || !isTerminal() {
		scanDone := make(chan *model.ScanResult, 1)
		go func() {
			s := scanner.New(opts)
			scanDone <- s.Scan()
		}()

		animStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
		fmt.Print(animStyle.Render("  ⚡ Ohm  ══════════ Scanning"))

		var result *model.ScanResult
		ticker := time.NewTicker(500 * time.Millisecond)
	loop:
		for {
			select {
			case result = <-scanDone:
				break loop
			case <-ticker.C:
				fmt.Print(".")
			}
		}
		ticker.Stop()
		fmt.Println(" done.")

		if result.Count() == 0 {
			fmt.Println("No AI software found on this system.")
			return
		}

		// Save state
		state, _ := model.LoadState()
		state.LastScan = result.ScannedAt
		state.Findings = result.Findings
		state.Save()

		if hasFlag("--json") {
			outputJSON(result)
			return
		}

		printScanResult(result)
		return
	}

	// TUI mode: Bubble Tea runs the scan with built-in animation
	app := NewTUIScanner(opts, generator.Generate)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
	}
}

func cmdGenerate() {
	state, err := model.LoadState()
	if err != nil || len(state.Findings) == 0 {
		fmt.Println("No scan results found. Run 'ohm scan' first.")
		os.Exit(1)
	}

	result := &model.ScanResult{Findings: state.Findings}
	selected := 0
	for _, f := range result.Findings {
		if f.Selected {
			selected++
		}
	}

	if selected == 0 {
		fmt.Println("No items selected. Run 'ohm scan' and select items first.")
		os.Exit(1)
	}

	path, err := generator.Generate(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating script: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📝 Written: %s\n", path)
	fmt.Println()
	fmt.Println("⚠️  Review the script before running.")
	fmt.Println("⚠️  It may contain commands to delete config files with saved credentials.")
	fmt.Println()
	fmt.Printf("  cat %s\n", path)
}

func cmdStragglers() {
	s := scanner.New(scanner.Options{})
	result := s.Scan()

	count := 0
	for _, f := range result.Findings {
		if f.Category == model.CatStragglers {
			count++
		}
	}
	if count == 0 {
		fmt.Println("No stragglers found.")
		return
	}

	fmt.Printf("Found %d straggler(s):\n\n", count)
	for _, f := range result.Findings {
		if f.Category == model.CatStragglers {
			fmt.Printf("  %s %s\n", f.RiskLevel.Icon(), f.Name)
			fmt.Printf("     Path: %s\n", f.Path)
			fmt.Printf("     Size: %s\n", f.FormatSize())
			for _, sub := range f.SubItems {
				fmt.Printf("     - %s\n", sub)
			}
			fmt.Println()
		}
	}
}

func cmdHistory() {
	state, err := model.LoadState()
	if err != nil || len(state.Removed) == 0 {
		fmt.Println("No removal history found.")
		return
	}
	fmt.Println("Removal history:")
	for _, r := range state.Removed {
		fmt.Printf("  - %s (%s) removed at %s\n", r.Name, r.ID, r.RemovedAt.Format("2006-01-02 15:04"))
	}
}

func hasFlag(flag string) bool {
	for _, arg := range os.Args[2:] {
		if arg == flag {
			return true
		}
	}
	return false
}

// --- Styles ---

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("230")).
	Background(lipgloss.Color("62")).
	Padding(0, 2)

// --- TUI ---

type GenerateFunc func(*model.ScanResult) (string, error)

// scanCompleteMsg is sent when the background scan finishes.
type scanCompleteMsg struct{ result *model.ScanResult }

// tickMsg is sent on each animation tick while scanning.
type tickMsg struct{}

// TUIApp with viewport scrolling.
type TUIApp struct {
	scanning   bool
	dotCount   int
	result     *model.ScanResult
	flatItems  []*model.Finding
	cursor     int
	scroll     int // scroll offset for viewport
	width      int
	height     int
	generate   GenerateFunc
	scriptPath string
	err        string
}

const (
	headerLines = 4 // banner + privacy + blank + blank
	footerLines = 4 // totals + blank + keys + blank
)

// NewTUIScanner creates a TUI that runs the scan with a spinner.
func NewTUIScanner(opts scanner.Options, genFn GenerateFunc) *TUIApp {
	return &TUIApp{
		scanning: true,
		generate: genFn,
	}
}

// NewTUIApp creates a TUI with pre-loaded results.
func NewTUIApp(result *model.ScanResult, genFn GenerateFunc) *TUIApp {
	var flat []*model.Finding
	for i := range result.Findings {
		flat = append(flat, &result.Findings[i])
	}
	return &TUIApp{
		result:    result,
		flatItems: flat,
		generate:  genFn,
	}
}

func (a *TUIApp) Init() tea.Cmd {
	if a.scanning {
		return tea.Batch(a.startScan(), a.tick())
	}
	return nil
}

func (a *TUIApp) startScan() tea.Cmd {
	return func() tea.Msg {
		s := scanner.New(scanner.Options{
			ScanPATH:  hasFlag("--path") || hasFlag("--all"),
			ScanENV:   hasFlag("--env") || hasFlag("--all"),
			ScanShell: hasFlag("--shell") || hasFlag("--all"),
			ScanDeep:  hasFlag("--deep") || hasFlag("--all"),
		})
		return scanCompleteMsg{result: s.Scan()}
	}
}

func (a *TUIApp) tick() tea.Cmd {
	return tea.Tick(300*time.Millisecond, func(_ time.Time) tea.Msg { return tickMsg{} })
}

func (a *TUIApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle scan-in-progress state
	if a.scanning {
		switch msg := msg.(type) {
		case scanCompleteMsg:
			a.scanning = false
			a.result = msg.result
			var flat []*model.Finding
			for i := range a.result.Findings {
				flat = append(flat, &a.result.Findings[i])
			}
			a.flatItems = flat
			// Save state
			state, _ := model.LoadState()
			state.LastScan = a.result.ScannedAt
			state.Findings = a.result.Findings
			state.Save()
			return a, nil
		case tea.WindowSizeMsg:
			a.width = msg.Width
			a.height = msg.Height
		case tea.KeyMsg:
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return a, tea.Quit
			}
		case tickMsg:
			a.dotCount++
			return a, a.tick()
		}
		return a, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.clampScroll()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "up", "k":
			if a.cursor > 0 {
				a.cursor--
				a.clampScroll()
			}
		case "down", "j":
			if a.cursor < len(a.flatItems)-1 {
				a.cursor++
				a.clampScroll()
			}
		case "pgup":
			a.cursor -= a.viewportHeight()
			if a.cursor < 0 {
				a.cursor = 0
			}
			a.clampScroll()
		case "pgdown":
			a.cursor += a.viewportHeight()
			if a.cursor > len(a.flatItems)-1 {
				a.cursor = len(a.flatItems) - 1
			}
			a.clampScroll()
		case "home":
			a.cursor = 0
			a.scroll = 0
		case "end":
			a.cursor = len(a.flatItems) - 1
			a.clampScroll()
		case " ":
			if a.cursor < len(a.flatItems) {
				a.flatItems[a.cursor].Selected = !a.flatItems[a.cursor].Selected
			}
		case "a":
			allSelected := true
			for _, f := range a.flatItems {
				if !f.Selected {
					allSelected = false
					break
				}
			}
			for _, f := range a.flatItems {
				f.Selected = !allSelected
			}
		case "g", "enter":
			if a.result.SelectedCount() > 0 {
				path, err := a.generate(a.result)
				if err != nil {
					a.err = err.Error()
				} else {
					a.scriptPath = path
					a.err = ""
				}
			}
		}
	}
	return a, nil
}

// viewportHeight returns the number of lines available for items.
func (a *TUIApp) viewportHeight() int {
	h := a.height - headerLines - footerLines
	if h < 1 {
		return 1
	}
	return h
}

// cursorDisplayLine computes which display line the cursor item appears on.
func (a *TUIApp) cursorDisplayLine() int {
	groups := a.result.ByCategory()
	itemIdx := 0
	lineIdx := 0
	for gi := range groups {
		group := &groups[gi]
		if len(group.FindingIdxs) == 0 {
			continue
		}
		lineIdx++ // category header
		for range group.FindingIdxs {
			if itemIdx == a.cursor {
				return lineIdx
			}
			itemIdx++
			lineIdx++
		}
		lineIdx++ // blank after category
	}
	return lineIdx
}

// totalDisplayLines returns the total number of display lines.
func (a *TUIApp) totalDisplayLines() int {
	groups := a.result.ByCategory()
	total := 0
	for gi := range groups {
		group := &groups[gi]
		if len(group.FindingIdxs) == 0 {
			continue
		}
		total++ // category header
		total += len(group.FindingIdxs)
		total++ // blank after category
	}
	return total
}

// clampScroll ensures the viewport follows the cursor.
func (a *TUIApp) clampScroll() {
	vh := a.viewportHeight()
	cursorLine := a.cursorDisplayLine()
	if a.scroll > cursorLine {
		a.scroll = cursorLine
	}
	if a.scroll+vh <= cursorLine {
		a.scroll = cursorLine - vh + 1
	}
	if a.scroll < 0 {
		a.scroll = 0
	}
}

func (a *TUIApp) View() string {
	// Scanning state — show spinner
	if a.scanning {
		dots := strings.Repeat(".", a.dotCount%20)
		return "\n  ⚡ Ohm  ══════════ Scanning" + dots + "\n\n  🔒 All scanning is local. No data leaves this machine.\n"
	}

	var sb strings.Builder

	// --- Fixed Header ---
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Render("  ───┤   ⚡  O H M     ├───"))
	sb.WriteString("\n")
	sb.WriteString(helpStyle("🔒 All scanning is local. No data leaves this machine."))
	sb.WriteString("\n")

	// Platform warnings (e.g. WSL detected)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	for _, w := range a.result.Warnings {
		sb.WriteString(warnStyle.Render("⚠️  " + w))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// --- Scrollable Item List ---
	catStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	selStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	szStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("180"))

	groups := a.result.ByCategory()
	vh := a.viewportHeight()
	scrollEnd := a.scroll + vh

	itemIdx := 0 // indexes flatItems (cursor target)
	lineIdx := 0 // indexes display lines (headers + items + blanks)
	for gi := range groups {
		group := &groups[gi]
		if len(group.FindingIdxs) == 0 {
			continue
		}

		// Render emoji outside of lipgloss to avoid width calculation bugs
		catIcon := group.Category.Icon()
		catText := fmt.Sprintf("%s (%d found, %s)",
			group.Category,
			group.Count(),
			model.FormatBytes(group.TotalSize(a.result.Findings)),
		)

		// Category header — not selectable, skip in cursor
		if lineIdx >= a.scroll && lineIdx < scrollEnd {
			sb.WriteString(catIcon + " " + catStyle.Render(catText))
			sb.WriteString("\n")
		}
		lineIdx++

		for _, fi := range group.FindingIdxs {
			f := &a.result.Findings[fi]
			isCursor := itemIdx == a.cursor // compare finding index, not line index

			check := "[ ]"
			if f.Selected {
				check = "[x]"
			}

			// Don't render emojis through lipgloss — it miscalculates width
			riskIcon := "  "
			switch f.RiskLevel {
			case model.RiskDanger:
				riskIcon = "🔑"
			case model.RiskCaution:
				riskIcon = "⚠️ "
			}

			sizeStr := szStyle.Render(fmt.Sprintf("%-10s", f.FormatSize()))
			pathStr := f.Path
			if len(pathStr) > 45 {
				pathStr = "..." + pathStr[len(pathStr)-42:]
			}

			cursorStr := " "
			if isCursor {
				cursorStr = ">"
			}

			// Style text parts individually, keep emojis raw
			styledName := fmt.Sprintf("%-30s", f.Name)
			styledPath := pathStr
			styledCursor := cursorStr
			styledCheck := check
			if isCursor {
				styledName = selStyle.Render(styledName)
				styledPath = selStyle.Render(styledPath)
				styledCursor = selStyle.Render(styledCursor)
				styledCheck = selStyle.Render(styledCheck)
			} else {
				styledName = itemStyle.Render(styledName)
				styledPath = itemStyle.Render(styledPath)
			}

			line := fmt.Sprintf("%s %s %s %s %s %s",
				styledCursor, styledCheck, riskIcon, styledName, sizeStr, styledPath)

			if lineIdx >= a.scroll && lineIdx < scrollEnd {
				sb.WriteString(line)
				sb.WriteString("\n")
			}
			itemIdx++
			lineIdx++
		}

		// Blank line after category — not selectable
		if lineIdx >= a.scroll && lineIdx < scrollEnd {
			sb.WriteString("\n")
		}
		lineIdx++
	}

	// Pad remaining viewport with empty lines so footer stays fixed
	remaining := scrollEnd - lineIdx
	if remaining > 0 {
		for i := 0; i < remaining; i++ {
			sb.WriteString("\n")
		}
	} else if lineIdx < scrollEnd {
		sb.WriteString("\n")
	}

	// --- Fixed Footer ---
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
	sb.WriteString(footer.Render(fmt.Sprintf(
		"Total: %d items (%s) | Selected: %d items (%s)",
		a.result.Count(),
		model.FormatBytes(a.result.TotalSize()),
		a.result.SelectedCount(),
		model.FormatBytes(a.result.SelectedSize()),
	)))
	sb.WriteString("\n")

	if a.scriptPath != "" {
		sb.WriteString(fmt.Sprintf("📝 Script written: %s\n", a.scriptPath))
		sb.WriteString(helpStyle("   Review with: cat " + a.scriptPath))
		sb.WriteString("\n")
	}
	if a.err != "" {
		sb.WriteString(fmt.Sprintf("❌ Error: %s\n", a.err))
	}

	// Scroll indicator
	scrollInfo := ""
	if len(a.flatItems) > vh {
		scrollInfo = fmt.Sprintf(" [%d/%d] ", a.cursor+1, len(a.flatItems))
	}

	sb.WriteString(helpStyle(fmt.Sprintf("↑/k up • ↓/j down • pgup/pgdn • space select • a toggle all • g generate • q quit%s", scrollInfo)))
	sb.WriteString("\n")
	sb.WriteString(helpStyle(fmt.Sprintf("Ohm v%s · MIT License © 2026 Mathias Kosinski · Updates: github.com/derKosi/Ohm/releases", version)))
	sb.WriteString("\n")

	return sb.String()
}

func helpStyle(s string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(s)
}

// --- Text Output (no TUI) ---

func outputJSON(result *model.ScanResult) {
	type jsonOutput struct {
		Version   string           `json:"version"`
		ScannedAt time.Time        `json:"scanned_at"`
		Count     int              `json:"count"`
		Findings  []model.Finding  `json:"findings"`
	}

	out := jsonOutput{
		Version:   version,
		ScannedAt: result.ScannedAt,
		Count:     result.Count(),
		Findings:  result.Findings,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "JSON encode error: %v\n", err)
		os.Exit(1)
	}
}

func printScanResult(result *model.ScanResult) {
	printBanner()
	fmt.Println()
	groups := result.ByCategory()

	for _, group := range groups {
		if len(group.FindingIdxs) == 0 {
			continue
		}

		header := fmt.Sprintf("%s %s (%d found, %s)",
			group.Category.Icon(),
			group.Category,
			group.Count(),
			model.FormatBytes(group.TotalSize(result.Findings)),
		)
		fmt.Println(header)

		for _, fi := range group.FindingIdxs {
			f := &result.Findings[fi]
			riskIcon := "  "
			switch f.RiskLevel {
			case model.RiskDanger:
				riskIcon = "🔑"
			case model.RiskCaution:
				riskIcon = "⚠️ "
			}

			fmt.Printf("   %s %-30s %-10s %s\n",
				riskIcon, f.Name, f.FormatSize(), f.Path)

			for _, sub := range f.SubItems {
				fmt.Printf("     - %s\n", sub)
			}
		}
		fmt.Println()
	}

	fmt.Printf("Total: %d items (%s)\n", result.Count(), model.FormatBytes(result.TotalSize()))

	if len(result.Warnings) > 0 {
		fmt.Println()
		for _, w := range result.Warnings {
			fmt.Printf("⚠️  %s\n", w)
		}
	}

	fmt.Println()
	fmt.Println("Run 'ohm scan' (without --no-tui) for interactive selection.")
	fmt.Println("Run 'ohm generate' to create a cleanup script from the last scan.")
}
