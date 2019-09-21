package core

import(
	"crypto/sha256"
	// "encoding"
	// "time"
	"encoding/hex"
	"encoding/gob"
	"log"
	// "strconv"
	"bytes"
	"fmt"
	// "math/god"
)

const  subsidy=10   //奖励，矿工挖矿的奖励

//输入
type TXInput struct{
	Txid []byte  //交易的Id
	Vout int   //该交易中的output索引
	ScriptSig string  //钱包地址
}

//输出
type TXOutput struct{
	Value int  //币
	ScriptPubkey string  //脚本
}

//交易，编号，输入，输出
type Transaction struct{
	ID []byte
	Vin [] TXInput
	Vout [] TXOutput
}

//判断是否在币的基础上进行交易，检查失误是否为coinbase,挖矿的奖励
func (tx *Transaction) IsCoinBase() bool{
	return len(tx.Vin)==1 && len(tx.Vin[0].Txid)==0 && tx.Vin[0].Vout==-1
}

//设置交易ID,二进制
func (tx *Transaction)SetID(){
	var encoded bytes.Buffer //开辟内存
	var hash[32] byte //hash数组
	enc:=gob.NewEncoder(&encoded) //解码对象
	err:=enc.Encode(tx) //解码
	if err!=nil{
		log.Panic(err)
	}
	hash=sha256.Sum256(encoded.Bytes()) //计算hash
	tx.ID=hash[:]  //设置ID
}

//检查地址是否启动事物
func (input *TXInput) CanUnlockOutPutWith(unlockingData string) bool{
	return input.ScriptSig==unlockingData
}

//是否可以解锁输出
func (output *TXOutput) CanBeUnlockedWith(unlockingData string) bool{
	return output.ScriptPubkey==unlockingData 
}

//挖矿交易
func NewCoinBaseTX(to ,data string)*Transaction{
	if data==""{
		data=fmt.Sprintf("奖励个%s",to)
	}
	txin:=TXInput{Txid:[]byte{},Vout:-1,ScriptSig:data} //输入奖励
	txout:=TXOutput{Value:subsidy,ScriptPubkey:to} //输出奖励
	tx:=Transaction{ID:nil,Vin:[]TXInput{txin},Vout:[]TXOutput{txout}} //交易
	return &tx
}

//转账交易
func NewUTXOTransaction(from,to string,amount int,bc *Blockchain)*Transaction{
	var inputs []TXInput
	var outputs []TXOutput
	acc,validOutputs:=bc.FindSpendableOutputs(from,amount)
	if acc<amount{
		log.Panic("交易金额不足")
	}
	for txid,outs:=range validOutputs{
		txID,err:=hex.DecodeString(txid) //解码
		if err!=nil{
			log.Panic(err) //处理错误
		}
		for _,out:=range outs{
			input:=TXInput{Txid:txID,Vout:out,ScriptSig:from}
			inputs=append(inputs,input) //输出的交易
			// output:=TXOutput(Value:,ScriptPubkey:to)
			
		}
	}
	// Build a list of outputs
	// from := fmt.Sprintf("%s", wallet.GetAddress())
	output:=TXOutput{Value:amount,ScriptPubkey:to}
	outputs = append(outputs,output)
	if acc > amount {
		output=TXOutput{Value:acc-amount,ScriptPubkey:from}
		outputs = append(outputs, output)
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	// UTXOSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}