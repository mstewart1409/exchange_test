package main

import (
    "fmt"
)

var (
    userBalances = map[string]float64{
        "BTC":  5.0,
        "ETH":  10.0,
        "USDT": 1000.0,
    }
    
    orderBook = map[string]map[string][]map[string]float64{
        "BTC-ETH": {
            "bids": {{"price": 0.03, "amount": 10}, {"price": 0.029, "amount": 5}},
            "asks": {{"price": 0.032, "amount": 7}, {"price": 0.033, "amount": 12}},
        },
    }
)

func executeOrder(orderType, baseCurrency, quoteCurrency string, amount, limitPrice, stopPrice, price float64) (map[string]float64, map[string]map[string][]map[string]float64) {
    orderTypeLogic := map[string]map[string]string{
        "buy": {
            "balanceKey":           quoteCurrency,
            "otherBalanceKey":      baseCurrency,
            "orderBookKey":         fmt.Sprintf("%s-%s", baseCurrency, quoteCurrency),
            "orderBookAskPriceKey": "price",
            "orderBookBidPriceKey": "price",
        },
        "sell": {
            "balanceKey":           baseCurrency,
            "otherBalanceKey":      quoteCurrency,
            "orderBookKey":         fmt.Sprintf("%s-%s", baseCurrency, quoteCurrency),
            "orderBookAskPriceKey": "price",
            "orderBookBidPriceKey": "price",
        },
    }
    
    if orderType == "limit" {
        if limitPrice == 0 {
            fmt.Println("Limit order price cannot be zero")
            return userBalances, orderBook
        }
        
        if orderTypeLogic[orderType]["balanceKey"] == baseCurrency {
            price = limitPrice
        }
    }
    
    if orderType == "stop" {
        if stopPrice == 0 {
            fmt.Println("Stop order price cannot be zero")
            return userBalances, orderBook
        }
        
        if orderTypeLogic[orderType]["balanceKey"] == quoteCurrency {
            price = stopPrice
        }
    }
    
    if orderType == "market" {
        if orderTypeLogic[orderType]["balanceKey"] == quoteCurrency {
            bestAskPrice := orderBook[orderTypeLogic[orderType]["orderBookKey"]]["asks"][0]["price"]
            price = bestAskPrice
        } else {
            bestBidPrice := orderBook[orderTypeLogic[orderType]["orderBookKey"]]["bids"][0]["price"]
            price = bestBidPrice
        }
    }
    
    orderExecuted := false
    
    for i, order := range orderBook[orderTypeLogic[orderType]["orderBookKey"]][orderTypeLogic[orderType]["otherOrderBookKey"]] {
        if (orderType == "buy" && order[orderTypeLogic[orderType]["orderBookAskPriceKey"]] > price) || (orderType == "sell" && order[orderTypeLogic[orderType]["orderBookBidPriceKey"]] < price) {
            break
        }
        
        if order["amount"] >= amount {
            orderValue := amount * price
            
            userBalances[orderTypeLogic[orderType]["balanceKey"]] -= orderValue
            userBalances[orderTypeLogic[orderType]["otherBalanceKey"]] += amount
            
            orderBook[orderTypeLogic[orderType]["orderBookKey"]][orderTypeLogic[orderType]["otherOrderBookKey"]][i]["amount"] -= amount
            
            if orderBook[orderTypeLogic[orderType]["orderBookKey"]][orderTypeLogic[orderType]["otherOrderBookKey"]][i]["amount"] == 0 {
                orderBook[orderTypeLogic[orderType]["orderBookKey"]][orderTypeLogic[orderType]["otherOrderBookKey"]] = append(orderBook[orderTypeLogic[orderType]["orderBookKey"]][orderTypeLogic[orderType]["otherOrderBookKey"]][:i], orderBook[orderTypeLogic[orderType]["orderBookKey"]][orderTypeLogic[orderType]["otherOrderBookKey"]][i+1:]...)
            }
            
            orderExecuted = true
            break
        } else {
            orderValue := order["amount"] * price
            
            userBalances[orderTypeLogic[orderType]["balanceKey"]] -= orderValue
            userBalances[orderTypeLogic[orderType]["otherBalanceKey"]] += order["amount"]
            
            amount -= order["amount"]
            
            orderBook[orderTypeLogic[orderType]["orderBookKey"]][orderTypeLogic[orderType]["otherOrderBookKey"]][i]["amount"] = 0
            
            if orderExecuted == false {
                orderExecuted = true
            }
        }
    }
    
    if orderExecuted == false {
        fmt.Println("Order was not executed")
    }
    
    return userBalances, orderBook
}

func main() {
    userBalances, orderBook = executeOrder("limit", "BTC", "ETH", 3.0, 0.028, 0, 0)
    fmt.Println(userBalances, orderBook)
}

