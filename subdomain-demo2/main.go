package main

import (
	"dnsgrep/util"
	"dnsgrep/vars"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"strings"
	"time"
)

func main()  {
	start := time.Now()
	app := &cli.App{
		Name: "子域名查询&IP反查",
		Author: "作者",
		Version: "版本信息",
		Usage: "简介",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "t",
				Value: "",
				Usage: "target,命令格式：-t `baidu.com,127.0.0.1`",
				Destination: &vars.Target,
			},
			&cli.StringFlag{
				Name: "r",
				Value: "",
				Usage: "从文件导入目标，目标文件内一行一个，命令格式：-r `xx.txt`",
				Destination: &vars.FileName,
			},
			&cli.StringFlag{
				Name: "f",
				Value: "",
				Usage: "输出表格名称,命令格式：-f `result.xlsx`",
				Destination: &vars.ExcelName,
			},
		},
		Action: func(c *cli.Context) {
			if !c.IsSet("t") && !c.IsSet("r") {
				fmt.Println("请输入查询目标，-h查看参数！")
				os.Exit(0)
			}
		},
	}
	app.Run(os.Args)
	//检查输出文件是否已经存在
	if vars.ExcelName != ""{
		if util.CheckFileExist(vars.ExcelName){
			fmt.Println("保存的文件名已存在，请重新输入！")
			os.Exit(0)
		}
	}
	//根据输入的target生成最终target切片
	if vars.Target != "" || vars.FileName != "" {
		var tasks []string
		switch  {
		//命令行输入和文件输入同时存在
		case vars.Target != "" && vars.FileName != "":
			tasks = strings.Split(vars.Target,",")			//根据逗号分割字符串，返回切片
			tasks = append(tasks,util.ReadTarget(vars.FileName)...)		//将切片添加到切片中
		//只有命令行
		case vars.Target != "":
			tasks = strings.Split(vars.Target,",")
		//只有文件
		case vars.FileName != "":
			tasks = util.ReadTarget(vars.FileName)
		}
		fmt.Println("开始查询，待查询数：",len(tasks))
		util.Start(tasks)
		util.OutPut()
	}else {
		os.Exit(0)
	}
	//判断一下是否需要输出到文件
	if vars.ExcelName != ""{
		util.WriteExcel(vars.ExcelName)
	}
	end := time.Since(start)
	fmt.Printf("共用时：%s",end)
}