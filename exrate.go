package main

import (
	"fmt"
	"os"
	"net/http"
	"io"
	"encoding/json"
	"strconv"
	"strings"
	"sort"
)

type ExchangeRates struct{
	Time string `json:"time_last_update_utc"`
	Rates map[string] float64 `json:"rates"`
}

func validateArgs(supported_currencies []string) (float64, string, []string, error){
	// check required args
	if len(os.Args) < 3{
		return 0.0, "", nil, fmt.Errorf("Error: Not enough arguments")
	}

	// convert amount to float64
	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil{
		fmt.Printf("Error getting currency amount: %s\n", err)
		return 0.0, "", nil, err
	}

	// make currencis uppercase and check support
	currency := strings.ToUpper(os.Args[2])
	currencySupported := false
	for _, supported_currency := range supported_currencies{
		if currency == supported_currency{
			currencySupported = true
		}
	}
	if !currencySupported{
		return 0.0, "", nil, fmt.Errorf("Currency %s not supported\n", currency)
	}

	conversions := os.Args[3:]
	for _, conversion := range conversions {
		conversionSupported := false
		conversion = strings.ToUpper(conversion)
		for _, supported_currency := range supported_currencies{
			if conversion == supported_currency{
				conversionSupported = true
			}
		}
		if !conversionSupported{
			return 0.0, "", nil, fmt.Errorf("Currency %s not supported\n", conversion)
		}
	}

	return amount, currency, conversions, nil
}

func updateExchangeRates(url string) error{
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("API connection error: %s\n", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Can't read API response: %s\n", err)
		return err
	}

	file, err := os.Create("exrate.json")
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		fmt.Printf("Error writing to file: %s\n", err)
		return err
	}

	return nil
}

func getExchangeRates() (*ExchangeRates, error){
	file, err := os.Open("exrate.json")
	if err != nil{
		fmt.Printf("Error opening file: %s\n", err)
		return nil, err
	}
	defer file.Close()

	var exRate ExchangeRates

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&exRate)
	if err != nil{
		fmt.Printf("Error reading json: %s\n", err)
		return nil, err
	}
	return &exRate, nil
}

func getListOfSupportedCurrencies(exRate *ExchangeRates) []string{
	keys := make([]string, 0, len(exRate.Rates))
	for k := range exRate.Rates {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func printHelp(supported_currencies []string){
	fmt.Printf("Usage:\texrate <amount> <currency_from> [<currency_to> ...]\nExamples:\n\texrate 10 USD\n\texrate 150 BYN KZT\n\texrate 4900 KZT USD RUB BYN\nSupported currencies: %v", supported_currencies)
}

func main(){

	err := updateExchangeRates("https://open.er-api.com/v6/latest/USD")
	if err != nil{
		fmt.Printf("Will try to use information from old requests... (Offline mode)\n")
	}

	exRate, err := getExchangeRates()
	if err != nil{
		fmt.Printf("All hope is lost. You NEED Internet :(\n")
		return
	}

	supported_currencies := getListOfSupportedCurrencies(exRate)

	amount, currency, conversions, err := validateArgs(supported_currencies)
	if err != nil{
		fmt.Printf("%s\n", err)
		printHelp(supported_currencies)
		return
	}

	fmt.Printf("\n%s\n\n", exRate.Time)

	if len(conversions) < 1{
		for _, supported_currency := range supported_currencies{
			fmt.Printf("%s: %.2f\n", supported_currency, exRate.Rates[supported_currency] * amount / exRate.Rates[currency])
		}
	} else {
		fmt.Printf("%s: %.2f\n", currency, amount)
		for _, conversion := range conversions{
			conversion = strings.ToUpper(conversion)
			fmt.Printf("%s: %.2f\n", conversion, exRate.Rates[conversion] * amount / exRate.Rates[currency])
		}
	}

}
