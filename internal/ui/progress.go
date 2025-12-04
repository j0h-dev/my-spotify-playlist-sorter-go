package ui

type Progress interface {
	Start(total int, label string)
	Advance(n int)
	Finish()
}
