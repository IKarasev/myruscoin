# MyRuscoin

Emulator of blockchain network with little perpetration capabilities. Build in GO with web interface.

## What is done

- **Blocks** - Simplified block object
    - contains necessary block information (height, merkle root hash, hash etc...)
    - Coinbase amount at the moment block is mined
    - Transactions list
- **Transactions** - simmple as it is
    - Input and output Utxo
    - Sign and Wallet publick key
- **Wallet**
    - it has address!
    - and list of it's utxos
    - can inititate transactions
    - and sign it with it's private key
- **Node** - blockchain node aka miiner
    - can mine blocks
    - can validate blocks
    - ... accept blocks
    - ... send blocks
- **Attacker**
    - steals the candidate block
    - changes block's headers and contents
    - injects block back to miner
    - or mines and sends it to other nodes

New block is finalized and mined on a tick triggered by the user.

# Running

0. Install [GO!](https://go.dev/doc/install)

1. Clone the repo

```
git clone https://github.com/IKarasev/myruscoin.git
```

2. Move to cloned repo directory

```bash
cd ./myruscoin
```

3. Download GO dependensies

```bash
go mod tidy
```

4. Run it with go

```
go run ./cmd/werbsrv/main.go
```

4. Or with make:

```bash
make run-web
```

5. Enjoy!

## Enviroment variables

You can set some settings with next enviroment variables (set in system or in .env file if you run via make)

| Variable | Default value | Description |
| -------------- | --------------- | ---- |
| COINBASE_START_AMOUNT | 1000000 | Coinbase amount on system start |
| MINE_DIFF | 40000 | Mining difficulty |
| RUSCOIN_HTTP_ADDR | 127.0.0.1 | ip address the web server will listen to |
| RUSCOIN_HTTP_PORT | 8080 | port the web server will listen to |
| RUSCOIN_RSS_UPDATE | 100 | Milliseconds, RSS queue update period |
| OP_PAUSE_MILISEC | 500 | Milliseconds, pause between node operations |
| WITH_LOG | true | show web server log or not |

# For development

To change or develop it you need to
- install GO dependesies
- install GO Templ
```bash
go install github.com/a-h/templ/cmd/templ@latest
```
- install [tailwindcss](https://tailwindcss.com/docs/installation)
- optionally install [Go Air](https://github.com/air-verse/air) to live rebuild-restart on code change
```bash
go install github.com/air-verse/air@latest
```

## Make commands:

- Run web server
```bash
make run-web
```
- Build templ and tailwindcss:
```bash
make web-gen
```
- Build templ, tailwindcss and run
```bash
make web-gen-run
```
- Run air (uses .air.toml as settings for air, rebuilds templ and tailwindcss before start)
```bash
make air
```




