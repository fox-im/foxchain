package core

import(
	"bytes"
	"encoding/binary"
	"log"
)

//整数转化为16进制
func IntToHex(num int64)[]byte  {
	buff:=new(bytes.Buffer) //开辟内存，存储字节集
	err:=binary.Write(buff,binary.BigEndian,num) //num转化为字节集合写入
	if err!=nil{
		log.Panic(err)
	}
	return buff.Bytes()
}