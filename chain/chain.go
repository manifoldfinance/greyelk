package chain

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/sambacha/gethulent/client"
)

// Blockchain contains pointer to a RPC client object
type Blockchain struct {
	RPCAgent client.Agent
}

var (
	// ErrGetBlockData cannot get block data from server
	ErrGetBlockData = fmt.Errorf("cannot get block data")
	// ErrGetCurrentBlockNum cannot get current block number from server
	ErrGetCurrentBlockNum = fmt.Errorf("cannnot get current block number from RPC response")
)

// Init initializes blockchain object
func Init(url string) (*Blockchain, error) {
	agent, err := client.New(url)
	if err != nil {
		return nil, err
	}
	return &Blockchain{RPCAgent: agent}, nil
}

//GetBlock gets the block information
func (bc *Blockchain) GetBlock(blockNum int) (common.MapStr, error) {
	var result interface{}
	params := []string{strconv.FormatInt(int64(blockNum), 10), "true"}
	err := bc.RPCAgent.CallMethod(&result, "eth_getBlockByNumber", params)
	if err != nil {
		logp.Err("GetBlock, failed to get block, number %d, error: %v", blockNum, err)
		return nil, err
	}
	if result == nil || reflect.ValueOf(result).IsNil() {
		logp.Warn("GetBlock, cannot get block data for block %d", blockNum)
		return nil, ErrGetBlockData
	}
	return result.(map[string]interface{}), nil
}

// GenerateBlockEvent returns beat event
func (bc *Blockchain) GenerateBlockEvent(blockNum int) beat.Event {
	block, err := bc.GetBlock(blockNum)
	if err != nil {
		logp.Err("GenerateBlockEvent, error: %v", err)
		return beat.Event{}
	}

	// workaround elastic mapping name conflict
	if h, ok := block["hash"]; ok {
		delete(block, "hash")
		block["blockhash"] = h
	}

	event := beat.Event{
		Timestamp: time.Now(),
		Fields:    block,
	}
	return event
}

//GetCurrentBlock gets the current block number
func (bc *Blockchain) GetCurrentBlock() (int, error) {
	var result interface{}
	err := bc.RPCAgent.CallMethod(&result, "eth_blockNumber", []string{})
	if err != nil || result == nil {
		logp.Err("GetCurrentBlock, cannot get current block number")
		return -1, ErrGetCurrentBlockNum
	}
	if result != nil && reflect.TypeOf(result).Kind() == reflect.String {
		numStr := result.(string)
		if strings.HasPrefix(numStr, "0x") {
			numStr = numStr[2:]
		}
		if num, err := strconv.ParseInt(numStr, 16, 32); err == nil {
			return int(num), nil
		}
	}
	return -1, ErrGetCurrentBlockNum
}
