package core
import(
	// "log"
	// "fmt"
)
/*
	结构体定义区块 Blockchain
*/
type Blockchain struct{
	Blocks []*Block   //一个存储Block指针地址的数组,
}

//创建新的区块链
func NewBlockchain() *Blockchain{
	return &Blockchain{Blocks:[]*Block{GenerateGenesisBlock()}}
}


//将区块添加到链上 接口
func (blocks *Blockchain) AddBlock(data string){
		prevBlock:=blocks.Blocks[len(blocks.Blocks)-1]  //取出最后一个区块
		newBlock:=NewBlock(prevBlock.Hash,data) //创建一个新的区块
		blocks.Blocks=append(blocks.Blocks,newBlock)
}


// func (bc *Blockchain)SendData(data string){
// 	preBlock := bc.Blocks[len(bc.Blocks)-1]
// 	newBlock :=GenerateNewBlock(*preBlock,data)
// 	bc.ApendBlock(&newBlock)
// }

// func (bc *Blockchain) ApendBlock(newBlock *Block){
// 	if len(bc.Blocks) ==0{
// 		bc.Blocks=append(bc.Blocks,newBlock)
// 		return 
// 	}
// 	if isValid(*newBlock,*bc.Blocks[len(bc.Blocks)-1]){
// 		bc.Blocks=append(bc.Blocks,newBlock)
// 	}else {
// 		log.Fatal("invalid block")
// 	}
// }

// func (bc *Blockchain) Print(){
// 	for _,block:=range bc.Blocks{
// 		fmt.Print("Index:",block.Index)
// 		fmt.Print("\n")
// 		fmt.Print("Prev.Hash:"+block.PrevBlockHash+"\n")
// 		fmt.Print("Curr.Hash:"+block.Hash+"\n")
// 		fmt.Print("Curr.Data:"+block.Data+"\n")
// 		fmt.Print("Curr.Timestamp:\n",block.Timestamp)
// 		fmt.Print("\n")
// 		fmt.Print("\n")
// 	}
// }

// //校验新的区块
// func isValid(newBlock Block,oldBlock Block) bool{
// 	if newBlock.Index-1 != oldBlock.Index{
// 		return false
// 	}

// 	if newBlock.PrevBlockHash != oldBlock.Hash{
// 		return false
// 	}

// 	if calculateHash(newBlock) != newBlock.Hash{
// 		return false
// 	}

// 	return true
// }