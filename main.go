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
	matrixCols  int
	matrixRows  int
	matrixGrid  [][]rune // The character at each cell
	matrixHeads []int    // Y position of the 'head'
	matrixTails []int    // Length of the trail
	matrixSpeed []int    // Speed factor
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
		spinner:   s,
		cpuBar:    p1,
		pwrBar:    p2,
		netBar:    p3,
		viewport:  vp,
		logs:      []string{"Initializing J.A.R.V.I.S. Protocol..."},
		cpuVal:    0.2,
		pwrVal:    0.8,
		netVal:    0.5,
		resonance: make([]float64, 20), // Initial buffer
		// Matrix fields initialized by Update loop
		matrixCols: 0,
		matrixRows: 0,
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
		matrixView(m),
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
