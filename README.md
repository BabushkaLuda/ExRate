# ExRate
This CLI application uses <a href="https://www.exchangerate-api.com">This API</a> to show you the exchange rate of your chosen currency.
<br>It makes an API call on startup. If it fails to get a response, it loads the last known exchange rate from the config file located at ~/.config/exrate/exrate.json.
## Usage:
```
exrate <amount> <currency_from> [<currency_to> ...]
```
## Examples:
Output how much is 10 USD in all other currencies:
```
  exrate 10 USD
```
Output how much is 150 BYN in KZT:
```
  exrate 150 BYN KZT
```
Output how much is 4900 KZT in USD, RUB and BYN:
```
  exrate 4900 KZT USD RUB BYN
```
