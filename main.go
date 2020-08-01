package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/fatih/color"
)

//Stock struct
type Stock struct {
	Stocks []string
}

var stock Stock

func main() {
	userDir, err := os.UserHomeDir()
	handleErr(err)
	makeFiles()

	parser := argparse.NewParser("goStock", "goStock Commands")
	addStonk := parser.String("a", "add", &argparse.Options{Help: "Add a Stonk Using Its Symbol", Required: false})
	removeStonk := parser.String("r", "remove", &argparse.Options{Help: "Remove a Stonk Using Its Symbol", Required: false})
	repeatStonk := parser.Int("i", "refresh", &argparse.Options{Help: "Refresh Stonks x Amount of Seconds", Required: false})
	var clearBool *bool = parser.Flag("c", "clear", &argparse.Options{Help: "Clear All Stonks", Required: false})
	p := parser.Parse(os.Args)
	handleErr(p)

	if *addStonk != "" {
		stripped := strings.Trim(*addStonk, " ")
		readJSON()
		stock.Stocks = append(stock.Stocks, stripped)
		writeJSON(stock)
	}

	if *removeStonk != "" {
		stripped := strings.Trim(*removeStonk, " ")
		remove(stripped)
	}

	if *clearBool {
		color.Green("All Stonks Have Been Cleared")
		os.Remove(userDir + "\\.goStocks\\config.json")
		makeFiles()
	}

	jsonData := readJSON()
	if *repeatStonk != 0 {
		for {
			showStonks(jsonData)
			time.Sleep(time.Duration(*repeatStonk) * time.Second)
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
	} else {
		showStonks(jsonData)
	}

}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func showStonks(stonks []string) {
	for i := 0; i < len(stonks); i++ {
		resp, err := http.Get("https://query1.finance.yahoo.com/v11/finance/quoteSummary/" + stonks[i] + "?modules=summaryDetail,price")
		handleErr(err)
		b, err := ioutil.ReadAll(resp.Body)
		handleErr(err)
		body := string(b)
		var t APIResponse
		err = json.Unmarshal([]byte(string(body)), &t)
		if t.QuoteSummary.Error.Description != "" {
			remove(stonks[i])
			fmt.Println(color.RedString("\n!!!Stonk Symbol '" + stonks[i] + "' Is Invalid!!!"))
			os.Exit(0)
		}
		symbol := t.QuoteSummary.Result[0].Price.Symbol
		bid := fmt.Sprintf("%.2f", t.QuoteSummary.Result[0].SummaryDetail.Bid.Raw)
		regMarketChange := t.QuoteSummary.Result[0].Price.RegularMarketChange.Raw
		regMarketChangePercent := t.QuoteSummary.Result[0].Price.RegularMarketChangePercent.Raw
		volume := t.QuoteSummary.Result[0].SummaryDetail.Volume.Fmt
		marketChange := fmt.Sprintf("%.2f", regMarketChange)
		marketChangePercent := fmt.Sprintf("%.2f", regMarketChangePercent)

		if regMarketChange < 0 {
			fmt.Printf(color.RedString("%-8s %-8s @ %-8s ▼ %-8s %-8s"), symbol, string(volume), bid, marketChange, marketChangePercent)
			fmt.Print(color.RedString("%\n"))
		} else if regMarketChange == 0 {
			fmt.Printf("%-8s %-8s @ %-8s ▬ %-8s %-8s", symbol, string(volume), bid, marketChange, marketChangePercent)
			fmt.Print("%\n")
		} else {
			fmt.Printf(color.GreenString("%-8s %-8s @ %-8s ▲ %-8s +%-7s"), symbol, string(volume), bid, marketChange, marketChangePercent)
			fmt.Print(color.GreenString("%\n"))
		}
	}
}

func readJSON() []string {

	userDir, err := os.UserHomeDir()
	jsonFile, err := os.Open(userDir + "\\.goStocks\\config.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &stock)
	jsonFile.Close()
	return stock.Stocks

}

func writeJSON(jsonData Stock) {
	userDir, err := os.UserHomeDir()
	handleErr(err)
	file, _ := json.MarshalIndent(jsonData, "", " ")
	_ = ioutil.WriteFile(userDir+"\\.goStocks\\config.json", file, 0644)

}

func makeFiles() {
	userDir, err := os.UserHomeDir()
	handleErr(err)
	_, makeDir := os.Stat(userDir + "\\.goStocks")
	if os.IsNotExist(makeDir) {
		errDir := os.Mkdir(userDir+"\\.goStocks", 0755)
		fmt.Println("Your Portfolio Files Have Been Created! Get Help By Using The Argument -h")
		handleErr(errDir)
	}
	_, makeFile := os.Stat(userDir + "\\.goStocks\\config.json")
	if os.IsNotExist(makeFile) {
		errFile, err := os.Create(userDir + "\\.goStocks\\config.json")
		if err != nil {
			fmt.Println(errFile)
		}
		stockVar := Stock{
			Stocks: []string{},
		}
		file, _ := json.MarshalIndent(stockVar, "", " ")
		_ = ioutil.WriteFile(userDir+"\\.goStocks\\config.json", file, 0644)
	}
}

func remove(stonk string) {
	readJSON()
	userDir, err := os.UserHomeDir()
	handleErr(err)
	newData := []string{}
	for i := 0; i < len(stock.Stocks); i++ {
		if stock.Stocks[i] != stonk {
			newData = append(newData, stock.Stocks[i])
		}
	}
	stock.Stocks = newData
	os.Remove(userDir + "\\.goStocks\\config.json")
	makeFiles()
	writeJSON(stock)
}
