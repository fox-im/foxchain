package core
import(
	"log"
	"fmt"
	"github.com/boltdb/bolt"
	"encoding/hex"
	"os"
)


const dbFile="blockchain.db" //数据库文件名目录
const blockBucket="blocks" //名称
const genesisCoinbaseData="sssssdkdk"

/*
	结构体定义区块 Blockchain
*/
type Blockchain struct{
	// Blocks []*Block   //一个存储Block指针地址的数组,
	Tip []byte //二进制数组
	DB *bolt.DB //数据库
}

type BlockchainIterator struct{
	 currentHash []byte  //当前的hash
	 db *bolt.DB //数据库
}

//挖矿带来的交易
func (blockchain *Blockchain)MineBlock(transactions []*Transaction){
	var Lasthash []byte //最后的哈希
	err:=blockchain.DB.View(func (tx *bolt.Tx) error{
		bucket:=tx.Bucket([]byte (blockBucket))
		// lastHash=bucket.Get([]byte("1"))  //取出最后区块Hash
		Lasthash=bucket.Get([]byte("1"))
		return nil
	})
	if err!=nil{
		log.Panic(err) 
	}

	newBlock:=NewBlock(Lasthash,transactions) //创建一个新的区块
	err=blockchain.DB.Update(func (tx *bolt.Tx)error{
		bucket:=tx.Bucket([]byte (blockBucket))
		err:=bucket.Put(newBlock.Hash,newBlock.Serialize())  //存入
		if err!=nil{
			log.Panic(err) 
		}
		err=bucket.Put([]byte("1"),newBlock.Hash)  //存入
		if err!=nil{
			log.Panic(err) 
		}
		blockchain.Tip=newBlock.Hash //保存上一块的hash
		return nil
	})
	if err!=nil{
		log.Panic(err) 
	}

}

//获取没使用输出的交易列表
func (blockchain *Blockchain) FindUnspentTransactions(address string) []Transaction{
	var unspentTXs [] Transaction //交易实务
	spentTXOS:=make(map[string][]int) //开辟内存
	bci:=blockchain.Iterator() //迭代器
	for{
		block:=bci.next()
		for _,tx:=range block.Transaction{ //循坏每一个交易
			txID:=hex.EncodeToString(tx.ID) //获取交易编号
			Outputs:
			for outindex,out:=range tx.Vout{
				if spentTXOS[txID]!=nil{
					for _,spentOut:=range spentTXOS[txID]{
						if spentOut==outindex{
							continue Outputs //循坏到不等为止
						}
					}
				}
				if out.CanBeUnlockedWith(address){
					unspentTXs=append(unspentTXs,*tx) //加入列表
				}
			}
			if tx.IsCoinBase()==false{
				for _,in:=range tx.Vin{
					if in.CanUnlockOutPutWith(address){ //判断是否可以绑定
						inTxID:=hex.EncodeToString(in.Txid) //
						spentTXOS[inTxID]=append(spentTXOS[inTxID],in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash)==0{ //最后一块，跳出
			break
		}
	}
	return unspentTXs
}
//获取所有没有使用的交易
func (blockchain *Blockchain) FindUTXO(address string)[]TXOutput{
	var UTXOs []TXOutput
	unspentTransactions:=blockchain.FindUnspentTransactions(address) //查找所有的
	for _,tx:=range unspentTransactions{ //循环所有的交易
		for _,out:=range tx.Vout{
			if out.CanBeUnlockedWith(address){ //判断是否锁定
				UTXOs=append(UTXOs,out) //加入数据
			}
		}
	}
	return UTXOs
}
//获取没有使用的输出以参考输入

func (blockchain *Blockchain) FindSpendableOutputs(address string,amount int)(int,map[string][]int){
	unspentOutputs:=make(map[string][]int) //输出
	unspentTXs:=blockchain.FindUnspentTransactions(address)  //根据地质查找所有的交易
	accmulated:=0 //累计
	Work:
		for _,tx:=range unspentTXs{
			txID:=hex.EncodeToString(tx.ID) //
			for outindex,out:=range tx.Vout{
				if out.CanBeUnlockedWith(address) && accmulated<amount{
					accmulated+=out.Value   //统计金额
					unspentOutputs[txID]=append(unspentOutputs[txID],outindex) //序列叠加
					if accmulated>=amount{
						break Work
					}
				}
			}
		}
	return  accmulated,unspentOutputs
}

//将区块添加到链上
// func (block *Blockchain) AddBlock(data string){
// 		var lastHash []byte  //上一块hash
// 		err:=block.DB.View(func (tx *bolt.Tx) error{
// 			block:=tx.Bucket([]byte(blockBucket))  //取得数据
// 			lastHash=block.Get([]byte("1")) //取得第一块
// 			return nil
// 		})
// 		if err!=nil{
// 			log.Panic(err)  //处理打开错误
// 		}
// 		newBlock:=NewBlock(lastHash,data) //创建一个新的区块
// 		err=block.DB.Update(func (tx *bolt.Tx) error{
// 			bucket :=tx.Bucket([]byte(blockBucket))  //取出
// 			err:=bucket.Put(newBlock.Hash,newBlock.Serialize()) //压入数据
// 			if err!=nil{
// 				log.Panic(err)  //处理压入错误
// 			}
// 			err=bucket.Put([]byte("1"),newBlock.Hash) //压入数据
// 			if err!=nil{
// 				log.Panic(err)  //处理压入错误
// 			}
// 			block.Tip=newBlock.Hash
// 			return nil
// 		})
// }

//迭代器
func (block *Blockchain) Iterator() *BlockchainIterator{
	bcit:=&BlockchainIterator{currentHash:block.Tip,db:block.DB}  //根据区块链创建区块链迭代器
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

//判断数据库是否存在
func dbExists() bool{
	if _,err:=os.Stat(dbFile);os.IsNotExist(err){
		return false
	}
	return true
}


//新建一个区块链
func NewBlockchain(address string) *Blockchain{
	if dbExists()==false{
		fmt.Println("数据库不存在，创建一个")
		os.Exit(1)
	}
	// fmt.Println("开始")
	var tip []byte //存储区块链的二进制数据
	db,err:=bolt.Open(dbFile,0600,nil) //打开数据库
	if err!=nil{
		log.Panic(err)  //处理数据库打开错误
	}
	err=db.Update(func (tx *bolt.Tx) error{
		bucket:=tx.Bucket([]byte(blockBucket)) //按照名称打开数据库的表格
		tip=bucket.Get([]byte("1")) //
		// if bucket==nil{
		// 	fmt.Println("*******当前数据库没有区块链，创建一个新的")
		// 	genesis:=GenerateGenesisBlock() //创建创世区块
		// 	bucket,err:=tx.CreateBucket([]byte(blockBucket)) //创建一个数据库的表格
		// 	if err!=nil{
		// 		log.Panic(err)  //处理数据库表格创建错误
		// 	}
		// 	err=bucket.Put(genesis.Hash,genesis.Serialize()) //存入数据
		// 	if err!=nil{
		// 		log.Panic(err)  //处理数据库数据存入错误
		// 	}
		// 	err=bucket.Put([]byte("1"),genesis.Hash) //存入数据
		// 	if err!=nil{
		// 		log.Panic(err)  //处理数据库数据存入错误
		// 	}
		// 	tip=genesis.Hash //取得Hash
		// }else{
		// 	//创建区块


		// 	tip=bucket.Get([]byte("1"))
		// }
		return nil
	}) //更新数据
	if err!=nil{
		log.Panic(err)  //处理数据库更新错误
	}
	
	bc:=Blockchain{Tip:tip,DB:db}

	return &bc
}

func createBlockChain(address string)*Blockchain{
	if dbExists()==false{
		fmt.Println("数据库已经存在，无需创建")
		os.Exit(1)
	}
	var tip []byte //存储区块链的二进制数据
	db,err:=bolt.Open(dbFile,0600,nil) //打开数据库
	if err!=nil{
		log.Panic(err)  //处理数据库打开错误
	}
	err=db.Update(func (tx *bolt.Tx) error{
		cbtx:=NewCoinBaseTX(address,genesisCoinbaseData) //创建创世区块的事务交易
		genesis:=GenerateGenesisBlock(cbtx)  //创建创世区块
		// bucket,err:=tx.Bucket([]byte(blockBucket)) //按照名称打开数据库的表格
		bucket,err:=tx.CreateBucket([]byte(blockBucket))
		if err!=nil{
			log.Panic(err)  //处理数据库打开错误
		}
		err=bucket.Put(genesis.Hash,genesis.Serialize()) //存储
		if err!=nil{
			log.Panic(err)
		}
		err=bucket.Put([]byte("1"),genesis.Hash)  //记录最后一个区块的Hash
		if err!=nil{
			log.Panic(err)
		}
		tip=genesis.Hash
		return nil
	})


	bc:=Blockchain{Tip:tip,DB:db}
	return &bc
}

// func FindSpendableOutPuts(){

// }

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