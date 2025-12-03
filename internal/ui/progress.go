package ui

import "github.com/schollz/progressbar/v3"

type ProgressBar struct {
	bar *progressbar.ProgressBar
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{}
}

func (p *ProgressBar) Start(total int, label string) {
	p.bar = progressbar.Default(int64(total), label)
}

func (p *ProgressBar) Advance(n int) {
	if p.bar != nil {
		p.bar.Add(n)
	}
}

func (p *ProgressBar) Finish() {
	if p.bar != nil {
		p.bar.Finish()
	}
}
