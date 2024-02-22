package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	bb "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"math"
	"os"
	"path/filepath"
	"strings"
)

var vite_templates = [][]interface{}{
	{"svelte", [][]string{
		{"Javascript", "svelte"},
		{"Typescript", "svelte-ts"},
	}},
	{"react", [][]string{
		{"Javascript", "react"},
		{"Typescript", "react-ts"},
	}},
	{"vue", [][]string{
		{"Javascript", "vue"},
		{"Typescript", "vue-ts"},
	}},
	{"preact", [][]string{
		{"Javascript", "preact"},
		{"Typescript", "preact-ts"},
	}},
	{"solid", [][]string{
		{"Javascript", "solid"},
		{"Typescript", "solid-ts"},
	}},
	{"lit", [][]string{
		{"Javascript", "lit"},
		{"Typescript", "lit-ts"},
	}},
	{"vanilla", [][]string{
		{"Javascript", "vanilla"},
		{"Typescript", "vanilla-ts"},
	}},
	{"react-swc", [][]string{
		{"Javascript", "react-swc"},
		{"Typescript", "react-swc-ts"},
	}},
	{"qwik", [][]string{
		{"Javascript", "qwik"},
		{"Typescript", "qwik-ts"},
	}},
}
var did_exit bool = false

func main() {
	mo := model{
		vite_templates:         vite_templates,
		selected_head_template: -1,
		cursor:                 0,
	}
	p := bb.NewProgram(&mo, bb.WithAltScreen())
	m, e := p.Run()
	if e != nil {
		panic(e)
	}
	if did_exit {
		return
	}

	selectedTemplate := m.(*model).selectedTemplate
	selectedTemplateValue := vite_templates[m.(*model).selected_head_template][1].([][]string)[selectedTemplate][1]

	textinputmodel := textinput.New()
	textinputmodel.Placeholder = "./"
	textinputmodel.Focus()
	textinputmodel.Prompt = ""
	textinputmodel.ShowSuggestions = true
	textinputmodel.Width = 50
	textinputmodel.CompletionStyle = lg.NewStyle().Foreground(lg.Color("#ff00ff"))
	textinputmodel.CursorStart()

	mo2 := selectPathModel{
		pathInput: textinputmodel,
	}
	p2 := bb.NewProgram(&mo2, bb.WithAltScreen())
	m2, e2 := p2.Run()
	if e2 != nil {
		panic(e2)
	}
	if did_exit {
		return
	}
	path := m2.(*selectPathModel).pathInput.Value()

	fmt.Println("vite create", selectedTemplateValue, path)
}

type model struct {
	vite_templates         [][]interface{}
	selected_head_template int
	cursor                 int
	selectedTemplate       int
}

type selectPathModel struct {
	pathInput       textinput.Model
	suggestionIndex int
}

func (m *model) Init() bb.Cmd {
	return nil
}

func (m *model) Update(msg bb.Msg) (bb.Model, bb.Cmd) {
	is_in_sub_menu := m.selected_head_template >= 0
	switch msg := msg.(type) {
	case bb.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, bb.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				if is_in_sub_menu {
					m.cursor = len(vite_templates[m.selectedTemplate]) - 1
				} else {
					m.cursor = len(vite_templates) - 1
				}
			}
		case "down", "j":
			if is_in_sub_menu {
				if m.cursor < len(vite_templates[m.selectedTemplate])-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			} else {
				if m.cursor < len(vite_templates)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			}
		case "enter":
			if is_in_sub_menu {
				m.selectedTemplate = m.cursor
				return m, bb.Quit
			} else {
				m.selected_head_template = m.cursor
			}
		case "left", "h":
			if is_in_sub_menu {
				m.selected_head_template = -1
			} else {
				did_exit = true
				return m, bb.Quit
			}
		}

	}
	return m, nil
}

func (m *model) View() string {
	is_sub_menu := m.selected_head_template >= 0
	var s []string
	header := lg.NewStyle().Bold(true).PaddingBottom(1).Foreground(lg.Color("#3399ff"))
	if is_sub_menu {
		items := vite_templates[m.selected_head_template][1].([][]string)
		s = append(s, capitalize(header.Render(vite_templates[m.selected_head_template][0].(string))))
		for i, item := range items {
			item[0] = capitalize(item[0])
			if i == m.cursor {
				s = append(s, lg.NewStyle().Foreground(lg.Color("#ff00ff")).Render(item[0]))
			} else {
				s = append(s, item[0])
			}
		}
	} else {
		s = append(s, capitalize(header.Render("Vite Templates")))
		for i, item := range vite_templates {
			item[0] = capitalize(item[0].(string))
			if i == m.cursor {
				s = append(s, lg.NewStyle().Foreground(lg.Color("#ff00ff")).Render(item[0].(string)))
			} else {
				s = append(s, item[0].(string))
			}
		}
	}
	bottom := lg.NewStyle().Foreground(lg.Color("#9e9e9e")).PaddingTop(1).Faint(true)
	if is_sub_menu {
		s = append(s, bottom.Render("press left or h to go back"))
	} else {
		s = append(s, bottom.Render("press left,h,esc or ctrl-c to go exit"))
	}

	w, h, _ := term.GetSize(0)

	longest := w * 10 / 5 / 10
	for _, item := range s {
		longest = int(math.Max(float64(longest), float64(lg.Width(item))))
	}
	for i, item := range s {
		s[i] = lg.NewStyle().Width(longest).Render(item)
	}
	app := lg.NewStyle().Border(lg.InnerHalfBlockBorder()).Background(lg.Color("#131313")).MaxWidth(w-w*10/9).Padding(1, 2).BorderForeground(lg.Color("#25b86e")).Width(longest + 2).Render(lg.JoinVertical(lg.Left, s...))

	return lg.Place(w, h, lg.Center, lg.Center, app)
}

func capitalize(str string) string {
	first := str[0]
	return strings.ToUpper(string(first)) + str[1:]
}

func (m *selectPathModel) Init() bb.Cmd {
	return textinput.Blink
}

func (m *selectPathModel) Update(msg bb.Msg) (bb.Model, bb.Cmd) {
	if !m.pathInput.Focused() {
		return m.UpdateSuggestions(msg)
	}
	switch msg := msg.(type) {
	case bb.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			did_exit = true
			return m, bb.Quit
		case "enter":
			return m, bb.Quit
		}
	}
	// get files and directories in the current directory
	// and set the suggestions to the input
	p := m.pathInput.Value()
	// replace ~ with the home directory
	if strings.HasPrefix(p, "~") {
		home, _ := os.UserHomeDir()
		p = strings.Replace(p, "~", home, 1)
	}
	files, err := os.ReadDir(p)
	if err == nil {
		var suggestions []string
		for _, file := range files {
			suggestions = append(suggestions, file.Name())
		}
		was_empty := len(m.pathInput.AvailableSuggestions()) == 0
		if was_empty {
			m.suggestionIndex = 0
		}
		m.pathInput.SetSuggestions(suggestions)
	} else {
		base := filepath.Base(p)
		dir := filepath.Dir(p)
		files, err := os.ReadDir(dir)
		if err == nil {
			var suggestions []string
			for _, file := range files {
				if strings.HasPrefix(file.Name(), base) {
					suggestions = append(suggestions, file.Name())
				}
			}
			was_empty := len(m.pathInput.AvailableSuggestions()) == 0
			if was_empty {
				m.suggestionIndex = 0
			}
			m.pathInput.SetSuggestions(suggestions)
		} else {
			m.pathInput.SetSuggestions([]string{})
		}
	}

	var cmd bb.Cmd
	m.pathInput, cmd = m.pathInput.Update(msg)
	return m, cmd
}

func (m *selectPathModel) View() string {
	var s []string
	s = append(s, lg.NewStyle().Foreground(lg.Color("#3399ff")).PaddingBottom(1).Bold(true).Render("Where do you want to create the project?"))
	s = append(s, m.pathInput.View())
	sugs := m.pathInput.AvailableSuggestions()[0:int(math.Min(float64(len(m.pathInput.AvailableSuggestions())), 5))]
	for i, sug := range sugs {
		if i == m.suggestionIndex {
			sugs[i] = lg.NewStyle().Foreground(lg.Color("#ff00ff")).Render(sug)
		}
	}
	sug := lg.JoinVertical(lg.Left, sugs...)
	s = append(s, sug)
	w, h, _ := term.GetSize(0)
	longest := w * 10 / 5 / 10
	for _, item := range s {
		longest = int(math.Max(float64(longest), float64(lg.Width(item))))
	}
	for i, item := range s {
		s[i] = lg.NewStyle().Width(longest).Background(lg.Color("#131313")).Render(item)
	}
	app := lg.NewStyle().Border(lg.InnerHalfBlockBorder()).Background(lg.Color("#131313")).MaxWidth(100).Padding(1, 2).BorderForeground(lg.Color("#25b86e")).Render(
		lg.JoinVertical(lg.Left, s...),
	)
	return lg.Place(w, h, lg.Center, lg.Center, app)
}

func (m *selectPathModel) UpdateSuggestions(msg bb.Msg) (bb.Model, bb.Cmd) {
	switch msg := msg.(type) {
	case bb.KeyMsg:
		switch msg.String() {
		case "up", "k", "shift+tab":
			if m.suggestionIndex > 0 {
				m.suggestionIndex--
			}
		case "down", "j", "tab":
			if m.suggestionIndex < len(m.pathInput.AvailableSuggestions())-1 {
				m.suggestionIndex++
			}
		case "enter":
			m.pathInput.SetValue(m.pathInput.AvailableSuggestions()[m.suggestionIndex])
		case "esc", "ctrl+c":
			m.suggestionIndex = 0
			m.pathInput.SetSuggestions([]string{})
			m.pathInput.Focus()
		}
	}
	return m, nil

}
