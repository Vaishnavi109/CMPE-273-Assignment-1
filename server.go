package main

import (
	"encoding/json"
	"strconv"
	"log"
	"net"
	"net/rpc"
	"io/ioutil"
	"strings"
	"math/rand"
	"time"
	"net/rpc/jsonrpc"
	"fmt"
	"net/http"
)
type ShareEngine struct{}
//structure for Request Of Transaction
type ArgsForBuyingStocks struct {
    StockList []string
	StockShare []string
	TransactionBudget float32
}
//structure for Loss/Gain Request
type GainLossRequest struct{
	TradeId int64
}
//structure for Response of Transaction
type ReplyForBoughtStocks struct {
    TradeId int64
	StocksList []string
	Count []int32
	StocksPrice []float32 
	UnvestedAmount float32
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

type OldTradesList struct {
	tradeDetails []*ReplyForBoughtStocks
}

type Response struct {
  List struct {
    Resources []struct {
      Resource struct {
        Fields struct {
          Name    string `json:"name"`
          Price   string `json:"price"`
          Symbol  string `json:"symbol"`
          Ts      string `json:"ts"`
          Type    string `json:"type"`
          UTCTime string `json:"utctime"`
          Volume  string `json:"volume"`
        } `json:"fields"`
      } `json:"resource"`
    } `json:"resources"`
  } `json:"list"`
}

var oldTradeList OldTradesList
//Funtion to get details of Trade
func (s *ShareEngine) RetrieveTradeDetails(args *GainLossRequest, reply *GainLossResponce) error{
	isFound:=false
	fmt.Print("Count of trades: ")
	fmt.Println(len(oldTradeList.tradeDetails))
	for l:=0;l<len(oldTradeList.tradeDetails);l++{
		
		var trade *ReplyForBoughtStocks
		trade = oldTradeList.tradeDetails[l]
		
		if(trade.TradeId==args.TradeId){
			isFound=true
			reply.CurrentMarketValue = 0.0
			for i:=0; i< len(trade.StocksList); i++{	
				resp, err := http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+trade.StocksList[i]+"/quote?format=json")
				if err != nil {
					// handle error
					fmt.Println("Error")
				}
				defer resp.Body.Close()
				
				body, _ := ioutil.ReadAll(resp.Body)
				
				var msg Response
				_ = json.Unmarshal(body, &msg)
				
				reply.StocksList=append(reply.StocksList,msg.List.Resources[0].Resource.Fields.Symbol)
				
				
				stkPrice,_:= strconv.ParseFloat(msg.List.Resources[0].Resource.Fields.Price,64)
				reply.StocksPrice=append(reply.StocksPrice,float32(stkPrice))
				
				if(float32(stkPrice) == trade.StocksPrice[i]){
					reply.StocksGainLoss=append(reply.StocksGainLoss,"") 
				}else if(float32(stkPrice)<trade.StocksPrice[i]){
					reply.StocksGainLoss=append(reply.StocksGainLoss,"-") 
				}else{
					reply.StocksGainLoss=append(reply.StocksGainLoss,"+") 
				}
				reply.Count=append(reply.Count,trade.Count[i])
				
				reply.CurrentMarketValue= reply.CurrentMarketValue + float32(reply.Count[i]) *reply.StocksPrice[i]
			}
			reply.UnvestedAmount= trade.UnvestedAmount
		}
	}
	if(isFound==false){
		var emptyReply *GainLossResponce
		reply = emptyReply
	}
	return nil
}

//Funtion to Buy stock
func (s *ShareEngine) CreateTransactionId(args *ArgsForBuyingStocks, reply *ReplyForBoughtStocks) error{
	shouldBreak:=false
	rand.Seed(time.Now().UTC().UnixNano())
	reply.TradeId = rand.Int63n(9999)
	reply.UnvestedAmount = 0.0
	
	for i:=0; i< len(args.StockList); i++{
		resp, err := http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+args.StockList[i]+"/quote?format=json")
		if err != nil {
			// handle error
			fmt.Println("Error")
		}
		defer resp.Body.Close()
		
		body, _ := ioutil.ReadAll(resp.Body)
		
		var msg Response
		_ = json.Unmarshal(body, &msg)
		
		if(len(msg.List.Resources)!=0){
			reply.StocksList=append(reply.StocksList,msg.List.Resources[0].Resource.Fields.Symbol)
			stkPrice,_:= strconv.ParseFloat(msg.List.Resources[0].Resource.Fields.Price,64)
			reply.StocksPrice=append(reply.StocksPrice,float32(stkPrice))
			stkShare,_:= strconv.ParseFloat(strings.Replace(args.StockShare[i],"%","",-1),64)
			indBudget := (float32(stkShare) * args.TransactionBudget)/100.0
			reply.Count=append(reply.Count,int32(indBudget/float32(reply.StocksPrice[i])))
			individualUnvestedBalance := indBudget - float32(reply.Count[i])* float32(reply.StocksPrice[i])
			
			reply.UnvestedAmount=reply.UnvestedAmount+individualUnvestedBalance
		}else{
			var emptyReply *ReplyForBoughtStocks
			reply = emptyReply
			shouldBreak = true
			break
		}
		
	}
	
	if(shouldBreak==false){
		oldTradeList.tradeDetails = append(oldTradeList.tradeDetails,reply)
		fmt.Println(reply.TradeId)
	}else{
		fmt.Println("No transaction")
	}
	
	
	return nil
}
//funtion to start server
func StartServer() {
	stockEng :=new(ShareEngine)
	server := rpc.NewServer()
	server.Register(stockEng)

    l, e := net.Listen("tcp", ":8222")
    if e != nil {
		log.Fatal("listen error:", e)
    }

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go server.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
}

func main(){
	go StartServer()
	var input string
	fmt.Scanln(&input)
}
