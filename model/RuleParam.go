package model

type RuleParam struct {
	Pid        int64
	Rid        int64
	Type       int64
	Field      string
	Rule       string
	ValueType  int64
	ValueAttr  string
	FilterType int64
	FilterReg  string
	Sort       int64
	IsDown     int64
	CreateDate string
}
