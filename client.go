package main

import (
	"strings"
	"strconv"
	"unsafe"
	"net/rpc/jsonrpc"
	"log"
	"fmt"
    "os"
	"net"	
)

//structure for Response of Transaction
type ReplyForBoughtStocks struct {
    TradeId int64
	StocksList []string
	Count []int32
	StocksPrice []float32 
	UnvestedAmount float32
}

//structure for Loss/Gain Request
type GainLossRequest struct{
	TradeId int64
}

//structure for Loss/Gain Response
type GainLossResponce struct{
	StocksList []string
	Count []int32
	StocksPrice []float32
	StocksGainLoss []string
	UnvestedAmount float32
	CurrentMarketValue float32
}

//structure for Request Of Transaction
type ArgsForBuyingStocks struct {
    StockList []string
	StockShare []string
	TransactionBudget float32
}
//Start Client
func StartClient(CmdArgs []string)  {
	// connect to the server
	//stocksDesc string,budget float64
	connect, err := net.Dial("tcp", "localhost:8222")
	if err != nil {
		panic(err)
	}
	defer connect.Close()
	//create new client
	client := jsonrpc.NewClient(connect)
	
	if(len(CmdArgs)==2){
		s,_ := strconv.ParseFloat(CmdArgs[1],64)
		var replyStock ReplyForBoughtStocks
		var argsStockStruct ArgsForBuyingStocks
		var argsStock *ArgsForBuyingStocks
		argsStockStruct.TransactionBudget = float32(s)
		individualStocks := strings.Split(CmdArgs[0],",")
		
		k:=0
		for i:=0; i < len(individualStocks); i++ {
			individualStockInfo:=strings.Split(individualStocks[i],":")
			for j:=0; j< len(individualStockInfo);j++{
				if(j%2==0){
					argsStockStruct.StockList=append(argsStockStruct.StockList,individualStockInfo[j])
				}else{
					argsStockStruct.StockShare=append(argsStockStruct.StockShare,individualStockInfo[j])
				}
				k++
			}
		}
	
		argsStock = &argsStockStruct
		//Call for Transaction
		err = client.Call("ShareEngine.CreateTransactionId",argsStock,&replyStock)
		if err!=nil{
			log.Fatal("stringed error:", err)
		}
		
		if(unsafe.Sizeof(replyStock)!=0){
			fmt.Printf("TradeId: %d\n",replyStock.TradeId)
			
			for j:=0;j<len(replyStock.StocksList);j++{
				if(j!=len(replyStock.StocksList)-1){
					fmt.Printf("%s:%d:$%.2f,",replyStock.StocksList[j],replyStock.Count[j], replyStock.StocksPrice[j])	
				}else{
					fmt.Printf("%s:%d:$%.2f\n",replyStock.StocksList[j],replyStock.Count[j], replyStock.StocksPrice[j])	
				}
			}
			fmt.Printf("UnvestedAmount: %.2f\n",replyStock.UnvestedAmount)
		}else{
			fmt.Printf("Invalid input!")
		}
			
	}else if(len(CmdArgs)==1){
		
		reqTradeId,_ := strconv.ParseInt(CmdArgs[0],10,64)
		
		var gainLossRequest *GainLossRequest
		var gainLossResponce GainLossResponce 
		gainLossRequest = &GainLossRequest{reqTradeId}
		//Call for Retieving Transaction details
		err = client.Call("ShareEngine.RetrieveTradeDetails",gainLossRequest,&gainLossResponce)
		if err!=nil{
			log.Fatal("stringed error:", err)
		}
		
		if(unsafe.Sizeof(gainLossResponce)==0){
			fmt.Printf("No such transaction!\n")
		}else{
			
			for j:=0;j<len(gainLossResponce.StocksList);j++{
				if(j!=len(gainLossResponce.StocksList)-1){
					fmt.Printf("%s:%d:%s$%.2f,",gainLossResponce.StocksList[j],gainLossResponce.Count[j],gainLossResponce.StocksGainLoss[j],gainLossResponce.StocksPrice[j])	
				}else{
					fmt.Printf("%s:%d:%s$%.2f\n",gainLossResponce.StocksList[j],gainLossResponce.Count[j],gainLossResponce.StocksGainLoss[j], gainLossResponce.StocksPrice[j])	
				}
			}
			fmt.Printf("Current Market Value: %.2f\n",gainLossResponce.CurrentMarketValue)
			fmt.Printf("Unvested Amount: %.2f\n",gainLossResponce.UnvestedAmount)
		}
				
	}
}
func main(){
	cmdArgs:=os.Args[1:]
	StartClient(cmdArgs)
}
