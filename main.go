package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type model struct {
	width, height int
	Tabs          []string
	activeTab     int

	requestHistory []RequestTime
	requestTime    RequestTime
	inputField     textinput.Model
	feedback       string
	helpListItems  []helpItem
	responseTable  table.Model
	historyTable   table.Model
	styles         *Styles

	debounceEnterKey bool
}

type helpItem struct {
	title, desc string
}

func initialModel() model {
	m := model{}

	m.width = 0
	m.height = 0

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

	historyTable := table.New(
		table.WithColumns(
			[]table.Column{
				{Title: "ID", Width: 3},
				{Title: "URL", Width: 22},
				//{Title: "IP", Width: 16},
				{Title: "DNS", Width: 10},
				{Title: "TCP", Width: 10},
				{Title: "TLS", Width: 10},
				{Title: "Server Processing", Width: 10},
				{Title: "Transfer", Width: 10},
				{Title: "Total", Width: 10},
			}),
		table.WithHeight(20),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	historyTable.SetStyles(s)
	m.historyTable = historyTable

	m.styles = DefaultStyles()
	m.inputField = ti

	m.helpListItems = []helpItem{
		{title: "DNS Lookup", desc: "The time it takes to look up the IP-Address(es) associated with this domain name"},
		{title: "TCP Connection", desc: "The time it takes to connect to the previously found IP-Address over TCP"},
		{title: "TLS Handshake", desc: "The time it takes the two communicating sides to exchange messages to verify each other and establish the cryptographic algorithms they will use"},
		{title: "Server Processing", desc: "The time that the server needs to process the request before sending a response."},
		{title: "Content Transfer", desc: "The time that the server needs to send an response"},
		{title: "Total", desc: "Combined Time of this whole Request"},
	}

	tabs := []string{"Request", "History", "Help", "About"}
	m.Tabs = tabs

	return m
}

func tick() tea.Msg {
	time.Sleep(time.Second / 4)
	return tickMsg(1)
}

func (m model) Init() tea.Cmd {
	return tick
}

type debounceEnterKey bool

type tickMsg int

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
				m.feedback = m.styles.Error.Width(100).Align(lipgloss.Center).Render(err.Error())
			} else {
				m.feedback = m.styles.Success.Render("Successfully")
				m.requestTime = resp
				m.requestHistory = append(m.requestHistory, resp)
				m.historyTable = m.generateHistoryTableView()
			}
			return m, tea.Tick(m.requestTime.contentTransfer, func(_ time.Time) tea.Msg {
				return debounceEnterKey(m.debounceEnterKey)
			})
		case "tab":
			m.activeTab = m.getNextTab()
			return m, nil
		case "shift+tab":
			m.activeTab = m.gePreviousTab()
			return m, nil
		case "down":
			m.historyTable.Focus()
			if m.historyTable.Cursor() == -1 {
				m.historyTable.SetCursor(0)
			}
			m.historyTable.MoveDown(0)
		case "up":
			m.historyTable.Focus()
			if m.historyTable.Cursor() == -1 {
				m.historyTable.SetCursor(0)
			}
			m.historyTable.MoveUp(0)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		w, h, _ := term.GetSize(int(os.Stdout.Fd()))
		return m, tea.Batch(tick, func() tea.Msg { return tea.WindowSizeMsg{Width: w, Height: h} })
	}

	m.inputField, _ = m.inputField.Update(msg)
	m.historyTable, cmd = m.historyTable.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width < 110 || m.height < 38 {
		s := ""
		s = "Terminal too small to display content\n"
		s += "Please resize your terminal to at least 110x38\n\n"
		width := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render(fmt.Sprintf("Width: %d", m.width))
		if m.width < 110 {
			width = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Width: %d", m.width))
		}
		height := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render(fmt.Sprintf("Height: %d", m.height))
		if m.height < 38 {
			height = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Height: %d", m.height))
		}
		s += "Your current terminal size is: " + fmt.Sprintf("%s x %s", width, height) + "\n\n"
		return fmt.Sprintf(s)
	}
	doc := strings.Builder{}

	tabContent := m.selectTabContent() + "\n"
	tabView := m.tabView()

	doc.WriteString(tabView)
	doc.WriteString("\n")
	doc.WriteString(m.styles.WindowStyle.Width(tabContentWidth).Render(tabContent))

	return m.styles.DocStyle.Render(doc.String() + "\n" + time.Now().String())
	//return m.viewport.View()
}

var tabContentWidth = 106

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

func (m model) selectTabContent() string {
	var view string
	if m.activeTab == 0 {
		view = lipgloss.JoinVertical(
			lipgloss.Center,
			m.styles.InputField.Render(m.inputField.View()),
			m.feedback,
			m.responseTableView(),
		)
	} else if m.activeTab == 1 {
		view = m.historyTable.View()
		m.historyTable.Focus()
	} else if m.activeTab == 2 {
		for _, item := range m.helpListItems {
			title := lipgloss.NewStyle().Width(100).Align(lipgloss.Left).Render(item.title)
			desc := lipgloss.NewStyle().Width(100).Align(lipgloss.Left).Foreground(m.styles.SubduedColor).Render(item.desc)
			view = view + title + "\n" + desc + "\n\n"
		}
	} else {
		view = lipgloss.NewStyle().
			Padding(0, 1).
			Italic(true).
			Background(lipgloss.Color("#fc4c92")).
			Foreground(lipgloss.Color("#ffffff")).
			Render("Created by nikita-t1") +
			"\n\nLink to Source Code: https://github.com/nikita-t1/responseTime\n\n"
	}
	return view
}

func (m model) generateHistoryTableView() table.Model {
	var rows []table.Row
	for _, requestTime := range m.requestHistory {
		rows = append(rows, table.Row{
			lpad(strconv.Itoa(requestTime.id), "0", 3),
			strings.ReplaceAll(requestTime.url, "https://", ""),
			//requestTime.ip,
			requestTime.dnsLookup.String(),
			requestTime.connectTime.String(),
			requestTime.tlsHandshake.String(),
			requestTime.serverProcessing.String(),
			requestTime.contentTransfer.String(),
			requestTime.total.String(),
		})
	}
	m.historyTable.SetRows(rows)
	return m.historyTable
}

func (m model) tabView() string {
	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, _, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = m.styles.ActiveTabStyle.Copy()
		} else {
			style = m.styles.InactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	rowFill := strings.Repeat("─", (tabContentWidth+1)-lipgloss.Width(row)) + "╮"
	row = row + lipgloss.NewStyle().Foreground(m.styles.HighlightColor).Render(rowFill)
	return row
}

func (m model) responseTableView() string {
	r := m.requestTime

	var ipList []table.Row
	for index, addr := range r.addrs {
		ipList = append(ipList, table.Row{"        " + strconv.Itoa(index+1), addr.IP.String()})
	}

	statusColor := lipgloss.NewStyle()
	if strings.HasPrefix(m.requestTime.status, "2") {
		statusColor = m.styles.Success
	} else if strings.HasPrefix(m.requestTime.status, "3") {
		statusColor = m.styles.Warning
	} else if strings.HasPrefix(m.requestTime.status, "4") || strings.HasPrefix(m.requestTime.status, "5") {
		statusColor = m.styles.Error
	}

	var rows []table.Row
	rows = []table.Row{
		{"ID", lpad(strconv.Itoa(r.id), "0", 3)},
		{"URL", r.url},
		{"IP", r.ip},
		{"Status", statusColor.Render(r.status)},
		{"DNS Lookup", r.dnsLookup.String()},
		{"TCP Connection", r.connectTime.String()},
		{"TLS Handshake", r.tlsHandshake.String()},
		{"Server Processing", r.serverProcessing.String()},
		{"Content Transfer", r.contentTransfer.String()},
		{"Total", r.total.String()},
		{"", ""},
		{"Alternative Addrs:", ""},
	}
	rows = append(rows, ipList...)
	m.responseTable.SetRows(rows)

	m.responseTable.SetStyles(m.styles.Table)

	return m.responseTable.View()
}

func (m model) getNextTab() int {
	activeTab := m.activeTab
	if activeTab < len(m.Tabs)-1 {
		return activeTab + 1
	} else {
		return 0
	}
}

func (m model) gePreviousTab() int {
	activeTab := m.activeTab
	if activeTab == 0 {
		return len(m.Tabs) - 1
	} else {
		return activeTab - 1
	}
}

func lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}
