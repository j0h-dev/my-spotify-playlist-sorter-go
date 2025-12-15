package components

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type ProgressBar struct {
	bar progress.Model

	label   string
	total   int
	percent float64

	updateList []tea.Cmd
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		bar: progress.New(progress.WithDefaultGradient()),
	}
}

func (p *ProgressBar) Start(total int, label string) {
	p.label = label
	p.total = total
	p.percent = 0.0
}

func (p *ProgressBar) Advance(n int) {
	if p.total > 0 {
		p.percent += float64(n) / float64(p.total)
	}
}

func (p *ProgressBar) Finish() {
	p.bar.IncrPercent(100)
}

func (p *ProgressBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p.bar.Update(msg)
}

func (p *ProgressBar) View() string {
	return p.bar.ViewAs(p.percent)
}

func (p *ProgressBar) Width() int {
	return p.bar.Width
}

func (p *ProgressBar) SetWidth(width int) {
	p.bar.Width = width
}

func (p *ProgressBar) Label() string {
	return p.label
}

func (p *ProgressBar) SetLabel(label string) {
	p.label = label
}
