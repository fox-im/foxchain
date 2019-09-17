package core
import(
	"log"
	"fmt"
	"github.com/boltdb/bolt"
)


const dbFile="blockchain.db" //数据库文件名目录
const blockBucket="blocks" //名称


/*
	结构体定义区块 Blockchain
*/
type Blockchain struct{
	// Blocks []*Block   //一个存储Block指针地址的数组,
	tip []byte //二进制数组
	db *bolt.DB //数据库
}

type BlockchainIterator struct{
	 currentHash []byte  //当前的hash
	 db *bolt.DB //数据库
}

//将区块添加到链上
func (block *Blockchain) AddBlock(data string){
		var lastHash []byte  //上一块hash
		err:=block.db.View(func (tx *bolt.Tx) error{
			block:=tx.Bucket([]byte(blockBucket))  //取得数据
			lastHash=block.Get([]byte("1")) //取得第一块
			return nil
		})
		if err!=nil{
			log.Panic(err)  //处理打开错误
		}
		newBlock:=NewBlock(lastHash,data) //创建一个新的区块
		err=block.db.Update(func (tx *bolt.Tx) error{
			bucket :=tx.Bucket([]byte(blockBucket))  //取出
			err:=bucket.Put(newBlock.Hash,newBlock.Serialize()) //压入数据
			if err!=nil{
				log.Panic(err)  //处理压入错误
			}
			err=bucket.Put([]byte("1"),newBlock.Hash) //压入数据
			if err!=nil{
				log.Panic(err)  //处理压入错误
			}
			block.tip=newBlock.Hash
			return nil
		})
}

//迭代器
func (block *Blockchain) Iterator() *BlockchainIterator{
	bcit:=&BlockchainIterator{currentHash:block.tip,db:block.db}  //根据区块链创建区块链迭代器
	return bcit
}

//根据迭代器取得下一个区块
func (it *BlockchainIterator) next() *Block{
	var block *Block
	err:=it.db.View(func (tx *bolt.Tx) error{
		// block:=tx.Bucket([]byte(blockBucket))  //取得数据
		// lastHash=block.Get([]byte("1")) //取得第一块
		bucket:=tx.Bucket([]byte(blockBucket))
		encodedBlock:=bucket.Get(it.currentHash) //抓取二进制数据
		block=DeserializeBlock(encodedBlock) //解码
		return nil
	})
	if err!=nil{
		log.Panic(err)  //处理打开错误
	}
	it.currentHash=block.PrevBlockHash //哈希赋值
	return block
}

//新建一个区块链
func NewBlockchain() *Blockchain{
	var tip []byte //存储区块链的二进制数据
	db,err:=bolt.Open(dbFile,0600,nil) //打开数据库
	if err!=nil{
		log.Panic(err)  //处理数据库打开错误
	}
	err=db.Update(func (tx *bolt.Tx) error{
		bucket:=tx.Bucket([]byte(blockBucket)) //按照名称打开数据库的表格
		if bucket==nil{
			fmt.Println("当前数据库没有区块链，创建一个新的")
			genesis:=GenerateGenesisBlock() //创建创世区块
			bucket,err:=tx.CreateBucket([]byte(blockBucket)) //创建一个数据库的表格
			if err!=nil{
				log.Panic(err)  //处理数据库表格创建错误
			}
			err=bucket.Put(genesis.Hash,genesis.Serialize()) //存入数据
			if err!=nil{
				log.Panic(err)  //处理数据库数据存入错误
			}
			err=bucket.Put([]byte("1"),genesis.Hash) //存入数据
			if err!=nil{
				log.Panic(err)  //处理数据库数据存入错误
			}
			tip=genesis.Hash //取得Hash
		}else{
			tip=bucket.Get([]byte("1"))
		}
		return nil
	}) //更新数据
	if err!=nil{
		log.Panic(err)  //处理数据库更新错误
	}
	bc:=Blockchain{tip:tip,db:db}
	return &bc
}

// //创建新的区块链
// func NewBlockchain() *Blockchain{
// 	return &Blockchain{Blocks:[]*Block{GenerateGenesisBlock()}}
// }


// //将区块添加到链上 接口
// func (blocks *Blockchain) AddBlock(data string){
// 		prevBlock:=blocks.Blocks[len(blocks.Blocks)-1]  //取出最后一个区块
// 		newBlock:=NewBlock(prevBlock.Hash,data) //创建一个新的区块
// 		blocks.Blocks=append(blocks.Blocks,newBlock)
// }


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