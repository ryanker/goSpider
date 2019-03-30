package model

type Rule struct {
	Rid    int64
	CateId int64
	Name   string
	Brief  string

	ListTable     string
	ListUrl       string
	ListPageStart int64
	ListPageEnd   int64
	ListPageSize  int64
	ListRange     string
	ListRule      string

	ContentTable string
	ContentUrl   string

	UpdateDate string
	CreateDate string
}
