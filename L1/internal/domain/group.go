package domain

type GroupType int

const (
	Project GroupType = iota
	Org
	CorpGroup
)

type Group struct {
	ID    int
	Name  string
	Type  GroupType
	Users []int
}
