package mcphost

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"

	_ "embed"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:embed openapi-go.json
var openapiSpecGo string

//go:embed openapi-python.json
var openapiPython string

type queryResult struct {
	response string
	isError  bool
}

type queryResultCmd tea.Cmd

func queryResultCmdFunc(qr queryResult) queryResultCmd {
	response := infoStyle.Render(qr.response)
	if qr.isError {
		response = errStyle.Render(fmt.Sprintf("Error: %s", qr.response))
	}
	return queryResultCmd(
		func() tea.Msg {
			return queryResult{
				response: response,
				isError:  qr.isError,
			}
		})
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Background(lipgloss.Color("yellow"))
	}()

	errStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("red"))
	}()

	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("blue"))
)

type model struct {
	llmApp *MCPHost
	ctx    context.Context

	openapiGo     string
	openapiPython string

	showViewport  bool
	viewPortTitle *string
	viewport      viewport.Model

	query    textinput.Model
	response *textinput.Model
	working  *bool
	spinner  spinner.Model

	appChan chan tea.Msg

	mux *sync.RWMutex
}

func (m *model) setQueryResult(response string, isError bool) {
	m.appChan <- queryResultCmdFunc(queryResult{
		response: response,
		isError:  isError,
	})
}

func initializeModel(ctx context.Context, llmApp *MCPHost) model {
	var m model

	m.appChan = make(chan tea.Msg, 10) // Buffered channel to handle messages
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-m.appChan:
				_, _ = m.Update(msg)
			default:
				continue
			}
		}
	}()

	m.working = new(bool)
	m.viewPortTitle = new(string)

	m.llmApp = llmApp
	m.ctx = ctx
	m.openapiGo = openapiSpecGo
	m.openapiPython = openapiPython
	m.response = &textinput.Model{}
	m.query = textinput.New()
	m.query.Placeholder = "Ask me anything..."
	m.query.Focus()

	m.mux = &sync.RWMutex{}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" {
			return m, tea.Quit
		}
		if k := msg.String(); k == "ctrl+o" || k == "ctrl+y" {
			m.showViewport = !m.showViewport
			m.viewport = viewport.New(80, 24) // default size, will be updated on resize
			content := m.openapiGo
			title := "Go MCP Server"

			if k == "ctrl+y" {
				content = m.openapiPython
				title = "Python MCP Server"
			}

			m.viewport.SetContent(content)
			m.viewPortTitle = &title
		}
		if k := msg.String(); k == "enter" && !m.isWorking() {
			m.switchWorkingFlag()
			m.spinner.Spinner = spinner.Line

			go func() {
				response, err := m.llmApp.RunUserQuery(m.ctx, m.query.Value())
				if err != nil {
					response = fmt.Sprintf("Error: %s", err.Error())
				}
				m.setQueryResult(response, err != nil)
			}()

			return m, tea.Batch(
				m.spinner.Tick,
			)
		}

		if k := msg.String(); k == "ctrl+r" {
			// Reset the query and response fields
			m.query.Reset()
			m.response.Reset()
		}

	case queryResultCmd:
		qrCmd := msg().(queryResult)
		var res string
		if qrCmd.isError {
			res = errStyle.Render(qrCmd.response)
		} else {
			res = qrCmd.response
		}
		m.response.SetValue(res)
		m.switchWorkingFlag()
		return m, cmd

	case spinner.TickMsg:
		if !m.isWorking() {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		if m.showViewport {
			headerHeight := lipgloss.Height(m.headerView())
			footerHeight := lipgloss.Height(m.footerView())
			verticalMarginHeight := headerHeight + footerHeight
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	}

	// Handle keyboard and mouse events in the viewport only if it's shown
	if m.showViewport {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	if *m.working {
		cmds = append(cmds, m.spinner.Tick)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.showViewport {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	}

	sb := strings.Builder{}
	sb.WriteString("Enter your query:\n\n")
	sb.WriteString(m.query.View())
	sb.WriteString("\n\n")

	if *m.working {
		sb.WriteString("Processing your request...\n\n")
		sb.WriteString(spinnerStyle.Render(m.spinner.View()))
		sb.WriteString("\n\n")
	}

	if m.response.Value() != "" {
		sb.WriteString(m.response.View())
		sb.WriteString("\n\n")
	}

	sb.WriteString("\n Ctrl+o to show/hide spec from Go MCP server")
	sb.WriteString("\n Ctrl+y to show/hide spec from Python MCP server")
	sb.WriteString("\n Ctrl+r to reset query, ")
	sb.WriteString("\n Q to quit \n")
	return sb.String()
}

func (m model) headerView() string {
	title := titleStyle.Render(*m.viewPortTitle)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *model) switchWorkingFlag() {
	m.mux.Lock()
	*m.working = !*m.working
	m.mux.Unlock()
}

func (m *model) isWorking() bool {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return *m.working
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

// RunUI starts the user interface for the MCPHost.
func RunUI(ctx context.Context, llmApp *MCPHost) error {
	p := tea.NewProgram(
		initializeModel(ctx, llmApp),
		tea.WithAltScreen(), // use the full size of the terminal in its "alternate screen buffer"
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run UI: %w", err)
	}
	return nil
}
