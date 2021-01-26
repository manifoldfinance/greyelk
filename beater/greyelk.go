package beater

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/manifoldfinance/greyelk/chain"
	"github.com/manifoldfinance/greyelk/config"
)

// GreyELK configuration.
type GreyELK struct {
	done    chan struct{}
	config  config.Config
	client  beat.Client
	chain   *chain.Blockchain
	nextNum int
}

var (
	//ErrStartBlockTooLarge is returned when start_block is greater than the current chain height
	ErrStartBlockTooLarge = fmt.Errorf("Current Ethereum block number is less than start_block")

	evChannelSize = 3
)

// New creates an instance of greyelk.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	chain, err := chain.Init(c.EthRPCAddr)
	if err != nil {
		return nil, err
	}

	currBlockNum, err := chain.GetCurrentBlock()
	if err != nil {
		return nil, err
	}
	if currBlockNum < c.StartBlock {
		return nil, ErrStartBlockTooLarge
	}
	if c.StartBlock < 0 {
		c.StartBlock = currBlockNum
		logp.Info("Set the start_block to the current blockchain height %d", currBlockNum)
	}

	bt := &GreyELK{
		done:    make(chan struct{}),
		config:  c,
		chain:   chain,
		nextNum: c.StartBlock,
	}
	return bt, nil
}

// Run starts greyelk.
func (bt *GreyELK) Run(b *beat.Beat) error {
	logp.Info("greyelk is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(bt.config.Period)
	evBuffer := make(chan beat.Event, evChannelSize)

	//start the go routine which publishes events to Elasticsearch
	go bt.sendBlockEvent(evBuffer, bt.done)

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		bt.getBlocks(evBuffer, bt.done)
	}
}

// Stop stops greyelk.
func (bt *GreyELK) Stop() {
	bt.client.Close()
	bt.chain.RPCAgent.Close()
	close(bt.done)
}

func (bt *GreyELK) getBlocks(ec chan<- beat.Event, done <-chan struct{}) {
	currBlockNum, err := bt.chain.GetCurrentBlock()
	if err != nil {
		logp.Err("cannot get current block number %v", err)
		return
	}
	if bt.nextNum > currBlockNum {
		logp.Warn("current block height has been imported %d", currBlockNum)
		return
	}

	for i := bt.nextNum; i <= currBlockNum; i++ {
		select {
		case <-bt.done:
			logp.Info("Received exit signal. Stop fetching blocks")
			return
		default:
			event := bt.chain.GenerateBlockEvent(i)
			if event.Fields == nil {
				logp.Warn("Unable to get block %d", i)
				bt.nextNum = i
				return
			}
			ec <- event
		}
	}

	logp.Info("Created beat event(s) for block number %d - %d", bt.nextNum, currBlockNum)
	bt.nextNum = currBlockNum + 1
}

func (bt *GreyELK) sendBlockEvent(ec <-chan beat.Event, done <-chan struct{}) {
	var events []beat.Event

	for {
		select {
		case <-done:
			if len(events) > 0 {
				bt.publishPendingEvents(events)
			}
			return
		case event := <-ec:
			if len(events)%evChannelSize == 0 {
				if len(events) == evChannelSize {
					bt.publishPendingEvents(events)
				}
				events = make([]beat.Event, 0, evChannelSize)
			}
			events = append(events, event)
		}
	}
}

func (bt *GreyELK) publishPendingEvents(events []beat.Event) {
	size := len(events)
	if size == 0 {
		return
	}
	if _, ok := events[0].Fields["number"]; ok {
		blockNumHexStr := events[0].Fields["number"].(string)
		if strings.HasPrefix(blockNumHexStr, "0x") || strings.HasPrefix(blockNumHexStr, "0X") {
			blockNumHexStr = blockNumHexStr[2:]
		}
		if blockNum, err := strconv.ParseInt(blockNumHexStr, 16, 32); err == nil {
			logp.Info("Publishing block(s) from %d to %d", int(blockNum), int(blockNum)+size-1)
		}
	}

	bt.client.PublishAll(events)
}
