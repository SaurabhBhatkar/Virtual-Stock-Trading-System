package main

import (
	"fmt"
	"net/rpc"
	"strconv"
	"strings"
)

var Sr1 StockReqObj
var Sresp StockResObj

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

func client() {

	

	
	var i3 string
	fmt.Scanln(&i3)
}
func main() {
	

	GetInput()

	
	var input string
	fmt.Scanln(&input)

}

func GetInput() {
	var Sr1 StockReqObj
	var Sresp StockResObj
	var InputFromUser int
	InputFromUser = 0
	fmt.Println(InputFromUser)

	
	fmt.Println("For Buying stocks Enter: 1 or For getting TraderProfile Enter: 2")
	fmt.Scanln(&InputFromUser)

	

	if InputFromUser == 2 {
		fmt.Println("Enter the trade ID")
		Trid := ""
		fmt.Scanln(&Trid)
		NewInt, err := strconv.Atoi(Trid)
		if err != nil {
			fmt.Println("Error occurred. Please enter appropriate trade id")
		}
		Sr1.TradeId = NewInt
		fmt.Println("Sr1.TradeId ", Sr1.TradeId)

	} else if InputFromUser == 1 {
		var BudgetData float32
		fmt.Println("Enter the budget")
		fmt.Scanln(&BudgetData)

		Sr1.Budget = BudgetData

		fmt.Println("Enter the stock values in the format given below : ")
		fmt.Println("StockName1,Percentage1,StockName2,Percentage2,StockName3,Percentage3")
		fmt.Println("Example:")
		fmt.Println("Goog,50,YHOO,50")

		UserString := ""
		fmt.Scanln(&UserString)

		var s []string
		s = strings.Split(UserString, ",")
		Latestindex := 0
		for index := 0; Latestindex < len(s); index++ {
			fmt.Println("index ", index)
			Sr1.Name[index] = s[Latestindex]
			s3, _ := strconv.Atoi(s[Latestindex+1])
			Sr1.Percentage[index] = int(s3)
			fmt.Println("Sr1.Name[index] ", Sr1.Name[index])
			fmt.Println("Sr1.Percentage[index] ", int(Sr1.Percentage[index+1]))
			Latestindex = Latestindex + 2
		}
	} else {
		fmt.Println(InputFromUser == 1)
		fmt.Println(InputFromUser)

		fmt.Println("Nothing found")
	}
	
	fmt.Println("InputFromUser", InputFromUser)
	
	c, err := rpc.Dial("tcp", "127.0.0.1:9999")

	if err != nil {
		fmt.Println(err)
		return
	}

	if InputFromUser == 1 {
		fmt.Println(Sr1.Name[0])
		fmt.Println(Sr1.Percentage[0])

		fmt.Println(Sr1.Name[1])
		fmt.Println(Sr1.Percentage[1])

		fmt.Println(Sr1.Name[2])
		fmt.Println(Sr1.Percentage[2])

		err = c.Call("Server.Receive", Sr1, &Sresp)
		if err != nil {
			fmt.Println("Error", err)
		} else {
			fmt.Println("Trade ID: ", Sresp.TradeId, " and Remaining amount: ", Sresp.UnvestedAmount)
			for index := 0; index < len(Sresp.NumberOfStocks); index++ {
				if Sresp.Name[index] != "" {
					fmt.Println("Name of Stock", Sresp.Name[index], " Number of Stocks", Sresp.NumberOfStocks[index], " Value of Stock", Sresp.ProfitLoss[index], Sresp.StockValue[index])
				}
			}

		} 

	} else if InputFromUser == 2 {
		fmt.Println("GetTradeProfile")

		
		fmt.Println("Second function here")
		err = c.Call("Server.GetTradeProfile", Sr1, &Sresp)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Trade ID: ", Sresp.TradeId, " and Remaining amount: ", Sresp.UnvestedAmount)

		for index := 0; index < len(Sresp.NumberOfStocks); index++ {
			if Sresp.Name[index] != "" {
				fmt.Println("Name of Stock", Sresp.Name[index], " Number of Stocks", Sresp.NumberOfStocks[index], " Value of Stock", Sresp.ProfitLoss[index], Sresp.StockValue[index])
			}
		}

		fmt.Println("Current total value ", Sresp.CurrentMarketValue)
	}
}

func GetTradeProfile() {
	var InputFromUser int
	InputFromUser = 0
	fmt.Println(InputFromUser)
	c, err := rpc.Dial("tcp", "127.0.0.1:9999")

	if err != nil {
		fmt.Println(err)
		return
	}

	if InputFromUser == 2 {
		Sr1.TradeId = 1
		fmt.Println("Second function here")
		err = c.Call("Server.GetTradeProfile", Sr1, &Sresp)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Trade ID: ", Sresp.TradeId, " and Remaining amount: ", Sresp.UnvestedAmount)

		for index := 0; index < len(Sresp.NumberOfStocks); index++ {
			if Sresp.Name[index] != "" {
				fmt.Println("Name of Stock", Sresp.Name[index], " Number of Stocks", Sresp.NumberOfStocks[index], " Value of Stock", Sresp.ProfitLoss[index], Sresp.StockValue[index])
			}
		}

		fmt.Println("Current total value ", Sresp.CurrentMarketValue)
	}

}
