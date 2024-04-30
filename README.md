# DistriAI-Index

## System Requirements
- Linux-amd64

## Build
Requires Go1.21 or higher.
- Windows
```
set GOOS=linux
set GOARCH=amd64
go build
```
- Linux or Mac
```
GOOS=linux GOARCH=amd64 go build
```
If all goes well, you will get a program called `distriai-index-solana`.

## Run
### Step 1: Prepare configuration file
- Copy `/config` folder to where `distriai-index-solana` program locate.
- Edit configuration in `./config/config.yml`.

https://github.com/distri-group/DistriAI-Index/blob/dbe7077a310e6a9a6c0fec4f41e7199004f26df5/config/config.yml#L1-L35

### Step 2: Start the distriai-index-solana service
- New a screen window
```
screen -S distriai-index-solana
```
- start service
```
./distriai-index-solana
```
When the service is started, the `machines` and `orders` tables in database will be cleared, the latest data on the chain will be pulled.
- Detach the screen window

`CTRL +  A` + `D`

## Stop
- Attach the screen window
```
screen -r distriai-index-solana
```
- stop service

`CTRL +  C`

- Exit the screen window
```
exit
```
