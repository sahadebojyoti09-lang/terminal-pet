package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	petStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true) // Orange Accent
	statusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))        // Grey text
	actionStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)  // Purple Accent
	bubbleStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("242")).Width(40).Padding(0, 1)
	frameStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("240"))
)

type tickMsg time.Time
type statusTickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func statusTick() tea.Cmd {
	return tea.Tick(time.Second*3, func(t time.Time) tea.Msg { return statusTickMsg(t) })
}

func getFortune() string {
	cmd := exec.Command("fortune", "-s")
	out, err := cmd.Output()
	if err != nil {
		return "Puchi is staring blankly into the terminal..."
	}
	return strings.TrimSpace(string(out))
}

type model struct {
	frame         int
	hunger        int
	happiness     int
	speech        string
	actionTimer   int
	isBlinking    bool
	isInteracting bool 
}

func initialModel() model {
	return model{
		frame:         0,
		hunger:        30,
		happiness:     70,
		speech:        "Mwahaha! I am alive inside your terminal!",
		actionTimer:   12, 
		isBlinking:    false,
		isInteracting: false,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), statusTick())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "f": // Feed
			m.hunger = max(0, m.hunger-20)
			m.happiness = min(100, m.happiness+5)
			m.speech = "✨ Nom nom! *crunchy noises* ✨"
			m.actionTimer = 8 
			m.isInteracting = true 
		case "p": // Pet
			m.happiness = min(100, m.happiness+15)
			m.speech = "❤️ Puchi purrs like a well-optimized system! ❤️"
			m.actionTimer = 8 
			m.isInteracting = true 
		case "s": // Speak
			if m.hunger > 70 {
				m.speech = "❌ (Puchi grumbles... too hungry to talk!)"
				m.actionTimer = 10 // Clears after 2.5 seconds so his sad face can show!
				m.isInteracting = false
			} else {
				m.speech = getFortune()
				m.actionTimer = -1 // Freezes the quote safely for reading
				m.isInteracting = false
			}
		}

	case tickMsg:
		m.frame++
		if m.frame%12 == 0 {
			m.isBlinking = true
		} else if m.frame%12 == 1 {
			m.isBlinking = false
		}

		if m.actionTimer > 0 {
			m.actionTimer--
			if m.actionTimer == 0 {
				m.speech = "..."
				m.isInteracting = false 
			}
		}
		return m, tick()

	case statusTickMsg:
		m.hunger = min(100, m.hunger+3)
		m.happiness = max(0, m.happiness-2)
		return m, statusTick()
	}

	return m, nil
}

func (m model) View() string {
	var face string
	
	// Layout logic prioritization
	if m.isInteracting {
		// Beautiful clean custom happy bounce dance (No tricky backslashes!)
		if m.frame%2 == 0 {
			face = "  (っ◕‿◕)っ \n   |═❤️═|  " 
		} else {
			face = "  (っ◕‿◕)っ \n  ~|═❤️═|~ " 
		}
	} else if m.hunger > 70 {
		face = "  (╥﹏╥)  \n  -|═  ═|- " 
	} else if m.isBlinking {
		face = "  (-‿ -)  \n  -|═  ═|- " 
	} else if m.frame%2 == 0 {
		face = "  (^‿ ^)  \n  -|═  ═|- " 
	} else {
		face = "  (^‿ ^)  \n  ~|═  ═|~ " 
	}

	hungerBar := fmt.Sprintf("Hunger:    [%-10s] %d%%", repeatChar("█", m.hunger/10), m.hunger)
	happyBar := fmt.Sprintf("Happiness: [%-10s] %d%%", repeatChar("█", m.happiness/10), m.happiness)

	renderedBubble := bubbleStyle.Render(actionStyle.Render(m.speech))
	petRender := petStyle.Render(face)
	statusRender := statusStyle.Render(fmt.Sprintf("%s\n%s", hungerBar, happyBar))
	helpRender := statusStyle.Render("\n[f] Feed  •  [p] Pet  •  [s] Speak  •  [q] Quit")

	body := fmt.Sprintf("%s\n\n%s\n\n%s\n%s", renderedBubble, petRender, statusRender, helpRender)

	return frameStyle.Render(body) + "\n"
}

func repeatChar(char string, count int) string {
	s := ""
	for i := 0; i < count; i++ { s += char }
	return s
}
func min(a, b int) int { if a < b { return a }; return b }
func max(a, b int) int { if a > b { return a }; return b }

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, an error occurred: %v", err)
		os.Exit(1)
	}
}
