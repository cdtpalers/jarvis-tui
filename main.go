package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	logLabel  = lipgloss.NewStyle().Foreground(cOrange).Bold(true)
	logText   = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	coreStyle = lipgloss.NewStyle().Foreground(cCyan)
)

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
	matrixCols  []int // Y position of the 'head'
	matrixSpeed []int // Speed factor
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
		spinner:     s,
		cpuBar:      p1,
		pwrBar:      p2,
		netBar:      p3,
		viewport:    vp,
		logs:        []string{"Initializing J.A.R.V.I.S. Protocol..."},
		cpuVal:      0.2,
		pwrVal:      0.8,
		netVal:      0.5,
		resonance:   make([]float64, 20), // Initial buffer
		matrixCols:  make([]int, 20),     // Initial placeholder size
		matrixSpeed: make([]int, 20),
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
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
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
		matrixWidth := colWidth - 2
		if len(m.matrixCols) != matrixWidth {
			newCols := make([]int, matrixWidth)
			newSpeed := make([]int, matrixWidth)
			for i := range newCols {
				newCols[i] = rand.Intn(m.height) * -1 // Start above view
				newSpeed[i] = rand.Intn(2) + 1
			}
			m.matrixCols = newCols
			m.matrixSpeed = newSpeed
		}

	case tickMsg:
		// Simulate live data changes
		m.cpuVal += (rand.Float64() - 0.5) * 0.1
		if m.cpuVal > 1 {
			m.cpuVal = 1
		}
		if m.cpuVal < 0 {
			m.cpuVal = 0
		}

		m.pwrVal += (rand.Float64() - 0.5) * 0.05
		if m.pwrVal > 1 {
			m.pwrVal = 1
		}
		if m.pwrVal < 0 {
			m.pwrVal = 0
		}

		m.netVal += (rand.Float64() - 0.5) * 0.08
		if m.netVal > 1 {
			m.netVal = 1
		}
		if m.netVal < 0 {
			m.netVal = 0
		}

		// Update Resonance
		// Shift left
		if len(m.resonance) > 0 {
			m.resonance = m.resonance[1:]
			m.resonance = append(m.resonance, rand.Float64())
		}

		// Update Matrix
		for i := range m.matrixCols {
			m.matrixCols[i] += m.matrixSpeed[i]
			if m.matrixCols[i] > m.height {
				m.matrixCols[i] = rand.Intn(10) * -1 // Reset to top with random delay
				m.matrixSpeed[i] = rand.Intn(2) + 1
			}
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
	vitals := lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render("SYSTEM VITALS"),
		"",
		lipgloss.NewStyle().Foreground(cCyan).Render("CPU INTEGRITY"),
		m.cpuBar.ViewAs(m.cpuVal),
		"\n",
		lipgloss.NewStyle().Foreground(cOrange).Render("THRUSTER POWER"),
		m.pwrBar.ViewAs(m.pwrVal),
		"\n",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("NETWORK STATUS"),
		m.netBar.ViewAs(m.netVal),
		"\n",
		lipgloss.NewStyle().Foreground(cDim).Render("Mark LXXXV // Online"),
	)
	leftPanel := boxStyle.Width(panelWidth).Height(panelHeight).Render(vitals)

	// --- CENTER PANEL: ARC REACTOR ---
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

	coreContent := lipgloss.JoinVertical(lipgloss.Center,
		m.spinner.View(),
		"\n",
		lipgloss.NewStyle().Bold(true).Foreground(cCyan).Render("ARC REACTOR"),
		lipgloss.NewStyle().Foreground(cDim).Render("Output: 4.8 GJ/s"),
		"\n",
		lipgloss.NewStyle().Foreground(nordTeal).Bold(true).Render("NEURAL LINK"),
		matrixView(m.matrixCols, m.matrixSpeed, panelWidth, panelHeight/2),
	)
	centerPanel := boxStyle.Width(panelWidth).Height(panelHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(coreContent)

	// --- RIGHT PANEL: LOGS ---
	logHeader := headerStyle.Render("TELEMETRY STREAM")
	rightPanel := boxStyle.Width(panelWidth).Height(panelHeight).
		Render(lipgloss.JoinVertical(lipgloss.Left, logHeader, m.viewport.View()))

	// Combine Columns
	ui := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, centerPanel, rightPanel)

	// Add Master Header
	title := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Foreground(cCyan).
		Background(cDark).
		Render("/// STARK INDUSTRIES INTERFACE ///")

	return lipgloss.JoinVertical(lipgloss.Top, title, ui)
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

func matrixView(cols, speeds []int, width, height int) string {
	var sb strings.Builder

	// Create visual grid
	for y := 0; y < height; y++ {
		for x := 0; x < width && x < len(cols); x++ {
			headY := cols[x]
			char := " "

			if y == headY {
				// Head char
				char = "█" // or a random digit
				sb.WriteString(lipgloss.NewStyle().Foreground(nordTeal).Render(char))
			} else if y < headY && y > headY-5 {
				// Trail
				char = fmt.Sprintf("%d", rand.Intn(10))
				sb.WriteString(lipgloss.NewStyle().Foreground(nordGreen).Faint(true).Render(char))
			} else {
				sb.WriteString(" ")
			}
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
