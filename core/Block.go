package core
import(
	// "crypto/sha256"
	// "encoding/hex"
	"time"
	"encoding/gob"
	"log"
	// "strconv"
	"bytes"
	// "math/god"
)
/*
	结构体定义区块 Block
*/
type Block struct{
	Index int64 //区块编号
	Timestamp int64 //时间戳
	PrevBlockHash []byte //上一个区块的哈希值
	Hash []byte //当前区块哈希值
	Data []byte //区块数据
	Nonce int //工作量证明
}

// //生成hash
// func calculateHash(b Block) string {
// 	blockData := string(b.Index) +string(b.Timestamp)+string(b.PrevBlockHash)+string(b.Data)
// 	hashInBytes := sha256.Sum256([]byte(blockData)) 
// 	return hex.EncodeToString(hashInBytes[:])
// }

// func (block *Block)SetHash(){
// 	//处理当前的时间，转为10进制的字符串，再转化为字节集合
// 	timestamp:=[]byte(strconv.FormatInt(block.Timestamp,10))
// 	//叠加需要去hash的数据
// 	headers:=bytes.Join([][]byte{block.PrevBlockHash,block.Data,timestamp},[]byte{})
// 	//计算出hash
// 	hash:=sha256.Sum256(headers)
// 	//赋值
// 	block.Hash=hash[:] 
// }

//生成一个区块
// func GenerateNewBlock(prevBlock Block,data string ) *Block{
// 	newBlock :=Block{}
// 	newBlock.Index=prevBlock.Index+1
// 	newBlock.PrevBlockHash=prevBlock.Hash
// 	newBlock.Timestamp=time.Now().Unix()
// 	newBlock.Data=data
// 	newBlock.Hash=calculateHash(newBlock)
// 	return newBlock
// }
func NewBlock(prevBlockHash []byte,data string ) *Block{
	//newBlock是指针，指向初始化对象
	newBlock :=&Block{Timestamp:time.Now().Unix(),Data:[]byte(data),PrevBlockHash:prevBlockHash,Hash:[]byte{},Nonce:0}
	// newBlock.SetHash()
	pow:=NewProofOfWork(newBlock) //挖矿附加这个区块
	nonce,hash:=pow.Run() //开始挖矿
	newBlock.Hash=hash[:]
	newBlock.Nonce=nonce
	return newBlock
}
//创世区块
func GenerateGenesisBlock() *Block{
	// prevBlock := Block{}
	// prevBlock.Index=-1
	// prevBlock.Hash=[]byte()
	// return GenerateNewBlock(prevBlock,"sss")
	return NewBlock([]byte{},"First Block")
}

//序列化,把对象转为二进制字节集，可以写入文件
func (block *Block) Serialize() []byte{
	var result bytes.Buffer  //开辟内存，存放字节集合
	encoder:=gob.NewEncoder(&result)  //编码
	err:=encoder.Encode(block)
	if err!=nil{
		log.Panic(err)
	}
	return result.Bytes()

}

//二进制字节集转化为对象
func DeserializeBlock(data []byte) *Block{
	var block Block //存储用于字节转化的对象
	decoder:=gob.NewDecoder(bytes.NewReader(data))  //解码
	err:=decoder.Decode(&block) //尝试解码
	if err!=nil{
		log.Panic(err) //处理错误
	}
	return &block
}