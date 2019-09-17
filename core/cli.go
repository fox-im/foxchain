package core

import(
	"fmt"
	"os"
)
//命令行接口
type CLI struct{
	blockchain *Blockchain
}

//用法
func (cli *CLI)printUsage()  {
	fmt.Println("用法如下")	
	fmt.Println("addblock 向区块链增加块")
	fmt.Println("showchain 显示区块链")
}


func (cli *CLI)validateArgs()  {
	if len(os.Args)<2{
		cli.printUsage() //显示用法
		os.Exit(1)  //退出
	}
}


func (cli *CLI)addBlock()  {
	cli.blockchain.AddBlock(data) // 增加区块
	fmt.Println("区块增加成功")
}

func (cli *CLI)showBlockchain()  {
	bci:=cli.blockchain.Iterator() //创建循环迭代器
	for{
		block:=bci.next()//取得下一个区块

	}
}


func (cli *CLI)Run()  {
	
}