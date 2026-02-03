package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// --- Styling Definitions ---
var (
	// Stark Palette
	cCyan   = lipgloss.Color("#00F0FF")
	cBlue   = lipgloss.Color("#0077BE")
	cOrange = lipgloss.Color("#FF5F1F")
	cDark   = lipgloss.Color("#1A1A1A")
	cDim    = lipgloss.Color("#444444")

	// Nord Palette
	nordGreen = lipgloss.Color("#A3BE8C")
	nordTeal  = lipgloss.Color("#8FBCBB")
	nordBlue  = lipgloss.Color("#81A1C1")
	nordDark  = lipgloss.Color("#2E3440")

	// Alert Colors
	alertRed    = lipgloss.Color("#FF4444")
	alertYellow = lipgloss.Color("#FFD700")
	alertGreen  = lipgloss.Color("#44FF44")
	pulsePurple = lipgloss.Color("#9B59B6")
	gridColor   = lipgloss.Color("#004444")

	// Layout Styles
	docStyle = lipgloss.NewStyle().Padding(1, 2).Background(cDark)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cCyan).
			Padding(1).
			Background(cDark)

	headerStyle = lipgloss.NewStyle().
			Foreground(cDark).
			Background(cCyan).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	// Text Styles
	logLabel    = lipgloss.NewStyle().Foreground(cOrange).Bold(true)
	logText     = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	coreStyle   = lipgloss.NewStyle().Foreground(cCyan)
	clockStyle  = lipgloss.NewStyle().Foreground(cCyan).Bold(true).Padding(0, 1)
	modeStyle   = lipgloss.NewStyle().Foreground(alertGreen).Bold(true).Padding(0, 1)
	alertStyle  = lipgloss.NewStyle().Foreground(alertRed).Background(lipgloss.Color("#330000")).Bold(true).Padding(0, 1)
	scanlineDim = lipgloss.NewStyle().Foreground(lipgloss.Color("#111111"))
	glitchStyle = lipgloss.NewStyle().Foreground(alertYellow)

	// HUD Badge Styles
	badgeGreen  = lipgloss.NewStyle().Background(alertGreen).Foreground(cDark).Padding(0, 1)
	badgeYellow = lipgloss.NewStyle().Background(alertYellow).Foreground(cDark).Padding(0, 1)
	badgeRed    = lipgloss.NewStyle().Background(alertRed).Foreground(cDark).Padding(0, 1)
	badgePurple = lipgloss.NewStyle().Background(pulsePurple).Foreground(cDark).Padding(0, 1)
)

// --- Theme System ---
type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Background lipgloss.Color
	Dim        lipgloss.Color
	Alert      lipgloss.Color
}

var themes = []Theme{
	{
		Name:       "STARK",
		Primary:    lipgloss.Color("#00F0FF"),
		Secondary:  lipgloss.Color("#0077BE"),
		Accent:     lipgloss.Color("#FF5F1F"),
		Background: lipgloss.Color("#1A1A1A"),
		Dim:        lipgloss.Color("#444444"),
		Alert:      lipgloss.Color("#FF4444"),
	},
	{
		Name:       "ARC REACTOR",
		Primary:    lipgloss.Color("#00D9FF"),
		Secondary:  lipgloss.Color("#0099FF"),
		Accent:     lipgloss.Color("#FFFFFF"),
		Background: lipgloss.Color("#0A0A1A"),
		Dim:        lipgloss.Color("#334466"),
		Alert:      lipgloss.Color("#00FFFF"),
	},
	{
		Name:       "STEALTH",
		Primary:    lipgloss.Color("#00FF00"),
		Secondary:  lipgloss.Color("#006600"),
		Accent:     lipgloss.Color("#88FF88"),
		Background: lipgloss.Color("#0A0A0A"),
		Dim:        lipgloss.Color("#223322"),
		Alert:      lipgloss.Color("#FFFF00"),
	},
	{
		Name:       "NEON CITY",
		Primary:    lipgloss.Color("#FF00FF"),
		Secondary:  lipgloss.Color("#9B59B6"),
		Accent:     lipgloss.Color("#00FFFF"),
		Background: lipgloss.Color("#1A0A1A"),
		Dim:        lipgloss.Color("#442244"),
		Alert:      lipgloss.Color("#FF0099"),
	},
	{
		Name:       "WAR MACHINE",
		Primary:    lipgloss.Color("#C0C0C0"),
		Secondary:  lipgloss.Color("#808080"),
		Accent:     lipgloss.Color("#FF0000"),
		Background: lipgloss.Color("#0A0A0A"),
		Dim:        lipgloss.Color("#404040"),
		Alert:      lipgloss.Color("#FF3333"),
	},
	{
		Name:       "RESCUE",
		Primary:    lipgloss.Color("#FFD700"),
		Secondary:  lipgloss.Color("#FFA500"),
		Accent:     lipgloss.Color("#FFFFFF"),
		Background: lipgloss.Color("#1A1410"),
		Dim:        lipgloss.Color("#665533"),
		Alert:      lipgloss.Color("#FF6600"),
	},
}

func (m model) getTheme() Theme {
	return themes[m.currentTheme%len(themes)]
}

// --- Messages ---
type tickMsg time.Time
type logMsg string

// --- Model ---
type model struct {
	width, height int

	// Components
	spinner  spinner.Model
	cpuBar   progress.Model
	pwrBar   progress.Model
	netBar   progress.Model
	viewport viewport.Model

	// Data
	logs      []string
	cpuVal    float64
	pwrVal    float64
	netVal    float64
	booted    bool
	resonance []float64

	// Matrix Data
	matrixCols  int
	matrixRows  int
	matrixGrid  [][]rune
	matrixHeads []int
	matrixTails []int
	matrixSpeed []int

	// HUD Features
	currentMode     string
	tickCount       int
	glitchActive    bool
	alertActive     bool
	alertMessage    string
	alertSeverity   int
	pulsePhase      float64
	scanlinePos     int
	targetAngles    []float64
	dataStreamChars []string
	hologramChars   []string

	// Interactive Controls
	paused          bool
	showHelp        bool
	currentTheme    int
	showSoundWave   bool
	audioLevels     []float64
	arcReactorPhase float64

	// Boot Sequence
	bootPhase    int
	bootComplete bool
	bootMessage  string
	systemScan   bool
	scanProgress float64
}

func initialModel() model {
	// 1. Arc Reactor Spinner
	s := spinner.New()
	s.Spinner = spinner.Globe
	s.Style = lipgloss.NewStyle().Foreground(cCyan)

	// 2. Progress Bars
	p1 := progress.New(progress.WithGradient(string(cBlue), string(cCyan)))
	p2 := progress.New(progress.WithGradient(string(cOrange), string(lipgloss.Color("#FF0000"))))
	p3 := progress.New(progress.WithGradient(string(lipgloss.Color("#00FF00")), string(cBlue)))

	// 3. Viewport (Log Stream)
	vp := viewport.New(40, 15) // Size updated on window resize

	return model{
		spinner:         s,
		cpuBar:          p1,
		pwrBar:          p2,
		netBar:          p3,
		viewport:        vp,
		logs:            []string{"Initializing J.A.R.V.I.S. Protocol..."},
		cpuVal:          0.2,
		pwrVal:          0.8,
		netVal:          0.5,
		resonance:       make([]float64, 20),
		matrixCols:      0,
		matrixRows:      0,
		currentMode:     "FLIGHT",
		tickCount:       0,
		glitchActive:    false,
		alertActive:     false,
		alertMessage:    "",
		alertSeverity:   0,
		pulsePhase:      0,
		scanlinePos:     0,
		targetAngles:    []float64{0, 120, 240},
		dataStreamChars: []string{"⬡", "⬢", "◈", "◇", "◆", "◊"},
		hologramChars:   []string{"█", "▓", "▒", "░", "▄", "▀"},
		paused:          false,
		showHelp:        false,
		currentTheme:    0,
		showSoundWave:   true,
		audioLevels:     make([]float64, 16),
		arcReactorPhase: 0,
		bootPhase:       0,
		bootComplete:    false,
		bootMessage:     "Initializing J.A.R.V.I.S. Protocol...",
		systemScan:      false,
		scanProgress:    0,
	}
}

// katakana returns a random half-width katakana rune or a digit/latin char
func randomMatrixChar() rune {
	// Half-width Katakana: FF66-FF9D
	// Digits: 0030-0039
	// Mix: 80% Katakana, 20% Digits/Latin
	if rand.Float64() < 0.8 {
		return rune(0xFF66 + rand.Intn(0xFF9D-0xFF66+1))
	}
	return rune('0' + rand.Intn(10))
}

// updateSystemStats fetches real system statistics
func (m *model) updateSystemStats() {
	// CPU Usage
	if percentages, err := cpu.Percent(0, false); err == nil && len(percentages) > 0 {
		m.cpuVal = percentages[0] / 100.0
	}

	// Memory Usage
	if vmem, err := mem.VirtualMemory(); err == nil {
		m.pwrVal = vmem.UsedPercent / 100.0
	}

	// Network I/O (simplified - just check if there's activity)
	if netStats, err := net.IOCounters(false); err == nil && len(netStats) > 0 {
		// Use bytes sent + received as a rough indicator
		totalBytes := float64(netStats[0].BytesSent + netStats[0].BytesRecv)
		// Normalize to 0-1 range (this is a simplified approach)
		m.netVal = math.Min(1.0, math.Mod(totalBytes, 1000000)/1000000.0)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tickCommand(),
		generateLogCommand(),
	)
}

// --- Logic ---

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "h", "?":
			m.showHelp = !m.showHelp

		case "t":
			m.currentTheme = (m.currentTheme + 1) % len(themes)
			theme := m.getTheme()
			newLog := fmt.Sprintf("Theme switched to: %s", theme.Name)
			m.logs = append(m.logs, logLabel.Render(">>")+" "+logText.Render(newLog))
			if len(m.logs) > 50 {
				m.logs = m.logs[1:]
			}
			m.viewport.SetContent(strings.Join(m.logs, "\n"))
			m.viewport.GotoBottom()

		case "p":
			m.paused = !m.paused
			status := "PAUSED"
			if !m.paused {
				status = "RESUMED"
			}
			newLog := fmt.Sprintf("System %s", status)
			m.logs = append(m.logs, logLabel.Render(">>")+" "+logText.Render(newLog))
			if len(m.logs) > 50 {
				m.logs = m.logs[1:]
			}
			m.viewport.SetContent(strings.Join(m.logs, "\n"))
			m.viewport.GotoBottom()

		case "s":
			m.showSoundWave = !m.showSoundWave
			status := "ENABLED"
			if !m.showSoundWave {
				status = "DISABLED"
			}
			newLog := fmt.Sprintf("Sound visualization %s", status)
			m.logs = append(m.logs, logLabel.Render(">>")+" "+logText.Render(newLog))
			if len(m.logs) > 50 {
				m.logs = m.logs[1:]
			}
			m.viewport.SetContent(strings.Join(m.logs, "\n"))
			m.viewport.GotoBottom()

		case "r":
			m.tickCount = 0
			m.pulsePhase = 0
			m.arcReactorPhase = 0
			newLog := "System reboot initiated"
			m.logs = append(m.logs, logLabel.Render(">>")+" "+logText.Render(newLog))
			if len(m.logs) > 50 {
				m.logs = m.logs[1:]
			}
			m.viewport.SetContent(strings.Join(m.logs, "\n"))
			m.viewport.GotoBottom()

		case " ":
			newLog := "Manual system scan initiated"
			m.logs = append(m.logs, logLabel.Render(">>")+" "+logText.Render(newLog))
			if len(m.logs) > 50 {
				m.logs = m.logs[1:]
			}
			m.viewport.SetContent(strings.Join(m.logs, "\n"))
			m.viewport.GotoBottom()

		case "up":
			m.viewport.LineUp(1)

		case "down":
			m.viewport.LineDown(1)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Responsive resizing
		colWidth := (m.width / 3) - 4
		m.cpuBar.Width = colWidth - 10
		m.pwrBar.Width = colWidth - 10
		m.netBar.Width = colWidth - 10
		m.viewport.Width = colWidth
		m.viewport.Height = m.height - 10

		// Resize Matrix
		// Ensure we have a grid that covers the center panel
		m.matrixCols = colWidth - 2
		m.matrixRows = m.height - 14 // Leave room for Arc Reactor (Spinner + Text) + Borders
		if m.matrixRows < 0 {
			m.matrixRows = 0
		}

		if len(m.matrixHeads) != m.matrixCols {
			// Re-initialize if width changed
			m.matrixGrid = make([][]rune, m.matrixCols)
			m.matrixHeads = make([]int, m.matrixCols)
			m.matrixTails = make([]int, m.matrixCols)
			m.matrixSpeed = make([]int, m.matrixCols)

			for x := 0; x < m.matrixCols; x++ {
				m.matrixGrid[x] = make([]rune, m.matrixRows)
				// Randomize start pos to be scattered off-screen or mid-screen
				m.matrixHeads[x] = rand.Intn(m.matrixRows*2) - m.matrixRows
				m.matrixTails[x] = rand.Intn(10) + 5
				m.matrixSpeed[x] = rand.Intn(3) + 1 // speed 1 to 3

				// Fill grid with random chars initially
				for y := 0; y < m.matrixRows; y++ {
					m.matrixGrid[x][y] = randomMatrixChar()
				}
			}
		} else if len(m.matrixGrid[0]) != m.matrixRows {
			// Height changed, resize columns
			for x := 0; x < m.matrixCols; x++ {
				newCol := make([]rune, m.matrixRows)
				copy(newCol, m.matrixGrid[x])
				// Fill new space
				for y := len(m.matrixGrid[x]); y < m.matrixRows; y++ {
					newCol[y] = randomMatrixChar()
				}
				m.matrixGrid[x] = newCol
			}
		}

	case tickMsg:
		// Skip updates if paused
		if m.paused {
			cmds = append(cmds, tickCommand())
			return m, tea.Batch(cmds...)
		}

		// Update real system stats every 10 ticks (~2 seconds)
		if m.tickCount%10 == 0 {
			m.updateSystemStats()
		} else {
			// Add small random variations between updates for smooth animation
			m.cpuVal += (rand.Float64() - 0.5) * 0.02
			if m.cpuVal > 1 {
				m.cpuVal = 1
			}
			if m.cpuVal < 0 {
				m.cpuVal = 0
			}

			m.pwrVal += (rand.Float64() - 0.5) * 0.01
			if m.pwrVal > 1 {
				m.pwrVal = 1
			}
			if m.pwrVal < 0 {
				m.pwrVal = 0
			}

			m.netVal += (rand.Float64() - 0.5) * 0.02
			if m.netVal > 1 {
				m.netVal = 1
			}
			if m.netVal < 0 {
				m.netVal = 0
			}
		}

		// Update Resonance
		// Shift left
		if len(m.resonance) > 0 {
			m.resonance = m.resonance[1:]
			m.resonance = append(m.resonance, rand.Float64())
		}

		// Update Matrix
		// 1. Move Heads
		for x := 0; x < m.matrixCols; x++ {
			// Move down based on speed (simple frame skip logic or increment)
			// For simplicity in this TUI loop, we just increment position.
			// To vary speed, we can use rand or a counter. Let's just use speed as step size.
			m.matrixHeads[x] += 1 // Always move 1 step per tick to keep it smooth?
			// Or move by speed? Speed might be too fast.
			// Let's use a probability based on speed to simulate variable update rates per column
			// Speed 1: 33% move, Speed 2: 66% move, Speed 3: 100% move
			if rand.Intn(4) < m.matrixSpeed[x] {
				m.matrixHeads[x]++
			}

			// Reset if trail is off bottom
			if m.matrixHeads[x]-m.matrixTails[x] > m.matrixRows {
				m.matrixHeads[x] = 0 - rand.Intn(10)
				m.matrixTails[x] = rand.Intn(15) + 5
				m.matrixSpeed[x] = rand.Intn(3) + 1
			}
		}

		// 2. Glitch Grid (Randomly change characters)
		// Change ~1% of visible characters per tick
		glitchCount := (m.matrixCols * m.matrixRows) / 100
		for i := 0; i < glitchCount; i++ {
			gx := rand.Intn(m.matrixCols)
			gy := rand.Intn(m.matrixRows)
			m.matrixGrid[gx][gy] = randomMatrixChar()
		}

		// HUD Updates
		m.tickCount++
		m.pulsePhase += 0.15
		m.scanlinePos = (m.scanlinePos + 1) % 10

		// Update target angles for radar effect
		for i := range m.targetAngles {
			m.targetAngles[i] += 3
			if m.targetAngles[i] >= 360 {
				m.targetAngles[i] = 0
			}
		}

		// Update audio levels for sound wave visualization
		for i := range m.audioLevels {
			m.audioLevels[i] = rand.Float64()
		}

		// Update Arc Reactor phase for pulsing animation
		m.arcReactorPhase += 0.1

		// Periodic glitch effect (every 50 ticks)
		if m.tickCount%50 == 0 {
			m.glitchActive = !m.glitchActive
		}

		// Random mode switching
		if m.tickCount%200 == 0 {
			modes := []string{"FLIGHT", "COMBAT", "STEALTH", "ANALYSIS", "NAVIGATION"}
			m.currentMode = modes[rand.Intn(len(modes))]
		}

		// Random alert generation (rare)
		if rand.Float64() < 0.005 {
			alerts := []string{
				"DETECTING HOSTILES",
				"ENERGY SPIKE",
				"INCOMING MISSILE",
				"TARGET LOCKED",
				"SYSTEM WARNING",
			}
			m.alertActive = true
			m.alertMessage = alerts[rand.Intn(len(alerts))]
			m.alertSeverity = rand.Intn(3) + 1
		}

		// Clear alert after 5 seconds (~25 ticks at 200ms)
		if m.alertActive && m.tickCount%25 == 0 {
			m.alertActive = false
		}

		cmds = append(cmds, tickCommand())

	case logMsg:
		// Add new log entry
		newLog := fmt.Sprintf("%s %s", logLabel.Render(">>"), logText.Render(string(msg)))
		m.logs = append(m.logs, newLog)
		if len(m.logs) > 50 {
			m.logs = m.logs[1:] // Keep buffer small
		}
		m.viewport.SetContent(strings.Join(m.logs, "\n"))
		m.viewport.GotoBottom()
		cmds = append(cmds, generateLogCommand())

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case progress.FrameMsg:
		newModel, cmd := m.cpuBar.Update(msg)
		if p, ok := newModel.(progress.Model); ok {
			m.cpuBar = p
		}
		cmds = append(cmds, cmd)

		newModel2, cmd2 := m.pwrBar.Update(msg)
		if p, ok := newModel2.(progress.Model); ok {
			m.pwrBar = p
		}
		cmds = append(cmds, cmd2)

		newModel3, cmd3 := m.netBar.Update(msg)
		if p, ok := newModel3.(progress.Model); ok {
			m.netBar = p
		}
		cmds = append(cmds, cmd3)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 {
		return "Calibrating Suits..."
	}

	// Calculate panel sizes based on window width
	panelWidth := (m.width / 3) - 2
	panelHeight := m.height - 4

	// --- LEFT PANEL: VITALS ---
	vitalsContent := lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render("SYSTEM VITALS"),
		"",
		m.renderClockHUD(),
		"\n",
		m.renderStatusBadges(),
		"\n",
		lipgloss.NewStyle().Foreground(cCyan).Render("CPU INTEGRITY"),
		m.cpuBar.ViewAs(m.cpuVal),
		"\n",
		lipgloss.NewStyle().Foreground(cOrange).Render("THRUSTER POWER"),
		m.pwrBar.ViewAs(m.pwrVal),
		"\n",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("NETWORK STATUS"),
		m.netBar.ViewAs(m.netVal),
		"\n",
		m.renderCircularGauge(m.cpuVal, 15, "POWER LEVEL"),
		"\n",
		lipgloss.NewStyle().Foreground(cDim).Render("Mark LXXXV // Online"),
	)

	if m.alertActive {
		vitalsContent = lipgloss.JoinVertical(lipgloss.Left, vitalsContent, "\n", m.renderAlert())
	}

	leftPanel := boxStyle.Width(panelWidth).Height(panelHeight).Render(vitalsContent)

	// --- CENTER PANEL: ARC REACTOR ---
	theme := m.getTheme()

	// Visualizer rendering
	var resonanceView string
	bars := []string{" ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	for _, v := range m.resonance {
		idx := int(v * float64(len(bars)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(bars) {
			idx = len(bars) - 1
		}
		resonanceView += bars[idx]
	}

	centerTop := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(panelWidth/2).Align(lipgloss.Center, lipgloss.Center).Render(
			lipgloss.JoinVertical(lipgloss.Center,
				m.renderEnhancedArcReactor(),
				"\n",
				lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).Render("ARC REACTOR"),
				lipgloss.NewStyle().Foreground(theme.Dim).Render("Output: 4.8 GJ/s"),
			),
		),
		lipgloss.NewStyle().Width(panelWidth/2).Align(lipgloss.Center, lipgloss.Center).Render(
			lipgloss.JoinVertical(lipgloss.Center,
				lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("TARGETING"),
				"\n",
				m.renderRadar(),
			),
		),
	)

	centerContent := lipgloss.JoinVertical(lipgloss.Left,
		centerTop,
		"\n",
		m.renderSoundWave(),
		"\n",
		lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("NEURAL LINK"),
		matrixView(m),
		"\n",
		lipgloss.NewStyle().Foreground(pulsePurple).Bold(true).Render("DATA STREAM"),
		m.renderDataStream(),
	)

	if m.glitchActive {
		centerContent = glitchStyle.Render(centerContent)
	}

	centerPanel := boxStyle.Width(panelWidth).Height(panelHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(centerContent)

	// --- RIGHT PANEL: LOGS ---
	logHeader := headerStyle.Render("TELEMETRY STREAM")
	rightContent := lipgloss.JoinVertical(lipgloss.Left,
		logHeader,
		m.viewport.View(),
		"\n",
		lipgloss.NewStyle().Foreground(gridColor).Faint(true).Bold(true).Render("HOLOGRAPHIC FEED"),
		m.renderHologramGrid(4),
	)
	rightPanel := boxStyle.Width(panelWidth).Height(panelHeight).
		Render(rightContent)

	// Combine Columns
	ui := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, centerPanel, rightPanel)

	// Add Master Header with theme
	title := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Foreground(theme.Primary).
		Background(theme.Background).
		Render("/// STARK INDUSTRIES INTERFACE - " + theme.Name + " ///")

	baseView := lipgloss.JoinVertical(lipgloss.Top, title, ui)

	// Show help menu overlay if active
	if m.showHelp {
		helpMenu := m.renderHelpMenu()
		// Center the help menu
		helpOverlay := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, helpMenu)
		return helpOverlay
	}

	return baseView
}

// --- HUD Helper Functions ---

func (m model) renderClockHUD() string {
	t := time.Now()
	timeStr := t.Format("15:04:05")
	dateStr := t.Format("Jan 02 2006")

	clockSection := lipgloss.JoinVertical(lipgloss.Left,
		clockStyle.Render("╭─ SYSTEM TIME ─╮"),
		clockStyle.Render("│ "+timeStr+" │"),
		clockStyle.Render("│ "+dateStr+" │"),
		clockStyle.Render("╰──────────────╯"),
	)
	return clockSection
}

func (m model) renderStatusBadges() string {
	statusColor := badgeGreen
	if m.cpuVal > 0.8 {
		statusColor = badgeYellow
	}
	if m.cpuVal > 0.9 {
		statusColor = badgeRed
	}

	networkBadge := badgeGreen
	if m.netVal < 0.3 {
		networkBadge = badgeYellow
	}

	modeBadge := modeStyle
	switch m.currentMode {
	case "COMBAT":
		modeBadge = badgeRed
	case "STEALTH":
		modeBadge = badgePurple
	}

	badges := lipgloss.JoinHorizontal(lipgloss.Top,
		statusColor.Render("◉ SYS"),
		networkBadge.Render("◉ NET"),
		modeBadge.Render(m.currentMode),
	)
	return badges
}

func (m model) renderCircularGauge(value float64, size int, label string) string {
	chars := []string{" ", "◔", "◑", "◕", "●"}
	idx := int(value * float64(len(chars)-1))
	if idx < 0 {
		idx = 0
	}
	if idx >= len(chars) {
		idx = len(chars) - 1
	}

	fillChars := ""
	for i := 0; i < size; i++ {
		if i < int(value*float64(size)) {
			fillChars += "█"
		} else {
			fillChars += "░"
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(cCyan).Bold(true).Render(label),
		lipgloss.NewStyle().Foreground(cCyan).Render("["+fillChars+"] "+fmt.Sprintf("%d%%", int(value*100))),
	)
}

func (m model) renderRadar() string {
	blipChar := "●"

	radarStr := lipgloss.NewStyle().Foreground(alertGreen).Bold(true).Render(
		"     ◉\n   ╱ | ╲\n  " + blipChar + "--R--" + blipChar + "\n   ╲ | ╱\n     ◉",
	)
	return radarStr
}

func (m model) renderAlert() string {
	if !m.alertActive {
		return ""
	}

	severity := ""
	alertColor := alertStyle

	switch m.alertSeverity {
	case 1:
		severity = "WARNING"
	case 2:
		severity = "CAUTION"
	case 3:
		severity = "CRITICAL"
	}

	if m.tickCount%10 < 5 {
		alertColor = alertStyle.Background(lipgloss.Color("#660000"))
	}

	alertBox := lipgloss.NewStyle().
		Width(30).
		Align(lipgloss.Center).
		Border(lipgloss.ThickBorder()).
		BorderForeground(alertRed).
		Padding(0, 1).
		Render(alertColor.Render("⚠ " + severity + "\n" + m.alertMessage))

	return alertBox
}

func (m model) renderScanlineOverlay(width, height int) string {
	var sb strings.Builder
	lineChar := "─"

	for i := 0; i < height; i++ {
		switch i {
		case m.scanlinePos, (m.scanlinePos + 5) % height:
			sb.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF00")).
				Faint(true).
				Render(strings.Repeat(lineChar, width)))
		default:
			sb.WriteString("")
		}
		if i < height-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func (m model) renderDataStream() string {
	stream := ""
	for i := 0; i < 3; i++ {
		for j := 0; j < 8; j++ {
			idx := (m.tickCount + i*8 + j) % len(m.dataStreamChars)
			stream += m.dataStreamChars[idx]
		}
		stream += "\n"
	}
	return lipgloss.NewStyle().Foreground(nordTeal).Faint(true).Render(stream)
}

func (m model) renderHologramGrid(rows int) string {
	var sb strings.Builder
	for i := 0; i < rows; i++ {
		line := ""
		for j := 0; j < 30; j++ {
			if (i+j)%2 == 0 {
				line += m.hologramChars[(m.tickCount+i)%len(m.hologramChars)]
			} else {
				line += " "
			}
		}
		sb.WriteString(lipgloss.NewStyle().Foreground(gridColor).Faint(true).Render(line))
		if i < rows-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func (m model) renderHelpMenu() string {
	theme := m.getTheme()

	helpStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theme.Primary).
		Padding(1, 2).
		Background(theme.Background).
		Foreground(theme.Primary)

	titleStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true).
		Align(lipgloss.Center)

	keyStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(theme.Dim)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("╔═══════════════════════════════════╗"),
		titleStyle.Render("║    J.A.R.V.I.S. CONTROLS HELP    ║"),
		titleStyle.Render("╚═══════════════════════════════════╝"),
		"",
		keyStyle.Render("  q / Ctrl+C")+" "+descStyle.Render("│ Exit Application"),
		keyStyle.Render("  h / ?      ")+" "+descStyle.Render("│ Toggle This Help"),
		keyStyle.Render("  t          ")+" "+descStyle.Render("│ Cycle Themes"),
		keyStyle.Render("  p          ")+" "+descStyle.Render("│ Pause/Resume"),
		keyStyle.Render("  s          ")+" "+descStyle.Render("│ Toggle Sound Wave"),
		keyStyle.Render("  r          ")+" "+descStyle.Render("│ Reboot System"),
		keyStyle.Render("  Space      ")+" "+descStyle.Render("│ Manual Scan"),
		keyStyle.Render("  ↑ / ↓      ")+" "+descStyle.Render("│ Scroll Logs"),
		"",
		titleStyle.Render("Current Theme: "+theme.Name),
	)

	return helpStyle.Render(content)
}

func (m model) renderSoundWave() string {
	if !m.showSoundWave {
		return ""
	}

	theme := m.getTheme()
	var sb strings.Builder

	// Render frequency bars
	bars := []string{" ", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

	for _, level := range m.audioLevels {
		idx := int(level * float64(len(bars)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(bars) {
			idx = len(bars) - 1
		}
		sb.WriteString(bars[idx])
	}

	waveStr := sb.String()

	header := lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("AUDIO ANALYSIS")
	wave := lipgloss.NewStyle().Foreground(theme.Primary).Render(waveStr)
	labels := lipgloss.NewStyle().Foreground(theme.Dim).Faint(true).Render("Bass  Mid  High")

	return lipgloss.JoinVertical(lipgloss.Left, header, wave, labels)
}

func (m model) renderEnhancedArcReactor() string {
	theme := m.getTheme()

	// Pulsing effect using sine wave
	pulseIntensity := (math.Sin(m.arcReactorPhase) + 1) / 2

	// Create concentric circles with pulsing effect
	var reactor strings.Builder

	if pulseIntensity > 0.7 {
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("       ╔═══════╗\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Primary).Bold(true).Render("     ╔═╝ ◉ ◉ ◉ ╚═╗\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Primary).Render("    ║ ◉ ▓▓▓▓▓▓▓ ◉ ║\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("    ║ ◉ ▓█████▓ ◉ ║\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Primary).Render("    ║ ◉ ▓▓▓▓▓▓▓ ◉ ║\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Primary).Bold(true).Render("     ╚═╗ ◉ ◉ ◉ ╔═╝\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Accent).Bold(true).Render("       ╚═══════╝"))
	} else {
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Primary).Render("       ╔═══════╗\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Primary).Render("     ╔═╝ ◉ ◉ ◉ ╚═╗\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Dim).Render("    ║ ◉ ▓▓▓▓▓▓▓ ◉ ║\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Primary).Render("    ║ ◉ ▓█████▓ ◉ ║\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Dim).Render("    ║ ◉ ▓▓▓▓▓▓▓ ◉ ║\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Primary).Render("     ╚═╗ ◉ ◉ ◉ ╔═╝\n"))
		reactor.WriteString(lipgloss.NewStyle().Foreground(theme.Primary).Render("       ╚═══════╝"))
	}

	return reactor.String()
}

// --- Simulations ---

func tickCommand() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func generateLogCommand() tea.Cmd {
	// Random delay for log generation
	duration := time.Duration(rand.Intn(1000)+200) * time.Millisecond
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		opts := []string{
			"Repulsor calibration complete",
			"Targeting array locked",
			"Scanning spectral analysis",
			"Flight systems check: PASS",
			"Incoming transmission blocked",
			"Auxiliary power rerouted",
			"Nanite density: 98%",
			"Weather pattern analyzing...",
		}
		return logMsg(opts[rand.Intn(len(opts))])
	})
}

func matrixView(m model) string {
	var sb strings.Builder

	// Render grid
	for y := 0; y < m.matrixRows; y++ {
		for x := 0; x < m.matrixCols; x++ {
			if x >= len(m.matrixGrid) || y >= len(m.matrixGrid[x]) {
				sb.WriteString(" ")
				continue
			}

			char := m.matrixGrid[x][y]
			headY := m.matrixHeads[x]
			tailLen := m.matrixTails[x]

			// Determine Color
			var style lipgloss.Style

			if y == headY {
				// Head: Bright White/Teal
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
			} else if y < headY && y > headY-tailLen {
				// Trail: Fade from Teal to Green to Dark
				dist := headY - y
				// Simple 3-step gradient
				if dist < tailLen/3 {
					style = lipgloss.NewStyle().Foreground(nordTeal)
				} else if dist < (tailLen*2)/3 {
					style = lipgloss.NewStyle().Foreground(nordGreen)
				} else {
					style = lipgloss.NewStyle().Foreground(cDim) // Fading out
				}
				// Determine boldness/faintness
				if dist > tailLen/2 {
					style = style.Faint(true)
				}

			} else {
				// Off-trail (invisible/dim background noise? or just empty)
				// True Matrix is empty black sans trail
				sb.WriteString(" ")
				continue
			}

			// Adjust for Rune Width (Katakana is nice but standard terminal grid is easiest with space)
			// But runes variable width. Let's just print.
			// Force 1 cell width?

			// Hack: Runewidth. Or just add a space after?
			// Katakana is often half-width in modern terms but might render wider.
			// Let's stick to simple rendering.
			sb.WriteString(style.Render(string(char)))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting J.A.R.V.I.S.:", err)
	}
}
