package ruscoin

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
)

const (
	GENESIS_BLOCK_PREV byte = 0b00
)

var (
	COINBASE_START_AMOUNT int    = 1000000
	REWARD_AMOUNT         int    = 5
	COINBASE_ADDR         string = "coinbase"
	MINE_DIFF             string = "20"
	MINE_BASE             string = "115792089237316195423570985008687907853269984665640564039457584007913129639936"
	NONCE_MAX             int    = 2147483647
)

func InitRuscoinSettings() error {
	errStr := ""
	if v := os.Getenv("COINBASE_START_AMOUNT"); v != "" {
		if c, err := strconv.Atoi(v); err == nil {
			COINBASE_START_AMOUNT = c
		} else {
			errStr += "Failed to parse COINBASE_START_AMOUNT env variable\n"
		}
	}

	if v := os.Getenv("REWARD_AMOUNT"); v != "" {
		if c, err := strconv.Atoi(v); err == nil {
			REWARD_AMOUNT = c
		} else {
			errStr += "Failed to parse REWARD_AMOUNT env variable\n"
		}
	}

	if v := os.Getenv("MINE_DIFF"); v != "" {
		if _, ok := new(big.Int).SetString(v, 10); ok {
			MINE_DIFF = v
		} else {
			errStr += "Failed to parse MINE_DIFF env variable\n"
		}
	}

	if v := os.Getenv("NONCE_MAX"); v != "" {
		if c, err := strconv.Atoi(v); err == nil {
			NONCE_MAX = c
		} else {
			errStr += "Failed to parse NONCE_MAX env variable\n"
		}
	}

	if errStr != "" {
		return fmt.Errorf(errStr)
	}
	return nil
}
