package vars

import "sync"

type JsonData struct {
	Status    int       `json:"status"`
	Data      Info      `json:"data"`
}
type Info struct {
	Data    []map[string]string  `json:"data"`
	Count   int
	Type    int
}

type Result struct {
	Target     []string
	ResultData []Info
}

var (
	V         Result
	Target    string
	FileName  string
	ExcelName string
	Mu        sync.Mutex
)
