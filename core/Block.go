package core
import(
	"crypto/sha256"
	"encoding/hex"
	"time"
)
type Block struct{
	Index int64 //区块编号
	Timestamp int64 //时间戳
	PrevBlockHash string //上一个区块的哈希值
	Hash string //当前区块哈希值
	Data string //区块数据
}

//生成hash
func calculateHash(b Block) string {
	blockData := string(b.Index) +string(b.Timestamp)+string(b.PrevBlockHash)+string(b.Data)
	hashInBytes := sha256.Sum256([]byte(blockData)) 
	return hex.EncodeToString(hashInBytes[:])
}

//生成下一个区块
func GenerateNewBlock(prevBlock Block,data string ) Block{
	newBlock :=Block{}
	newBlock.Index=prevBlock.Index+1
	newBlock.PrevBlockHash=prevBlock.Hash
	newBlock.Timestamp=time.Now().Unix()
	newBlock.Data=data
	newBlock.Hash=calculateHash(newBlock)
	return newBlock
}

//创世区块
func GenerateGenesisBlock() Block{
	prevBlock := Block{}
	prevBlock.Index=-1
	prevBlock.Hash=""
	return GenerateNewBlock(prevBlock,"sss")
}