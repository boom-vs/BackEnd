package webSocketModels

type Section struct {
	ID         uint
	Name       string
	Title      string
	Active     bool
	SortNumber uint
	Icon       string
	Sections   []*Section
}
