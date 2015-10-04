package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

var CurrentStockNames [5]string

var CurrentStockValues [5]float64

type Server struct{}

var TradeProfile = make(map[int]StockResObj)

type MyJsonName struct {
	Query struct {
		Count       int    `json:"count"`
		Created     string `json:"created"`
		Diagnostics struct {
			Build_version string `json:"build-version"`
			Cache         struct {
				Content              string `json:"content"`
				Execution_start_time string `json:"execution-start-time"`
				Execution_stop_time  string `json:"execution-stop-time"`
				Execution_time       string `json:"execution-time"`
				Method               string `json:"method"`
				Type                 string `json:"type"`
			} `json:"cache"`
			Javascript struct {
				Execution_start_time string `json:"execution-start-time"`
				Execution_stop_time  string `json:"execution-stop-time"`
				Execution_time       string `json:"execution-time"`
				Instructions_used    string `json:"instructions-used"`
				Table_name           string `json:"table-name"`
			} `json:"javascript"`
			PubliclyCallable string `json:"publiclyCallable"`
			Query            struct {
				Content              string `json:"content"`
				Execution_start_time string `json:"execution-start-time"`
				Execution_stop_time  string `json:"execution-stop-time"`
				Execution_time       string `json:"execution-time"`
				Params               string `json:"params"`
			} `json:"query"`
			Service_time string `json:"service-time"`
			URL          []struct {
				Content              string `json:"content"`
				Execution_start_time string `json:"execution-start-time"`
				Execution_stop_time  string `json:"execution-stop-time"`
				Execution_time       string `json:"execution-time"`
			} `json:"url"`
			User_time string `json:"user-time"`
		} `json:"diagnostics"`
		Lang    string `json:"lang"`
		Results struct {
			Quote struct {
				LastTradePriceOnly   string `json:"LastTradePriceOnly"`
				MarketCapitalization string `json:"MarketCapitalization"`
				Name2                string `json:"Name"`
				Name                 string `json:"symbol"`
			} `json:"quote"`
		} `json:"results"`
	} `json:"query"`
}

type StockReqObj struct {
	Name       [5]string
	Percentage [5]int
	Budget     float32
	TradeId    int
}

type StockResObj struct {
	TradeId            int
	Name               [5]string
	NumberOfStocks     [5]int
	StockValue         [5]float64
	UnvestedAmount     float64
	CurrentMarketValue float64
	ProfitLoss         [5]string
}

var names [5]string
var val [5]float64

func (this *Server) Negate(i int64, reply *int64) error {
	*reply = -i
	return nil
}

var tradeId int

func (this *Server) Receive(Sr1 StockReqObj, Sresp *StockResObj) error {


	for index := 0; index < len(Sr1.Name); index++ {
		if Sr1.Name[index] != "" {


			selectQuery := "https://query.yahooapis.com/v1/public/yql?q=select%20LastTradePriceOnly%2C%20Symbol%20from%20yahoo.finance.quote%20"
			whereQuery := "where%20symbol%20in%20("
			endQuery := ")&format=json&diagnostics=true&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="

			whereQuery = whereQuery + "%27" + Sr1.Name[index] + "%27"
			finalQuery := selectQuery + whereQuery + endQuery
			res, err := http.Get(finalQuery)

			fmt.Println(finalQuery)
			if err != nil {
				log.Fatal(err)
			}
			robots, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			var myjson MyJsonName
			err = json.Unmarshal(robots, &myjson)
			fmt.Println(myjson.Query.Results.Quote.Name)
			fmt.Println(myjson.Query.Results.Quote.LastTradePriceOnly)

			names[index] = myjson.Query.Results.Quote.Name
			val[index], err = strconv.ParseFloat(myjson.Query.Results.Quote.LastTradePriceOnly, 64)

			Sresp.Name[index] = names[index]
			fmt.Println("Stock name ", names[index])
			fmt.Println("Stock value ", val[index])
		}
	}
	var amountLeft float64
	amountLeft = 0

	for NewIndex := 0; NewIndex < len(Sr1.Name); NewIndex++ {
		if Sr1.Name[NewIndex] != "" {

			AllocatedAmount := float64((Sr1.Budget * float32(Sr1.Percentage[NewIndex])) / 100)
			fmt.Println("AllocatedAmount", AllocatedAmount)

			Sresp.NumberOfStocks[NewIndex] = int(AllocatedAmount / val[NewIndex])
			fmt.Println("Sresp.NumberOfStocks[NewIndex]", Sresp.NumberOfStocks[NewIndex])
			var tempSum float64

			var stValue float64
			stValue = float64(val[NewIndex]) * float64(Sresp.NumberOfStocks[NewIndex])
			tempSum = float64(AllocatedAmount - stValue)
			Sresp.StockValue[NewIndex] = stValue
			fmt.Println("tempSum", tempSum)
			amountLeft += tempSum
			fmt.Println("number of stocks", Sresp.NumberOfStocks[NewIndex])
		}
	}
	fmt.Println("Amount left: ", amountLeft)
	Sresp.UnvestedAmount = amountLeft
	fmt.Println(Sresp.UnvestedAmount)

	tradeId += 1
	Sresp.TradeId = tradeId
	TradeProfile[Sresp.TradeId] = *Sresp
	fmt.Println(TradeProfile[Sresp.TradeId])
	Test := TradeProfile[Sresp.TradeId]
	fmt.Println(Test.CurrentMarketValue)


	return nil
}

//start here

func (this *Server) GetTradeProfile(Sr1 StockReqObj, Sresp *StockResObj) error {
	fmt.Println("GetTradeProfile*******************************************************************")
	fmt.Println("Sr1.TradeId ", Sr1.TradeId)
	
	
	Test := TradeProfile[Sr1.TradeId]
	fmt.Println(Test)

	var StockNames [5]string
	var NumberOfStocks [5]int
	var StockValues [5]float64

	for index := 0; index < len(Test.Name); index++ {
		StockNames[index] = Test.Name[index]
		NumberOfStocks[index] = Test.NumberOfStocks[index]
		StockValues[index] = Test.StockValue[index]
	}

	//start
	for index := 0; index < len(StockNames); index++ {
		if StockNames[index] != "" {

	
			selectQuery2 := "https://query.yahooapis.com/v1/public/yql?q=select%20LastTradePriceOnly%2C%20Symbol%20from%20yahoo.finance.quote%20"
			whereQuery2 := "where%20symbol%20in%20("
			endQuery2 := ")&format=json&diagnostics=true&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="

			whereQuery2 = whereQuery2 + "%27" + StockNames[index] + "%27"
			finalQuery2 := selectQuery2 + whereQuery2 + endQuery2
			res, err := http.Get(finalQuery2)

			fmt.Println(finalQuery2)
			if err != nil {
				log.Fatal(err)
			}
			robots2, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			var myjson MyJsonName
			err = json.Unmarshal(robots2, &myjson)
			fmt.Println("New Name: ", myjson.Query.Results.Quote.Name)
			fmt.Println("New Price: ", myjson.Query.Results.Quote.LastTradePriceOnly)

			CurrentStockNames[index] = myjson.Query.Results.Quote.Name
			CurrentStockValues[index], err = strconv.ParseFloat(myjson.Query.Results.Quote.LastTradePriceOnly, 64)

			Sresp.Name[index] = names[index]
			Sresp.TradeId = Sr1.TradeId
			fmt.Println()
			fmt.Println("New Stock name ", CurrentStockNames[index])
			fmt.Println("New Stock value ", CurrentStockValues[index])

		} 
	} 

	
	Sresp.CurrentMarketValue = 0
	for NewIndex2 := 0; NewIndex2 < len(CurrentStockNames); NewIndex2++ {
		if CurrentStockNames[NewIndex2] != "" {

			Sresp.CurrentMarketValue += (CurrentStockValues[NewIndex2] * float64(Test.NumberOfStocks[NewIndex2]))
			fmt.Println("Sresp.CurrentMarketValue", Sresp.CurrentMarketValue)
			Sresp.NumberOfStocks[NewIndex2] = Test.NumberOfStocks[NewIndex2]
			Sresp.StockValue[NewIndex2] = (CurrentStockValues[NewIndex2] * float64(Test.NumberOfStocks[NewIndex2]))
			var TestProfitLoss float64
			TestProfitLoss = (CurrentStockValues[NewIndex2] * float64(Test.NumberOfStocks[NewIndex2])) - (StockValues[NewIndex2] * float64(Test.NumberOfStocks[NewIndex2]))
			if TestProfitLoss > 0 {
				Sresp.ProfitLoss[NewIndex2] = " + "
			} else if TestProfitLoss > 0 {
				Sresp.ProfitLoss[NewIndex2] = " - "
			} else {
				Sresp.ProfitLoss[NewIndex2] = " NoChange "
			} 
		} 
	} 
	Sresp.UnvestedAmount = Test.UnvestedAmount

	fmt.Println("New Unvested Amount", Sresp.UnvestedAmount)
	
	fmt.Println("Current total value of stocks", Sresp.CurrentMarketValue)

	return nil
}

//end here

func server() {
	rpc.Register(new(Server))
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(c)
	}
}

func main() {
	go server()
	
	var input string
	fmt.Scanln(&input)
}
