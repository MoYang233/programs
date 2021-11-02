package util

import (
	"dnsgrep/vars"
	"encoding/json"
	"fmt"
	"github.com/liushuochen/gotable"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Request(target string) []byte {
	start := time.Now()
	client := &http.Client{Timeout: 5*time.Second}
	url := "https://www.dnsgrep.cn/api/query?q="+target+"&token=[申请的token值]"
	resp,err := client.Get(url)
	if err != nil {
		fmt.Println("请求链接失败！",err)
		return nil
	}else if resp.Status != "200 OK" {
		fmt.Println("请求链接失败，状态码不为200！")
		return nil
	}
	defer resp.Body.Close()

	text,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败！",err)
		return nil
	}
	end := time.Since(start)
	fmt.Printf("请求数据用时：%s",end)
	return text
}

func GetInfo(jsonData []byte) {
	if jsonData == nil {
		fmt.Println("获取数据失败！")
		return
	}

	err := json.Unmarshal(jsonData, &vars.V)
	if err != nil {
		fmt.Println("数据解析失败！", err)
		return
	}
}

func OutPut()  {
	fmt.Printf("共检索到%d条数据\n",len(vars.V.Data.Data))
	tb,err := gotable.Create("domain","value","type","time")
	if err != nil{
		fmt.Println("创建表格失败！",err)
		os.Exit(0)
	}
	tb.AddRows(vars.V.Data.Data)
	tb.PrintTable()
}

func WriteExcel()  {
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "Domain")
	f.SetCellValue("Sheet1", "B1", "Value")
	f.SetCellValue("Sheet1", "C1", "Type")
	f.SetCellValue("Sheet1", "D1", "Time")

	for i,r := range vars.V.Data.Data{
		num := strconv.Itoa(i+2)
		f.SetCellValue("Sheet1", "A"+num, r["domain"])
		f.SetCellValue("Sheet1", "B"+num, r["value"])
		f.SetCellValue("Sheet1", "C"+num, r["type"])
		f.SetCellValue("Sheet1", "D"+num, r["time"])
	}

	if err := f.SaveAs("Result.xlsx"); err != nil {
		println(err.Error())
	}
}