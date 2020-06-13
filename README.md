# Mussum
> A telegram bot to help you fight spam written in Go

This bot is inspired by a Brazilian comedian called [Mussum](https://en.wikipedia.org/wiki/Mussum). He was very popular between the 70's and 90's and he had a particular way of speaking portuguese while acting and I tried to bring that style into the way the bot talks to the users. Hope you enjoy it.

<p align="center">
  <img width="300" height="400" src="mussum.jpg">
</p>


## Build
The following command should build Mussum as well as its docker image.
```
make build
```

## Running locally
Use the [BotFather](https://core.telegram.org/bots) on Telegram to generate an APIKEY.
Start the bot locally by running:
```
APIKEY=your_api_key make run
```
Mussum speaks portuguese by default.

## Running tests
```
make test
```

## Changing the language
Mussum currently supports portuguese(pt) and english(en).
```
APIKEY=your_api_key LANGUAGE=en make run
```
You can easily add more languages by implementing the language interface. Please see the `language` package.

