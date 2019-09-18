package core

import(
	"fmt"
	"os"
	"strconv"
	"log"
	"flag"
)
//命令行接口
type CLI struct{
	Blockchain *Blockchain
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


func (cli *CLI)addBlock(data string)  {
	cli.Blockchain.AddBlock(data) // 增加区块
	fmt.Println("区块增加成功")
}

func (cli *CLI)showBlockchain()  {
	bci:=cli.Blockchain.Iterator() //创建循环迭代器
	for{
		block:=bci.next()//取得下一个区块
		fmt.Printf("上一块hash:%x",block.PrevBlockHash)
		fmt.Println("\n")
		fmt.Printf("数据:%s",block.Data)
		fmt.Println("\n")
		fmt.Printf("当前hash:%x",block.Hash)
		pow:=NewProofOfWork(block)
		fmt.Printf("pow: %s",strconv.FormatBool(pow.Validate()))
		fmt.Println("\n")
		fmt.Println("\n")
		if len(block.PrevBlockHash)==0{ //遇到创世区块
			break
		}
	}
}


func (cli *CLI)Run()  {
	cli.validateArgs() //校验
	//处理命令行参数
	addblockcmd:=flag.NewFlagSet("addblock",flag.ExitOnError)
	showchaincmd:=flag.NewFlagSet("showchain",flag.ExitOnError)
	addBlockData:=addblockcmd.String("data","","Block data")
	switch os.Args[1] {
	case "addblock":
		err:=addblockcmd.Parse(os.Args[2:]) //解析参数
		if err!=nil{
			log.Panic(err)
		}
	case "showchain":
		err:=showchaincmd.Parse(os.Args[2:]) //解析参数
		if err!=nil{
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if addblockcmd.Parsed(){
		if *addBlockData==""{
			addblockcmd.Usage()
			os.Exit(1)
		}else{
			cli.addBlock(*addBlockData) //增加区块
		}
	}
	if showchaincmd.Parsed(){
		cli.showBlockchain() //显示区块链
	}
	
}