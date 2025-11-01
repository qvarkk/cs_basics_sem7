package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"qvarkk/huffman/internal/huffman"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorReset  = "\033[0m"
)

type state int

const (
	stateMenu state = iota
	statePicker
	stateDone
)

type model struct {
	state 				state
	choices				[]string
	cursor				int
	selected			int
	filepicker		filepicker.Model
	help					help.Model
	selectedFile	string
	quitting			bool
	keys 					keyMap
	output				*compressionInfo
	err 					error
}

type keyMap struct {
	Up			key.Binding
	Down		key.Binding
	Select	key.Binding
	Quit		key.Binding
	Help    key.Binding
}

var defaultKeyMap = keyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("⏎", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "exit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Select}, {k.Help, k.Quit}}
}

func initialModel() model {
	fp := filepicker.New()
	fp.SetHeight(5)
	fp.AllowedTypes = []string{".txt", ".md", ".go", ".bin"}
	fp.CurrentDirectory, _ = os.UserHomeDir()
	return model{
		choices:		[]string{"Compress file", "Decompress file"},
		state:      stateMenu,
		filepicker: fp,
		help:       help.New(),
		keys:       defaultKeyMap,
	}
}

type clearErrorMsg struct {}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateMenu:	
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.Quit):
				m.quitting = true
				return m, tea.Quit

			case key.Matches(msg, m.keys.Up):
				if m.cursor > 0 {
					m.cursor--
				}

			case key.Matches(msg, m.keys.Down):
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}

			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll

			case key.Matches(msg, m.keys.Select):
				m.selected = m.cursor
				m.state = statePicker
				return m, m.filepicker.Init()
			}
		}
	case statePicker:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, m.keys.Quit) {
				m.quitting = true
				return m, tea.Quit
			}

		case clearErrorMsg:
			m.err = nil
		}

		var cmd tea.Cmd
		m.filepicker, cmd = m.filepicker.Update(msg)

		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.selectedFile = path
			m.state = stateDone
			m.keys = keyMap{}
			if m.selected == 0 {
				m.output, m.err = handleHuffmanCoding(path)			
			} else {
				m.err = handleHuffmanDecoding(path)
			}
		}

		if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
			m.err = fmt.Errorf("%s is not valid", path)
			m.selectedFile = ""
			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}

		return m, cmd
	case stateDone:
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "quitting...\n"
	}
	var s strings.Builder

	switch m.state {
	case stateMenu:
		s.WriteString("Select an option\n\n")

		for i, option := range m.choices {
			cursor := " "
			if m.cursor == i {
					cursor = ColorBlue + ">"
			}

			s.WriteString(fmt.Sprintf(" %s %s%s\n", cursor, option, ColorReset))
		}
		s.WriteString("\n")
	case statePicker:
		if m.err != nil {
			s.WriteString(fmt.Sprintf("%sError:%s %s\n\n", ColorRed, ColorReset, m.err.Error()))
		}
		s.WriteString("Pick a file:\n\n")
		s.WriteString(m.filepicker.View())
		s.WriteString("\n")
	case stateDone:
		s.WriteString(fmt.Sprintf("You chose: %s%s%s\n", ColorYellow, m.choices[m.selected], ColorReset))
		s.WriteString(fmt.Sprintf("On file: %s%s%s\n\n", ColorYellow, m.selectedFile, ColorReset))

		if m.selected == 0 && m.output != nil {
			s.WriteString(fmt.Sprintf("Size before compression: %s%d%s bytes\n", ColorRed, m.output.initialSize, ColorReset))
			s.WriteString(fmt.Sprintf("Size after compression: %s%d%s bytes\n\n", ColorGreen, m.output.compressedSize, ColorReset))
		} else {
			s.WriteString(fmt.Sprintf("File was saved with name %sdecompressed.txt%s\n\n", ColorBlue, ColorReset))	
		}

		if m.err != nil {
			s.WriteString(fmt.Sprintf("%sError:%s %v\n", ColorRed, ColorReset, m.err))
		} else {
			s.WriteString(fmt.Sprintf("Operation %ssuccessful%s\n\n", ColorGreen, ColorReset))
		}
		s.WriteString(fmt.Sprintf("%sPress any key to exit\n\n%s", ColorYellow, ColorReset))
	}

	s.WriteString(m.help.View(m.keys))
	return s.String()
}

type compressionInfo struct {
	initialSize			int64
	compressedSize	int64
}

func handleHuffmanCoding(path string) (*compressionInfo, error) {
	fileStat, err := os.Stat(path)
	if err != nil {
			return nil, fmt.Errorf("could not stat input file: %w", err)
	}

  initialSize := fileStat.Size()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read input file: %w", err)
	}

	var buf bytes.Buffer
	err = huffman.Code(string(data), &buf)
	if err != nil {
		return nil, fmt.Errorf("compression error: %w", err)
	}

	outputPath := "output.bin"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("could not create output file: %w", err)
	}
	defer outputFile.Close()

	_, err = buf.WriteTo(outputFile)
	if err != nil {
		return nil, fmt.Errorf("could not write compressed data: %w", err)
	}

	outStat, err := os.Stat(outputPath)
	if err != nil {
			return nil, fmt.Errorf("could not stat output file: %w", err)
	}
	compressedSize := outStat.Size()

	return &compressionInfo{
			initialSize:    initialSize,
			compressedSize: compressedSize,
	}, nil
}

func handleHuffmanDecoding(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("couldn't open compressed file: %w", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	reader := bufio.NewReader(file)
	err = huffman.Decode(reader, &buf)
	if err != nil {
		return fmt.Errorf("couldn't decompress file: %w", err)
	}

	outputPath := "decompressed.txt"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer outputFile.Close()

	_, err = buf.WriteTo(outputFile)
	if err != nil {
		return fmt.Errorf("could not write decompressed data: %w", err)
	}

	return nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}