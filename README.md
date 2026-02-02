```
     ██╗ █████╗ ██████╗ ██╗   ██╗██╗███████╗
     ██║██╔══██╗██╔══██╗██║   ██║██║██╔════╝
     ██║███████║██████╔╝██║   ██║██║███████╗
██   ██║██╔══██║██╔══██╗╚██╗ ██╔╝██║╚════██║
╚█████╔╝██║  ██║██║  ██║ ╚████╔╝ ██║███████║
 ╚════╝ ╚═╝  ╚═╝╚═╝  ╚═╝  ╚═══╝  ╚═╝╚══════╝
```

<div align="center">

### 🤖 **J**ust **A** **R**ather **V**ery **I**ntelligent **S**ystem

*A Stark Industries-inspired Terminal User Interface*

[![Go Version](https://img.shields.io/badge/Go-1.25.5-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![Bubble Tea](https://img.shields.io/badge/Bubble_Tea-1.3.10-FF69B4?style=for-the-badge)](https://github.com/charmbracelet/bubbletea)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

</div>

---

## 🌟 **Overview**

**JARVIS-GO** is a stunning, cyberpunk-themed terminal user interface that brings the aesthetic of Tony Stark's AI assistant to your command line. Built with Go and the powerful Bubble Tea framework, it features real-time system monitoring, Matrix-style animations, and a sleek Iron Man-inspired design.

```
╔═══════════════════════════════════════════════════════════════╗
║                  /// STARK INDUSTRIES ///                     ║
║                                                               ║
║  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          ║
║  │   SYSTEM    │  │     ARC     │  │  TELEMETRY  │          ║
║  │   VITALS    │  │   REACTOR   │  │   STREAM    │          ║
║  │             │  │             │  │             │          ║
║  │  CPU ████░  │  │      ◉      │  │ >> Repulsor │          ║
║  │  PWR ██████ │  │   ▂▄▆█▆▄▂   │  │ >> Scanning │          ║
║  │  NET ████░░ │  │  ｱｲｳｴｵｶｷｸ  │  │ >> Systems  │          ║
║  └─────────────┘  └─────────────┘  └─────────────┘          ║
╚═══════════════════════════════════════════════════════════════╝
```

---

## ✨ **Features**

### 🎨 **Visual Excellence**
- **Cyberpunk Aesthetic** — Stark Industries color palette with cyan, orange, and blue gradients
- **Matrix Rain Effect** — Cascading Katakana characters with dynamic trails and glitch effects
- **Arc Reactor Animation** — Spinning globe visualization with real-time resonance bars
- **Responsive Layout** — Three-column interface that adapts to terminal size

### 📊 **Real-Time Monitoring**
- **CPU Integrity** — Live CPU usage visualization with gradient progress bars
- **Thruster Power** — Power level monitoring with orange-to-red gradients
- **Network Status** — Network activity tracking with green-to-blue gradients
- **Telemetry Stream** — Scrolling log viewport with system events

### 🎭 **Interactive Elements**
- **Smooth Animations** — 60 FPS updates with Bubble Tea's event loop
- **Dynamic Data** — Simulated live metrics that fluctuate realistically
- **Auto-scrolling Logs** — Continuous stream of system messages
- **Graceful Resizing** — Automatic layout adjustment on terminal resize

---

## 🚀 **Quick Start**

### **Prerequisites**
- Go 1.25.5 or higher
- A terminal with Unicode support (for best visual experience)

### **Installation**

```bash
# Clone the repository
git clone <your-repo-url>
cd jarvis-go

# Install dependencies
go mod download

# Build the executable
go build -o jarvis main.go

# Run JARVIS
./jarvis
```

### **Make it Global** (Optional)

To run `jarvis` from anywhere:

```bash
# Move to a directory in your PATH
sudo mv jarvis /usr/local/bin/

# Now run from anywhere
jarvis
```

---

## 🎮 **Usage**

```bash
# Start the interface
./jarvis

# Exit the interface
Press 'q' or 'Ctrl+C'
```

### **Controls**
| Key | Action |
|-----|--------|
| `q` | Quit the application |
| `Ctrl+C` | Force quit |

---

## 🏗️ **Architecture**

```
┌─────────────────────────────────────────────────────────┐
│                      JARVIS-GO                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   Bubble Tea │  │   Lipgloss   │  │   Bubbles    │ │
│  │  (Framework) │  │   (Styling)  │  │ (Components) │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│         │                  │                  │        │
│         └──────────────────┴──────────────────┘        │
│                            │                           │
│  ┌─────────────────────────▼────────────────────────┐  │
│  │              Model (State Management)            │  │
│  │  • Spinner  • Progress Bars  • Viewport         │  │
│  │  • Matrix Grid  • Logs  • Metrics               │  │
│  └──────────────────────────────────────────────────┘  │
│                            │                           │
│  ┌─────────────────────────▼────────────────────────┐  │
│  │           Update Loop (Event Handler)            │  │
│  │  • Tick Messages  • Log Messages                │  │
│  │  • Window Resize  • Spinner Updates             │  │
│  └──────────────────────────────────────────────────┘  │
│                            │                           │
│  ┌─────────────────────────▼────────────────────────┐  │
│  │              View (Rendering)                    │  │
│  │  • Left Panel   • Center Panel   • Right Panel  │  │
│  │  • Matrix Renderer  • Layout Composition        │  │
│  └──────────────────────────────────────────────────┘  │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## 🎨 **Color Palette**

The interface uses a carefully curated color scheme inspired by Iron Man's suit:

| Color | Hex | Usage |
|-------|-----|-------|
| **Cyan** | `#00F0FF` | Primary accent, borders, headers |
| **Blue** | `#0077BE` | CPU metrics, gradients |
| **Orange** | `#FF5F1F` | Power metrics, log labels |
| **Dark** | `#1A1A1A` | Background |
| **Dim** | `#444444` | Secondary text |
| **Nord Green** | `#A3BE8C` | Matrix trail fade |
| **Nord Teal** | `#8FBCBB` | Matrix characters |

---

## 📦 **Dependencies**

```go
require (
    github.com/charmbracelet/bubbles v0.21.0      // UI components
    github.com/charmbracelet/bubbletea v1.3.10    // TUI framework
    github.com/charmbracelet/lipgloss v1.1.0      // Styling engine
)
```

### **Key Libraries**
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** — The Elm Architecture for Go TUIs
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** — Style definitions and layout
- **[Bubbles](https://github.com/charmbracelet/bubbles)** — Pre-built TUI components (spinner, progress, viewport)

---

## 🔧 **Customization**

### **Modify Colors**
Edit the color constants in `main.go`:

```go
var (
    cCyan   = lipgloss.Color("#00F0FF")  // Change primary accent
    cOrange = lipgloss.Color("#FF5F1F")  // Change power color
    // ... etc
)
```

### **Adjust Update Speed**
Change the tick interval in the `tickCommand()` function:

```go
func tickCommand() tea.Cmd {
    return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}
```

### **Customize Log Messages**
Edit the log options in `generateLogCommand()`:

```go
opts := []string{
    "Your custom message here",
    "Another system event",
    // ... add more
}
```

---

## 🎯 **Technical Highlights**

### **Matrix Rain Algorithm**
- Each column has independent head position, tail length, and speed
- Characters randomly glitch and change (1% per tick)
- Gradient coloring from white (head) → teal → green → dim (tail)
- Trails reset and restart when off-screen

### **Responsive Design**
- Dynamic panel sizing based on terminal width
- Automatic viewport height adjustment
- Grid reinitialization on window resize
- Maintains aspect ratio across different terminal sizes

### **Performance Optimization**
- Efficient string building with `strings.Builder`
- Bounded log buffer (max 50 entries)
- Probabilistic matrix updates to reduce CPU load
- Minimal allocations in hot paths

---

## 📸 **Screenshots**

> **Note:** For best results, use a terminal with:
> - Unicode/UTF-8 support
> - True color (24-bit) support
> - Monospaced font (e.g., Fira Code, JetBrains Mono)

---

## 🤝 **Contributing**

Contributions are welcome! Feel free to:
- Report bugs
- Suggest new features
- Submit pull requests
- Improve documentation

---

## 📄 **License**

This project is open source and available under the [MIT License](LICENSE).

---

## 🙏 **Acknowledgments**

- **Marvel/Iron Man** — For the incredible design inspiration
- **[Charm](https://charm.sh/)** — For the amazing Bubble Tea ecosystem
- **The Matrix** — For the iconic digital rain effect

---

<div align="center">

```
    ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄               ▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄ 
   ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░▌             ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
    ▀▀▀▀█░█▀▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀█░▌ ▐░▌           ▐░▌  ▀▀▀▀█░█▀▀▀▀ ▐░█▀▀▀▀▀▀▀▀▀ 
        ▐░▌     ▐░▌       ▐░▌▐░▌       ▐░▌  ▐░▌         ▐░▌       ▐░▌     ▐░▌          
        ▐░▌     ▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄█░▌   ▐░▌       ▐░▌        ▐░▌     ▐░█▄▄▄▄▄▄▄▄▄ 
        ▐░▌     ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌    ▐░▌     ▐░▌         ▐░▌     ▐░░░░░░░░░░░▌
        ▐░▌     ▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀█░█▀▀      ▐░▌   ▐░▌          ▐░▌      ▀▀▀▀▀▀▀▀▀█░▌
        ▐░▌     ▐░▌       ▐░▌▐░▌     ▐░▌       ▐░▌ ▐░▌           ▐░▌               ▐░▌
        ▐░▌     ▐░▌       ▐░▌▐░▌      ▐░▌       ▐░▐░▌            ▐░▌      ▄▄▄▄▄▄▄▄▄█░▌
        ▐░▌     ▐░▌       ▐░▌▐░▌       ▐░▌       ▐░▌             ▐░▌     ▐░░░░░░░░░░░▌
         ▀       ▀         ▀  ▀         ▀         ▀               ▀       ▀▀▀▀▀▀▀▀▀▀▀ 
```

**Made with ❤️ and Go**

*"Sometimes you gotta run before you can walk."* — Tony Stark

</div>
