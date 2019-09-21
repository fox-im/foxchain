package main


import(
	"../core"
	// "strconv"
	// "fmt"
) 
func main(){
	// bc:=core.NewBlockchain()
	// bc.SendData("Send 1 BTC to fox1");
	// bc.SendData("Send 1 ETH to fox2");
	// bc.Print()
	// getSumAndSub(1,2);
	// fmt.Println("hello start!")
	// bcc:=core.NewBlockchain() //创建区块链
	// bcc.AddBlock("Send 1 BTC to fox1")
	// bcc.AddBlock("Send 2 BTC to fox1")
	// for _,block :=range bcc.Blocks{
	// 	fmt.Printf("上一块hash%x",block.PrevBlockHash)
	// 	fmt.Println("\n")
	// 	fmt.Printf("数据%s",block.Data)
	// 	fmt.Println("\n")
		// fmt.Printf("当前hash%x",block.Hash)
	// 	pow:=core.NewProofOfWork(block) //校验工作量
	// 	fmt.Printf("pow %s\n",strconv.FormatBool(pow.Validate()))
	// 	fmt.Println("\n")
	// 	fmt.Println("\n")
	// }
	block :=core.NewBlockchain("我是一个地址") //创建区块链
	// fmt.Println("sss")
	// fmt.Printf("%v",block)
	defer block.DB.Close() //延迟关闭数据库
	cli:=core.CLI{Blockchain:block} //创建命令行
	cli.Run() //开启

}


// func getSumAndSub(n1 int,n2 int) (int,int){
// 	sum:=n1+n2
// 	sub:=n1-n2
// 	return sum,sub
// }
