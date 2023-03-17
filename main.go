package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"strconv"
	"time"
)

type model struct {
	requestTime   RequestTime
	inputField    textinput.Model
	feedback      string
	responseTable table.Model
	styles        *Styles

	debounceEnterKey bool
}

func initialModel() model {
	m := model{}

	ti := textinput.New()
	ti.Placeholder = "https://google.com"
	ti.Focus()

	responseTable := table.New(
		table.WithColumns(
			[]table.Column{
				{Title: "", Width: 36},
				{Title: "", Width: 56},
			}),
		table.WithHeight(20),
	)
	m.responseTable = responseTable

	m.styles = DefaultStyles()
	m.inputField = ti
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

type debounceEnterKey bool

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case debounceEnterKey:
		if msg {
			m.debounceEnterKey = false
			return m, nil
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.debounceEnterKey {
				return m, nil
			}
			m.debounceEnterKey = true
			m.feedback = "Loading"
			resp, err := ExecuteRequest(m.inputField.Value())
			if err != nil {
				m.feedback = m.styles.Error.Render(err.Error())
			} else {
				m.feedback = m.styles.Success.Render("Successfully")
				m.requestTime = resp
			}
			return m, tea.Tick(m.requestTime.contentTransfer, func(_ time.Time) tea.Msg {
				return debounceEnterKey(m.debounceEnterKey)
			})
		}
	}
	m.inputField, cmd = m.inputField.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.styles.InputField.Render(m.inputField.View()),
		m.feedback,
		m.responseTableView(),
	)
}

func main() {
	f, _ := tea.LogToFile("debug.log", "")
	defer f.Close()
	log.Printf("Start")

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (m model) responseTableView() string {
	r := m.requestTime

	var ipList []table.Row
	for index, addr := range r.addrs {
		ipList = append(ipList, table.Row{"        " + strconv.Itoa(index+1), addr.IP.String()})
	}

	var rows []table.Row
	rows = []table.Row{
		{"ID", strconv.Itoa(r.id)},
		{"URL", r.url},
		{"IP", r.ip},
		{"Status", r.status},
		{"DNS Lookup", r.dnsLookup.String()},
		{"TCP Connection", r.connectTime.String()},
		{"TLS Handshake", r.tlsHandshake.String()},
		{"Server Processing", r.serverProcessing.String()},
		{"Content Transfer", r.contentTransfer.String()},
		{"", ""},
		{"Alternative Addrs:", ""},
	}
	rows = append(rows, ipList...)
	m.responseTable.SetRows(rows)

	m.responseTable.SetStyles(m.styles.Table)

	return m.responseTable.View()

}
