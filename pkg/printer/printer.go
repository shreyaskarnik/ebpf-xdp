package printer

import (
	"fmt"
	"net"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cilium/ebpf"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyIP    = "ip"
	columnKeyCount = "count"
)

type data struct {
	IP    string
	Count uint32
}

type Model struct {
	simpleTable table.Model
	data        []*data
	ebpfMap     *ebpf.Map
}

func NewModel(e *ebpf.Map) Model {
	return Model{
		simpleTable: table.New(
			[]table.Column{
				table.NewColumn(columnKeyIP, "IPv4 Address", 45),
				table.NewColumn(columnKeyCount, "Count", 10),
			},
		).BorderRounded().
			SortByDesc(columnKeyCount).
			Focused(true),
		ebpfMap: e,
	}
}

func (m Model) Init() tea.Cmd {
	return m.ProcessData
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.simpleTable, cmd = m.simpleTable.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			fmt.Println("quitting...")
			cmds = append(cmds, tea.Quit)
		}
	case []*data:
		m.simpleTable = m.simpleTable.WithRows(generateRows(msg))
		cmds = append(cmds, func() tea.Msg {
			time.Sleep(5 * time.Second)
			return m.ProcessData()
		})

	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}
	pad := lipgloss.NewStyle().Padding(1)
	body.WriteString(pad.Render(m.simpleTable.View()))
	return body.String()
}

func (m Model) ProcessData() tea.Msg {
	var (
		key []byte
		val uint32
		d   []*data
	)
	iter := m.ebpfMap.Iterate()
	for iter.Next(&key, &val) {
		sourceIP := net.IP(key)
		packetCount := val
		d = append(d, &data{
			IP:    sourceIP.String(),
			Count: packetCount,
		})
	}
	return d
}

func generateRows(d []*data) []table.Row {
	var rows []table.Row
	for _, d := range d {
		rows = append(rows, table.NewRow(
			table.RowData{
				columnKeyIP:    d.IP,
				columnKeyCount: fmt.Sprintf("%d", d.Count),
			},
		))
	}
	return rows
}
