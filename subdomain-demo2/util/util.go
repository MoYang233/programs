package util

import (
	"bufio"
	"dnsgrep/vars"
	"encoding/json"
	"fmt"
	"github.com/liushuochen/gotable"
	"github.com/xuri/excelize/v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ReadTarget(filename string) []string {
	var tasks []string
	f,_ := os.Open(filename)
	defer f.Close()
	buff := bufio.NewReader(f)
	for i := 1 ; ; i++ {
		target,_,err := buff.ReadLine()
		if err != nil && err != io.EOF{
			panic(err)
		}else if err == io.EOF {
			break
		}
		tasks = append(tasks,string(target))
	}
	return tasks
}

func Start(tasks []string)  {
	wg := &sync.WaitGroup{}
	taskChan := make(chan string,50)
	for i:=1; i<5; i++{
		go Run(taskChan,wg)
	}
	for _,target := range tasks {
		wg.Add(1)
		taskChan <- target
	}
	close(taskChan)
	wg.Wait()
}

func Run(taskChan chan string,wg *sync.WaitGroup)  {
	for target := range taskChan{
		GetInfo(Request(strings.TrimSpace(target)))
		wg.Done()
	}
}

func Request(target string) ([]byte,string) {
	start := time.Now()
	client := &http.Client{Timeout: 5*time.Second}
	url := "https://www.dnsgrep.cn/api/query?q="+target+"&token=6fecc6d76090e8fd4ff0ebaa9af30c7d"
	resp,err := client.Get(url)
	if err != nil {
		fmt.Println("请求链接失败！",err)
		return nil,target
	}else if resp.Status != "200 OK" {
		fmt.Println("请求链接失败，状态码不为200！")
		return nil,target
	}
	defer resp.Body.Close()

	text,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败！",err)
		return nil,target
	}
	end := time.Since(start)
	fmt.Printf("请求数据用时：%s\n",end)
	return text,target
}

func GetInfo(jsonData []byte,target string)  {
	if jsonData == nil {
		fmt.Println("获取数据失败！")
		return
	}
	var v vars.JsonData
	err := json.Unmarshal(jsonData, &v)
	if err != nil {
		fmt.Println("数据解析失败！", err)
		return
	}
	vars.Mu.Lock()
	vars.V.Target = append(vars.V.Target,target)
	vars.V.ResultData = append(vars.V.ResultData,v.Data)
	vars.Mu.Unlock()
}

func OutPut()  {
	for i,r := range vars.V.ResultData{
		fmt.Printf("\n==============目标：%s共检索到%d条数据================\n\n",vars.V.Target[i],len(r.Data))
		tb,err := gotable.Create("domain","value","type","time")
		if err != nil{
			fmt.Println("创建表格失败！",err)
			os.Exit(0)
		}
		tb.AddRows(r.Data)
		tb.PrintTable()
	}

}

func CheckFileExist(filename string) bool {
	_,err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}else {
		return true
	}
}

func WriteExcel(filename string)  {
	f := excelize.NewFile()
	//根据target名称创建工作表
	for _,r := range vars.V.Target{
		f.NewSheet(r)
	}
	//写入数据
	var sheetName string
	for i,obj := range vars.V.ResultData{
		sheetName = vars.V.Target[i]
		//指定单元格，设置第一行的值
		f.SetCellValue(sheetName, "A1", "Domain")
		f.SetCellValue(sheetName, "B1", "Value")
		f.SetCellValue(sheetName, "C1", "Type")
		f.SetCellValue(sheetName, "D1", "Time")
		//指定工作表、单元格输出数据
		for i,r := range obj.Data{
			num := strconv.Itoa(i+2)
			f.SetCellValue(sheetName, "A"+num, r["domain"])
			f.SetCellValue(sheetName, "B"+num, r["value"])
			f.SetCellValue(sheetName, "C"+num, r["type"])
			f.SetCellValue(sheetName, "D"+num, r["time"])
		}
	}
	//设置工作簿默认工作表，即打开时显示的表
	f.SetActiveSheet(1)
	// 根据指定路径保存文件
	if err := f.SaveAs(filename); err != nil {
		fmt.Println("写入文件失败：",err)
		return
	}
	fmt.Printf("写入文件：%s成功！\n",filename)
}