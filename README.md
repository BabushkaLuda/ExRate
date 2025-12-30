This CLI uses <a href="https://www.exchangerate-api.com">This API</a> to show you exchange rate of chosen currency.
It makes API call on start up? but if fails to get responce, loads last known exchange rate from config file (~/.config/exrate/exrate.json).
Usage:  exrate <amount> <currency_from> [<currency_to> ...]
Examples:
  `exrate 10 USD`
  `exrate 150 BYN KZT`
  `exrate 4900 KZT USD RUB BYN`
