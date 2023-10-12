package webSocketModels

type BuilderGetSections struct {
	ID         uint
	Name       string
	Title      string
	Active     bool
	SortNumber uint
	ParentID   *uint
	Icon       string
	Sections   []*BuilderGetSections
}
