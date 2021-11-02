package vars

type JsonData struct {
	Status    int       `json:"status"`
	Data      Info      `json:"data"`
}
type Info struct {
	Data    []map[string]string  `json:"data"`
	Count   int
	Type    int
}

var (
	//声明一个全局变量V，JsonData类型，用于存放数据
	V JsonData
	Target string
)