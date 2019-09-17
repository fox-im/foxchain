package core

import(
	"math/big"
	"math"
	// "encoding/binary"
	// "log"
	"bytes"
	"fmt"
	"crypto/sha256"
	// "match/big"
	// "atom"

)

var (
	maxNonce=math.MaxInt64 //最大的64位整数

)

const targetBits=24 //对比的位数，难度

type ProofOfWork struct{
	block *Block    //区块
	target * big.Int  //存储对比的hash的整数

}

//创建一个工作量证明的挖矿对象
func NewProofOfWork(block *Block) *ProofOfWork {
	target:=big.NewInt(1) //初始化目标整数
	target.Lsh(target,uint(256-targetBits)) //数据转换
	pow:=&ProofOfWork{block:block,target:target}  //创建对象
	return pow
}

//准备数据进行挖矿计算
func (pow *ProofOfWork) prepareData(nonce int)[]byte{
	data:=bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,  //上一块的hash
			pow.block.Data,  //当前数据
			IntToHex(pow.block.Timestamp),  //时间转为16进制
			IntToHex(int64(targetBits)), //位数
			IntToHex(int64(nonce)),  //保存工作量的证明
		},[]byte{},
	)
	return data
}

//挖矿执行的过程
func (pow *ProofOfWork) Run() (int,[]byte){
	var hashInt big.Int
	var hash [32]byte
	nonce :=0
	fmt.Printf("当前挖矿计算出来的数据%s",pow.block.Data)
	for nonce<maxNonce{
		data :=pow.prepareData(nonce) //准备好的数据
		hash=sha256.Sum256(data) //计算出hash
		fmt.Printf("\r%x",hash)
		hashInt.SetBytes(hash[:]) //获取要对比的数据
		if hashInt.Cmp(pow.target)==-1{
			break
		}else{
			nonce++
		}
		// fmt.Println("\n\n")
		
	}
	return nonce,hash[:]  //nonce解题的答案   hash当前的hash
}

//校验挖矿是否成功
func (pow *ProofOfWork) Validate() bool{
	var hashInt big.Int
	data :=pow.prepareData(pow.block.Nonce) //准备好的数据
	hash:=sha256.Sum256(data) //计算出hash
	hashInt.SetBytes(hash[:]) //获取要对比的数据
	isValid:=(hashInt.Cmp(pow.target)==-1) //校验数据
	return isValid

}
