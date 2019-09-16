package main

import "../core"
func main(){
	bc:=core.NewBlockchain()
	bc.SendData("Send 1 BTC to fox1");
	bc.SendData("Send 1 ETH to fox2");
	bc.Print()
}