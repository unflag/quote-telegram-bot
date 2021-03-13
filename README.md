# quote-telegram-bot

quote-telegram-bot is a simple bot to query [Yahoo Finance](https://finance.yahoo.com/) for stock market data and is meant to provide brief view into quote indicators.
When communicated, bot answers with basic quote indicators and additionally allows querying of some charts. Also, Yahoo Finance API allows querying of currency exchange rates, so here it is.

## How to use
* Register your bot instance via [BotFather](https://t.me/botfather)
* Build quote-telegram-bot:
```shell
$ make
```
* Run bot with obtained token. Token can be provided via command-line
```shell
$ ./quote-telegram-bot -botToken ${YOUR_TOKEN_HERE}
```
or environment variable
```shell
$ BOT_TOKEN=${YOUR_TOKEN_HERE} ./quote-telegram-bot
```
* Docker container:
```shell
$ docker run -d --restart=always -e "BOT_TOKEN=${YOUR_TOKEN_HERE}" unflag/quote-telegram-bot:latest
```  
* To allow debug logging provide `-debug` flag or `DEBUG=true` env variable
* Enjoy communicating your bot!

## What to ask
* Search symbol names by its company name using commands like `/apple` or `/tesla`.
* Search financial info by stock or ETF symbols like `AAPL` or `VOO`. Yahoo Finance supports variety of 
  [markets](https://help.yahoo.com/kb/exchanges-data-providers-yahoo-finance-sln2310.html), 
  so you must add a specific suffix, if your stock is trading on some of these markets.
  For example - `YNDX.ME` for Yandex shares that are traded on the Moscow Exchange.
* Currency exchange rates can be queried using `RUB=X` or `CNY=X` syntax for exchange rates of USD/RUB and USD/CNY respectively,
  or `RUBUSD=X`/`CNYUSD=X` syntax for a specific pair.
* Letter case does not matter.