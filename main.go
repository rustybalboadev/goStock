package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/fatih/color"
)

func main() {
	stonks := []string{}

	userDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	makeFiles(userDir)

	parser := argparse.NewParser("goStock", "goStock Commands")
	addStonk := parser.String("a", "add", &argparse.Options{Help: "Add a Stonk Using Its Symbol", Required: false})
	removeStonk := parser.String("r", "remove", &argparse.Options{Help: "Remove a Stonk Using Its Symbol", Required: false})
	p := parser.Parse(os.Args)
	if p != nil {
		fmt.Print(parser.Usage(p))
	}
	if *addStonk != "" {
		stripped := strings.Trim(*addStonk, " ")
		f, err := os.OpenFile(userDir+"\\.goStocks\\config.txt", os.O_APPEND, 0600)
		if err != nil {
			fmt.Println(err)
		}
		f.WriteString(stripped + "\n")
		defer f.Close()
	}

	if *removeStonk != "" {
		stripped := strings.Trim(*removeStonk, " ")
		data, err := ioutil.ReadFile(userDir + "\\.goStocks\\config.txt")
		if err != nil {
			fmt.Println(err)
		}
		splitted := strings.Split(string(data), "\n")
		newData := ""
		for i := 0; i < len(splitted)-1; i++ {
			if splitted[i] != stripped {
				newData += splitted[i] + "\n"
			}
		}
		os.Remove(userDir + "\\.goStocks\\config.txt")
		makeFiles(userDir)
		f, err := os.OpenFile(userDir+"\\.goStocks\\config.txt", os.O_APPEND, 0600)
		if err != nil {
			fmt.Println(err)
		}
		f.WriteString(newData)
	}

	data, err := ioutil.ReadFile(userDir + "\\.goStocks\\config.txt")
	if err != nil {
		fmt.Println(err)
	}

	split := strings.Split(string(data), "\n")
	for i := 0; i < len(split); i++ {
		stonks = append(stonks, split[i])
	}

	for i := 0; i < len(stonks)-1; i++ {
		resp, err := http.Get("https://query1.finance.yahoo.com/v11/finance/quoteSummary/" + stonks[i] + "?modules=summaryDetail,price")
		if err != nil {
			fmt.Println(err)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		body := string(b)
		var t APIResponse
		err = json.Unmarshal([]byte(string(body)), &t)

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

func makeFiles(userDir string) {

	_, makeDir := os.Stat(userDir + "\\.goStocks")
	if os.IsNotExist(makeDir) {
		errDir := os.Mkdir(userDir+"\\.goStocks", 0755)

		exe, err := os.Executable()
		if err != nil {
			fmt.Println(err)
		}
		data, err := ioutil.ReadFile(exe)
		if err != nil {
			fmt.Println(err)
		}
		err = ioutil.WriteFile(userDir+"\\.goStocks\\goStock.exe", data, 0644)
		if err != nil {
			fmt.Println(err)
		}
		os.Setenv("goStock", userDir+"\\.goStocks\\goStock.exe")
		fmt.Println("Your Portfolio Files Have Been Created! Get Help By Using The Argument -h")
		if errDir != nil {
			fmt.Println(errDir)
		}
	}
	_, makeFile := os.Stat(userDir + "\\.goStocks\\config.txt")
	if os.IsNotExist(makeFile) {
		errFile, err := os.Create(userDir + "\\.goStocks\\config.txt")
		if err != nil {
			fmt.Println(errFile)
		}
	}
}
