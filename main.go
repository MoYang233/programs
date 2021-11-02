package main

import (
	"dnsgrep/util"
	"dnsgrep/vars"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"time"
)

func main()  {
	start := time.Now()
	//设置参数部分
	app := &cli.App{
		Name: "脚本名称",
		Author: "作者",
		Version: "2021-11-01",
		Usage: "简介",
		Flags: []cli.Flag{			//设定参数的标志
			&cli.StringFlag{
				Name: "t",			//参数名称
				Value: "",			//设置参数的默认值
				Usage: "target，域名或ip",	//参数介绍，将显示在-h中
				Destination: &vars.Target,		//将输入的参数赋值给vars包中定义的全局变量Target
			},
		},
		Action: func(c *cli.Context) {
			//检查是否输入了目标参数，没输入直接退出
			if c.IsSet("t") == false{
				fmt.Println("请输入查询目标，-h查看参数！")
				os.Exit(0)
			}
		},
	}
	app.Run(os.Args)
	if vars.Target == "" {
		os.Exit(0)
	}
	//发送请求并将返回值用作GetInfo()函数的输入
	util.GetInfo(util.Request(vars.Target))
	util.OutPut()
	util.WriteExcel()
	//统计程序运行时间
	end := time.Since(start)
	fmt.Printf("共用时：%s",end)
}