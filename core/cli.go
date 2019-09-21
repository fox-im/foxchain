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

func (cli *CLI)createBlockChain(address string){
	bc:=createBlockChain(address) //创建区块链
	bc.DB.Close()
	fmt.Println("创建成功",address)
}

func (cli *CLI) getBalance(address string){
	bc:=NewBlockchain(address) //根据地址创建
	defer bc.DB.Close()
	balance:=0
	UTXOs:=bc.FindUTXO(address) //查找交易金额
	for _,out:=range UTXOs{
		balance+=out.Value //取出金额
	}
	fmt.Printf("查询的金额如下%s :%d \n",address,balance)
}

//用法
func (cli *CLI)printUsage()  {
	fmt.Println("用法如下")	
	fmt.Println("getbalance -address 你输入的地址   根据地址查询金额")
	fmt.Println("createblockchain -address 你输入的地址   根据地址创建区块链")
	// fmt.Println("addblock 向区块链增加块")
	fmt.Println("send -from From -to To -amount Amount  转账")
	fmt.Println("showchain 显示区块链")
}


func (cli *CLI)validateArgs()  {
	if len(os.Args)<2{
		cli.printUsage() //显示用法
		os.Exit(1)  //退出
	}
}


// func (cli *CLI)addBlock(data string)  {
// 	cli.Blockchain.AddBlock(data) // 增加区块
// 	fmt.Println("区块增加成功")
// }

func (cli *CLI)showBlockchain()  {
	bc:=NewBlockchain("")
	defer bc.DB.Close()

	bci:=bc.Iterator() //创建循环迭代器
	for{
		block:=bci.next()//取得下一个区块
		fmt.Printf("上一块hash:%x",block.PrevBlockHash)
		fmt.Println("\n")
		fmt.Printf("当前hash:%x",block.Hash)
		fmt.Println("\n")
		pow:=NewProofOfWork(block)
		fmt.Printf("pow: %s",strconv.FormatBool(pow.Validate()))
		fmt.Println("\n")
		fmt.Println("\n")
		
		// 	fmt.Printf("数据:%s",block.Data)
		// 	fmt.Println("\n")

		if len(block.PrevBlockHash)==0{ //遇到创世区块
			break
		}
	}
}

func (cli *CLI) send (from ,to string,amount int)  {
	bc:=NewBlockchain(from)
	defer bc.DB.Close()
	tx :=NewUTXOTransaction(from,to,amount,bc) //转账
	bc.MineBlock([]*Transaction{tx})  //挖矿确认交易
	fmt.Println("交易成功")
}

func (cli *CLI)Run()  {
	cli.validateArgs() //校验
	//处理命令行参数
	getbalancecmd:=flag.NewFlagSet("getbalance",flag.ExitOnError)
	createblockchaincmd:=flag.NewFlagSet("createblockchain",flag.ExitOnError)
	sendcmd:=flag.NewFlagSet("send",flag.ExitOnError)
	// addblockcmd:=flag.NewFlagSet("addblock",flag.ExitOnError)
	showchaincmd:=flag.NewFlagSet("showchain",flag.ExitOnError)

	getbalanceaddress:=getbalancecmd.String("address","","查询地址")
	createblockchainaddress:=createblockchaincmd.String("address","","地址")
	sendfrom:=sendcmd.String("from","","谁给的")
	sendto:=sendcmd.String("to","","给谁的")
	sendamount:=sendcmd.Int("amount",0,"金额")

	switch os.Args[1] {
	case "getbalance":
		err:=getbalancecmd.Parse(os.Args[2:]) //解析参数
		if err!=nil{
			log.Panic(err)
		}
	case "createblockchain":
		err:=createblockchaincmd.Parse(os.Args[2:]) //解析参数
		if err!=nil{
			log.Panic(err)
		}
	case "send":
		err:=sendcmd.Parse(os.Args[2:]) //解析参数
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
	// if addblockcmd.Parsed(){
	// 	if *addBlockData==""{
	// 		addblockcmd.Usage()
	// 		os.Exit(1)
	// 	}else{
	// 		cli.addBlock(*addBlockData) //增加区块
	// 	}
	// }
	if getbalancecmd.Parsed(){
		if *getbalanceaddress==""{
			getbalancecmd.Usage()
					os.Exit(1)
				}else{
					cli.getBalance(*getbalanceaddress) //增加区块
				}
	}
	if createblockchaincmd.Parsed(){
		if *createblockchainaddress==""{
			createblockchaincmd.Usage()
					os.Exit(1)
				}else{
					cli.createBlockChain(*createblockchainaddress) //增加区块
				}
	}
	if sendcmd.Parsed(){
		if *sendfrom=="" || *sendto=="" || *sendamount<=0{
			sendcmd.Usage()
					os.Exit(1)
				}else{
					cli.send(*sendfrom,*sendto,*sendamount) //增加区块
				}
	}
	if showchaincmd.Parsed(){
		cli.showBlockchain() //显示区块链
	}
	
}