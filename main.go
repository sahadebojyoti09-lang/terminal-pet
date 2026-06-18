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
	
	bubbleStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("242")).Width(42).Padding(0, 1)
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
	energy        int       
	isSleeping    bool      
	isAngry       bool      
	isDancing     bool 
	isBathing     bool 
	angryTimer    int       
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
		energy:        80, 
		isSleeping:    false,
		isAngry:       false,
		isDancing:     false,
		isBathing:     false,
		angryTimer:    0,
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

		case "e": 
			if m.isSleeping {
				m.isSleeping = false
				m.isAngry = true
				m.isDancing = false
				m.isBathing = false
				m.angryTimer = 16 
				m.happiness = max(0, m.happiness-25)
				m.speech = "💢 HEY! Why did you wake me up?! (╬◣_◢)"
				m.actionTimer = 12
			} else {
				m.isSleeping = true
				m.isAngry = false
				m.isInteracting = false
				m.isDancing = false
				m.isBathing = false
				m.speech = "💤 Goodnight... going to sleep mode."
				m.actionTimer = 8
			}

		case "f": 
			if m.isSleeping {
				m.speech = "zzz (Puchi is fast asleep, don't shove food...)"
				m.actionTimer = 8
				return m, nil
			}
			m.isDancing = false
			m.isBathing = false
			m.hunger = max(0, m.hunger-20)
			m.happiness = min(100, m.happiness+5)
			m.speech = "✨ Nom nom! *crunchy noises* ✨"
			m.actionTimer = 8 
			m.isInteracting = true 

		case "p": 
			if m.isSleeping {
				m.speech = "zzz (Puchi twitches his ears but stays asleep...)"
				m.actionTimer = 8
				return m, nil
			}
			m.isDancing = false
			m.isBathing = false
			m.happiness = min(100, m.happiness+15)
			m.speech = "❤️ Puchi purrs like a well-optimized system! ❤️"
			m.actionTimer = 8 
			m.isInteracting = true 

		case "d": 
			if m.isSleeping {
				m.speech = "zzz (Puchi is sleeping deep, he can't hear the beat...)"
				m.actionTimer = 8
				return m, nil
			}
			if m.energy < 20 {
				m.speech = "❌ (Puchi is way too tired to dance right now!)"
				m.actionTimer = 8
				return m, nil
			}
			m.isDancing = true
			m.isBathing = false
			m.isInteracting = false
			m.happiness = min(100, m.happiness+10)
			m.energy = max(0, m.energy-8) 
			m.speech = "🕺 *Bass Drops* Check out these layout moves! 🕺"
			m.actionTimer = 16 

		case "b": 
			if m.isSleeping {
				m.speech = "zzz (Don't splash water on a sleeping pet!)"
				m.actionTimer = 8
				return m, nil
			}
			m.isBathing = true
			m.isDancing = false
			m.isInteracting = false
			m.happiness = min(100, m.happiness+10)
			m.energy = max(0, m.energy-5)
			m.speech = "🫧 *Scrub Scrub* Squeaky clean compilation! 🫧"
			m.actionTimer = 12 

		case "s": 
			if m.isSleeping {
				m.speech = "zzz (Puchi snores softly...)"
				m.actionTimer = 8
				return m, nil
			}
			if m.isAngry {
				m.speech = "❌ (Puchi is too mad to talk right now!)"
				m.actionTimer = 10
				return m, nil
			}
			if m.hunger > 70 {
				m.speech = "❌ (Puchi grumbles... too hungry to talk!)"
				m.actionTimer = 10 
				m.isInteracting = false
			} else {
				m.isDancing = false
				m.isBathing = false
				m.speech = getFortune()
				m.actionTimer = -1 
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

		if m.isAngry && m.angryTimer > 0 {
			m.angryTimer--
			if m.angryTimer == 0 {
				m.isAngry = false
			}
		}

		if m.actionTimer > 0 {
			m.actionTimer--
			if m.actionTimer == 0 {
				m.speech = "..."
				m.isInteracting = false 
				m.isDancing = false
				m.isBathing = false
			}
		}
		return m, tick()

	case statusTickMsg:
		if m.isSleeping {
			m.energy = min(100, m.energy+10) 
			m.hunger = min(100, m.hunger+2)  
			if m.energy == 100 {
				m.isSleeping = false 
				m.speech = "🌅 (Puchi woke up refreshed on his own!)"
				m.actionTimer = 12
			}
		} else {
			if m.frame%2 == 0 {
				m.energy = max(0, m.energy-5)
			}
			m.hunger = min(100, m.hunger+3)
			m.happiness = max(0, m.happiness-2)

			if m.energy == 0 {
				m.isSleeping = true
				m.isInteracting = false
				m.isDancing = false
				m.isBathing = false
				m.speech = "⚠️ *Thud* (Puchi collapsed from pure exhaustion!)"
				m.actionTimer = 16
			}
		}
		return m, statusTick()
	}

	return m, nil
}

func (m model) View() string {
	var face string
	var emotionText string
	
	if m.isSleeping {
		emotionText = "Sleeping 💤"
		if m.frame%2 == 0 {
			face = "    (─‿─) zzZ \n  -|═  ═|- " 
		} else {
			face = "    (─‿─)  zZ \n  -|═  ═|- " 
		}
	} else if m.isAngry {
		emotionText = "Grumpy 💢"
		if m.frame%2 == 0 {
			face = "   (╬◣_◢)  \n  ~|═  ═|~ "
		} else {
			face = "   (╬◣_◢)  \n  -|═  ═|- "
		}
	} else if m.isDancing {
		emotionText = "Grooving 🕺"
		switch m.frame % 4 {
		case 0:
			face = " ♬  ヘ( ^‿^ )ノ \n           "
		case 1:
			face = "    ヽ( ^‿^ )ノ ♫\n           "
		case 2:
			face = " ♫  ＼( ^‿^ )ノ \n           "
		default:
			face = "    ノ( ^‿^ )ヘ ♬\n           "
		}
	} else if m.isBathing {
		emotionText = "Bathing 🧼"
		if m.frame%2 == 0 {
			face = "    o  ✧  °  \n   ( ˘ᵕ˘ ) 🫧 "
		} else {
			face = "   °  o  ✧   \n   ( ˘ᵕ˘ ) 🫧 "
		}
	} else if m.isInteracting {
		emotionText = "Happy ✨"
		if m.frame%2 == 0 {
			face = "  (っ◕‿◕)っ\n           " 
		} else {
			face = "   (っ◕‿◕)っ\n           " 
		}
	} else if m.hunger > 70 || m.happiness < 30 {
		emotionText = "Miserable 😭"
		face = "   (╥﹏╥)  \n  -|═  ═|- " 
	} else {
		emotionText = "Content 😐"
		if m.isBlinking {
			face = "   (-‿ -)  \n  -|═  ═|- " 
		} else if m.frame%2 == 0 {
			face = "   (^‿ ^)  \n  -|═  ═|- " 
		} else {
			face = "   (^‿ ^)  \n  ~|═  ═|~ " 
		}
	}

	hungerBar := fmt.Sprintf("Hunger:    [%-10s] %d%%", repeatChar("█", m.hunger/10), m.hunger)
	happyBar  := fmt.Sprintf("Happiness: [%-10s] %d%%", repeatChar("█", m.happiness/10), m.happiness)
	energyBar := fmt.Sprintf("Energy:    [%-10s] %d%%", repeatChar("█", m.energy/10), m.energy)
	emotionStr := fmt.Sprintf("Emotion:   %s", emotionText)

	wrappedSpeech := lipgloss.NewStyle().Width(40).Render(actionStyle.Render(m.speech))

	renderedBubble := bubbleStyle.Render(wrappedSpeech)
	petRender := petStyle.Render(face)
	statusRender := statusStyle.Render(fmt.Sprintf("%s\n%s\n%s\n\n%s", hungerBar, happyBar, energyBar, emotionStr))
	helpRender := statusStyle.Render("\n[f] Feed  •  [p] Pet  •  [d] Dance  •  [b] Bath  •  [s] Speak  •  [e] Sleep/Wake  •  [q] Quit")

	body := fmt.Sprintf("%s\n\n%s\n\n%s\n%s", renderedBubble, petRender, statusRender, helpRender)

	lockedFrame := frameStyle.Width(54).Render(body)

	return lockedFrame + "\n"
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
