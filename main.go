package main

import (
	"fmt"
	"math/rand"
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
		return "Puchi is contemplating the nature of source code..."
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
	isBrushing    bool 
	inStyleMode   bool 
	isLoafMode    bool 
	isPlayingBall bool 
	isUsingToilet bool
	hairIndex     int  
	clothIndex    int  
	angryTimer    int       
	speech        string
	actionTimer   int
	isBlinking    bool
	isInteracting bool 
	
	// Weight and Feeding Tracking Engine
	feedCount     int
	weight        int // Increases every 20 feeds
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
		isBrushing:    false,
		inStyleMode:   false,
		isLoafMode:    false,
		isPlayingBall: false,
		isUsingToilet: false,
		hairIndex:     0,
		clothIndex:    0,
		angryTimer:    0,
		speech:        "Mwahaha! I am alive inside your terminal! 🔮",
		actionTimer:   12, 
		isBlinking:    false,
		isInteracting: false,
		feedCount:     0,
		weight:        0,
	}
}

func (m model) Init() tea.Cmd {
	rand.Seed(time.Now().UnixNano())
	return tea.Batch(tick(), statusTick())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.feedCount >= 100 && msg.String() != "t" && msg.String() != "q" && msg.String() != "ctrl+c" {
			m.speech = "🛑 CRITICAL ERROR: System Full! Force toilet break needed [t]! 🚽"
			return m, nil
		}

		if m.inStyleMode {
			switch msg.String() {
			case "c", "esc":
				m.inStyleMode = false
				m.speech = "Exited wardrobe setup mode. 👕"
				m.actionTimer = 6
				return m, nil
			case "1": 
				if m.isSleeping { return m, nil }
				m.isBrushing = true
				m.isDancing = false
				m.isBathing = false
				m.isPlayingBall = false
				m.isUsingToilet = false
				m.happiness = min(100, m.happiness+10)
				m.speech = "✨ *Scrubba-scrubba* Keeping the bytecode clean! 🦷"
				m.actionTimer = 10
				return m, nil
			case "2": 
				if m.isSleeping { return m, nil }
				m.hairIndex = (m.hairIndex + 1) % 4
				hairs := []string{"Freshly shaved bald! 🥚", "Spiky Punk Fringe! ⚡", "Classy Top Hat! 🎩", "Minimalist Headband! 🎗️"}
				m.speech = "💇 " + hairs[m.hairIndex]
				m.actionTimer = 8
				return m, nil
			case "3": 
				if m.isSleeping { return m, nil }
				m.clothIndex = (m.clothIndex + 1) % 4
				clothes := []string{"Comfy birthday suit! 🐳", "Classic Utility Jacket! 🧥", "Sharp Tuxedo Cross-Tie! 👔", "Warm Grid Scarf! 🧣"}
				m.speech = "✨ " + clothes[m.clothIndex]
				m.actionTimer = 8
				return m, nil
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "l": 
			m.isLoafMode = !m.isLoafMode
			m.inStyleMode = false
			m.isDancing = false
			m.isBathing = false
			m.isBrushing = false
			m.isPlayingBall = false
			m.isUsingToilet = false
			if m.isLoafMode {
				m.speech = "🍞 Loaf Mode ON: Puchi has engaged autopilot."
			} else {
				m.speech = "🛑 Loaf Mode OFF: Manual overrides restored."
			}
			m.actionTimer = 12

		case "c": 
			if m.isLoafMode { m.speech = "(Disable Loaf Mode [l] to open manual style settings)"; return m, nil }
			if m.isSleeping { m.speech = "zzz (Wake him up before dressing him! 😴)"; return m, nil }
			m.inStyleMode = true
			m.speech = "🎨 Wardrobe Open! Select structural styles below."
			m.actionTimer = -1

		case "e": 
			if m.isLoafMode { m.speech = "(Disable Loaf Mode [l] to handle sleep schedules)"; return m, nil }
			if m.isSleeping {
				m.isSleeping = false
				m.isAngry = true
				m.angryTimer = 16 
				m.happiness = max(0, m.happiness-25)
				m.speech = "Academic routine offline! Why did you wake me up?! (╬◣_◢)"
				m.actionTimer = 12
			} else {
				m.isSleeping = true
				m.isAngry = false
				m.isInteracting = false
				m.isDancing = false
				m.isBathing = false
				m.isBrushing = false
				m.isUsingToilet = false
				m.speech = "💤 Goodnight... going to sleep mode."
				m.actionTimer = 8
			}

		case "f": 
			if m.isLoafMode { m.speech = "(Puchi will eat when he feels like it!)"; return m, nil }
			if m.isSleeping { m.speech = "zzz (Puchi is fast asleep, don't shove food... 💤)"; return m, nil }
			m.isDancing = false
			m.isBathing = false
			m.isBrushing = false
			m.isUsingToilet = false
			m.hunger = max(0, m.hunger-20)
			m.happiness = min(100, m.happiness+5)
			
			m.feedCount++
			m.weight = m.feedCount / 20

			m.speech = "✨ Nom nom! *crunchy noises* 🍖"
			m.actionTimer = 16 
			m.isInteracting = true 

		case "t": 
			if m.isSleeping { m.speech = "zzz (Wake him up first!)"; return m, nil }
			m.isDancing = false
			m.isBathing = false
			m.isBrushing = false
			m.isInteracting = false
			m.isPlayingBall = false
			m.isUsingToilet = true
			
			m.weight = 0
			m.feedCount = 0
			m.speech = "🚽 *Flushhhhh* System cache cleared! Back to baseline speed. ✨"
			m.actionTimer = 20

		case "p": 
			if m.isLoafMode { m.speech = "(Loaf Mode active: Puchi looks cozy and complete)"; return m, nil }
			if m.isSleeping { m.speech = "zzz (Puchi twitches his ears but stays asleep... 🐱)"; return m, nil }
			m.isDancing = false
			m.isBathing = false
			m.isBrushing = false
			m.isUsingToilet = false
			m.happiness = min(100, m.happiness+15)
			m.speech = "❤️ Puchi purrs like a well-optimized system! ❤️"
			m.actionTimer = 8 
			m.isInteracting = true 

		case "d": 
			if m.isLoafMode { m.speech = "(Puchi is in deep loaf format, no dancing allowed)"; return m, nil }
			if m.isSleeping { m.speech = "zzz (Puchi is sleeping deep, he can't hear the beat... 🎷)"; return m, nil }
			if m.energy < 20 { m.speech = "❌ (Puchi is way too tired to dance right now!)"; return m, nil }
			m.isDancing = true
			m.isBathing = false
			m.isBrushing = false
			m.isInteracting = false
			m.isUsingToilet = false
			m.happiness = min(100, m.happiness+10)
			m.energy = max(0, m.energy-8) 
			m.speech = "🕺 *Bass Drops* Check out these layout moves! 🕺"
			m.actionTimer = 16 

		case "b": 
			if m.isLoafMode { m.speech = "(Puchi washes behind his own ears in Loaf mode)"; return m, nil }
			if m.isSleeping { m.speech = "zzz (Don't splash water on a sleeping pet! 🌊)"; return m, nil }
			m.isBathing = true
			m.isDancing = false
			m.isBrushing = false
			m.isInteracting = false
			m.isUsingToilet = false
			m.happiness = min(100, m.happiness+10)
			m.energy = max(0, m.energy-5)
			m.speech = "🫧 *Scrub Scrub* Squeaky clean compilation! 🫧"
			m.actionTimer = 12 

		case "s": 
			if m.isSleeping { m.speech = "zzz (Puchi snores softly... 💭)"; return m, nil }
			if m.isAngry { m.speech = "❌ (Puchi is too mad to talk right now!)"; return m, nil }
			if m.hunger > 70 { m.speech = "❌ (Puchi grumbles... too hungry to talk!)"; return m, nil }
			
			m.isDancing = false
			m.isBathing = false
			m.isBrushing = false
			m.isPlayingBall = false
			m.isUsingToilet = false
			m.speech = "💬 " + getFortune()
			m.actionTimer = -1 
			m.isInteracting = false
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
				m.isBrushing = false
				m.isPlayingBall = false
				m.isUsingToilet = false
			}
		}
		return m, tick()

	case statusTickMsg:
		if m.feedCount >= 80 && m.feedCount < 100 && !m.isSleeping && !m.isLoafMode {
			m.speech = "⚠️ [System Warning]: High buffer pressure detected. Toilet break suggested! 🚽"
			m.actionTimer = 12
		}

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
				m.isBrushing = false
				m.isPlayingBall = false
				m.isUsingToilet = false
				m.inStyleMode = false
				m.speech = "⚠️ *Thud* (Puchi collapsed from pure exhaustion!)"
				m.actionTimer = 16
			}
		}

		// AUTOPILOT AI ENGINE (Now with proactive toilet logic)
		if m.isLoafMode && !m.isSleeping && m.energy > 0 {
			if m.actionTimer <= 0 {
				// 1. HIGH STORAGE EMERGENCIES: Dump data before hitting 100% panic threshold
				if m.feedCount >= 80 {
					m.isDancing = false
					m.isBathing = false
					m.isBrushing = false
					m.isInteracting = false
					m.isPlayingBall = false
					m.isUsingToilet = true
					
					m.weight = 0
					m.feedCount = 0
					m.speech = "🤖 [Autopilot] High buffer pressure! Flushing system cache... 🚽"
					m.actionTimer = 20
				} else if m.hunger > 65 && m.feedCount < 100 {
					m.hunger = max(0, m.hunger-25)
					m.feedCount++
					m.weight = m.feedCount / 20
					m.speech = "🤖 [Autopilot] Eating a snack by myself! 🍪"
					m.actionTimer = 24 
					m.isInteracting = true
				} else if m.energy < 25 {
					m.isSleeping = true
					m.speech = "🤖 [Autopilot] Taking a quick power nap... 💤"
					m.actionTimer = 32 
				} else if m.happiness < 40 {
					m.isBathing = true
					m.happiness = min(100, m.happiness+15)
					m.speech = "🤖 [Autopilot] Taking a self-care shower! 🧼"
					m.actionTimer = 24 
				} else if m.happiness < 60 && m.energy > 40 {
					m.isPlayingBall = true
					m.happiness = min(100, m.happiness+20)
					m.speech = "🤖 [Autopilot] Chasing a rubber ball! ⚽"
					m.actionTimer = 28 
				} else {
					roll := rand.Intn(20) 
					switch roll {
					case 0: 
						m.isBrushing = true
						m.hairIndex = rand.Intn(4)
						m.clothIndex = rand.Intn(4)
						m.speech = "🤖 [Autopilot] Redressing my system configs! ✨"
						m.actionTimer = 24
					case 1: 
						m.speech = "💬 " + getFortune()
						m.actionTimer = 32 
					default:
						m.speech = "..."
						m.isInteracting = false 
						m.isDancing = false
						m.isBathing = false
						m.isBrushing = false
						m.isPlayingBall = false
						m.isUsingToilet = false
					}
				}
			}
		}

		return m, statusTick()
	}

	return m, nil
}

func (m model) View() string {
	var hairLine string
	var faceLine string
	var armLine string
	var emotionText string

	isEatingAction := m.isInteracting && (strings.Contains(m.speech, "Nom nom") || strings.Contains(m.speech, "snack"))

	switch m.hairIndex {
	case 1: 
		hairLine = "|||||||"
	case 2: 
		hairLine = "  ___  \n |___|"
	case 3: 
		hairLine = "[=====]"
	default: 
		hairLine = ""
	}
	
	if m.isSleeping {
		emotionText = "Sleeping 💤"
		if m.frame%2 == 0 {
			faceLine = "(- ‿ -) zZ" 
		} else {
			faceLine = "(- ‿ -)" 
		}
	} else if m.isAngry {
		emotionText = "Grumpy 💢"
		faceLine = "(> ◣ <)"
	} else if m.isDancing {
		emotionText = "Grooving 🕺"
		switch m.frame % 4 {
		case 0:  faceLine = "/ ( ^ ‿ ^ )"
		case 1:  faceLine = "( ^ ‿ ^ )\\"
		case 2:  faceLine = "\\ ( ^ ‿ ^ )"
		default: faceLine = "( ^ ‿ ^ )/"
		}
	} else if m.isBathing {
		emotionText = "Bathing 🧼"
		if m.frame%2 == 0 {
			faceLine = "o . ( ˘ ᵕ ˘ ) ."
		} else {
			faceLine = ". o ( ˘ ᵕ ˘ )o."
		}
	} else if m.isBrushing {
		emotionText = "Brushing ✨"
		if m.frame%2 == 0 {
			faceLine = "* ( - ‿ - ) *"
		} else {
			faceLine = "*( - ‿ - )*"
		}
	} else if m.isPlayingBall {
		emotionText = "Playing ⚽"
		if m.frame%2 == 0 {
			faceLine = "(っ^ ‿ ^)っ ◯"
		} else {
			faceLine = "◯ (っ^ ‿ ^)っ"
		}
	} else if m.isUsingToilet {
		emotionText = "Relieved 🚽"
		faceLine = "🚽 [ ˘ ᵕ ˘ ]"
	} else if isEatingAction {
		emotionText = "Munching 🍪"
		if m.frame%2 == 0 {
			faceLine = "(っ^ ‿ ^)っ 🍪" 
		} else {
			faceLine = "(っ〇 ‿ 〇)っ 🍪" 
		}
	} else if m.isInteracting {
		emotionText = "Happy ✨"
		faceLine = "(っ. ‿ .)っ" 
	} else if m.isLoafMode {
		emotionText = "Loafing 🍞"
		if m.isBlinking {
			faceLine = "(- ‿ -)"
		} else {
			faceLine = "(っ- ‿ -)っ"
		}
	} else if m.feedCount >= 100 {
		emotionText = "PANIC 🚨"
		faceLine = "(╬〇 ◣ 〇)"
	} else if m.hunger > 70 || m.happiness < 30 {
		emotionText = "Miserable 😭"
		faceLine = "(; ﹏ ;)" 
	} else {
		emotionText = "Content 😐"
		if m.isBlinking {
			faceLine = "(- ‿ -)" 
		} else {
			faceLine = "(^ ‿ ^)" 
		}
	}

	var currentCloth string
	switch m.clothIndex {
	case 1:  currentCloth = "▓▓▓▓▓"
	case 2:  currentCloth = "[ X ]"
	case 3:  currentCloth = "#####"
	default: currentCloth = "====="
	}

	if m.isUsingToilet {
		armLine = "  [  🧻  ]  "
	} else if !isEatingAction {
		switch m.weight {
		case 0:  armLine = "---|" + currentCloth + "|---"             
		case 1:  armLine = "---|(" + currentCloth + ")|---"           
		case 2:  armLine = "---|( " + currentCloth + " )|---"         
		case 3:  armLine = "===(  " + currentCloth + "  )==="         
		default: armLine = "(((   " + currentCloth + "   )))"         
		}
	}

	if m.isDancing || m.isBathing || m.isBrushing || m.isPlayingBall || m.isInteracting || isEatingAction {
		armLine = ""
	} else if m.isLoafMode && !m.isUsingToilet {
		switch m.weight {
		case 0:  armLine = "[===========]"
		case 1:  armLine = "[(===========)]"
		case 2:  armLine = "[((===========))]"
		default: armLine = "(((============)))"
		}
	}

	faceWidth := lipgloss.Width(faceLine)
	var centeredHair string
	if hairLine != "" {
		centeredHair = lipgloss.NewStyle().Width(faceWidth).Align(lipgloss.Center).Render(hairLine)
	}

	var assembledPet string
	if centeredHair != "" {
		assembledPet = lipgloss.JoinVertical(lipgloss.Center, centeredHair, faceLine, armLine)
	} else {
		assembledPet = lipgloss.JoinVertical(lipgloss.Center, faceLine, armLine)
	}

	hungerBar := fmt.Sprintf("Hunger:    [%s] %d%%", renderBar(m.hunger), m.hunger)
	happyBar  := fmt.Sprintf("Happiness: [%s] %d%%", renderBar(m.happiness), m.happiness)
	energyBar := fmt.Sprintf("Energy:    [%s] %d%%", renderBar(m.energy), m.energy)
	
	weightText := fmt.Sprintf("%d kg", 4+m.weight*5)
	if m.weight >= 4 { weightText += " (🚨 Unit)" }
	
	emotionStr := fmt.Sprintf("Emotion:   %-12s Weight: %s", emotionText, weightText)

	wrappedSpeech := lipgloss.NewStyle().Width(40).Render(actionStyle.Render(m.speech))
	renderedBubble := bubbleStyle.Render(wrappedSpeech)
	
	boxWidth := 16 + (m.weight * 2)
	if isEatingAction { boxWidth = 22 }
	
	petRender := petStyle.Width(boxWidth).Align(lipgloss.Center).Render(assembledPet)
	statusRender := statusStyle.Render(fmt.Sprintf("%s\n%s\n%s\n\n%s", hungerBar, happyBar, energyBar, emotionStr))
	
	var helpRender string
	if m.inStyleMode {
		helpRender = actionStyle.Render(
			"\n[1] Brush Teeth  |  [2] Change Hair" +
			"\n[3] Clothes      |  [c] Return Back",
		)
	} else {
		helpRender = statusStyle.Render(
			"\n[f] Feed  |  [p] Pet    |  [d] Dance  |  [b] Bath" +
			"\n[t] Toilet|  [s] Speak  |  [c] Style  |  [l] Loaf",
		)
	}

	body := fmt.Sprintf("%s\n\n%s\n\n%s\n%s", renderedBubble, petRender, statusRender, helpRender)
	lockedFrame := frameStyle.Width(54).Render(body)

	return lockedFrame + "\n"
}

func renderBar(value int) string {
	filled := value / 10
	if filled < 0 { filled = 0 }
	if filled > 10 { filled = 10 }
	return strings.Repeat("█", filled) + strings.Repeat(" ", 10-filled)
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
