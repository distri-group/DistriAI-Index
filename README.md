# DistriAI-Backend-Solana

## System Requirements
- Linux-amd64

## Build
Requires Go1.20 or higher.
```
GOOS=linux GOARCH=amd64 go build
```
If all goes well, you will get a program called `distriai-backend-solana`.

## Run
### Step 1: Prepare configuration file
- Copy `/config` folder to where `distriai-backend-solana` program locate.
- Edit Database configuration in `./config/config.yml`.
```
Server:
  Mode: release
  Port: 8888

Database:
  Host:
  Port: 3306
  UserName:
  Password:
  Database: distri-ai-solana

Mailbox:
  Host:
  Port: 25
  Username:
  Password:

Chain:
  Rpc:
  ProgramId:
```

### Step 2: Start the distriai-backend-solana service
- New a screen window
```
screen -S distriai-backend-solana
```
- start service
```
./distriai-backend-solana
```
When the service is started, the `machines` and `orders` tables in database will be cleared, the latest data on the chain will be pulled.

## Stop
- Attach the screen window
```
screen -r distriai-backend-solana
```
- stop service

`CTRL +  C`
