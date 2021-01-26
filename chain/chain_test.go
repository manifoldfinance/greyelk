package chain

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeRPCAgent struct {
	fakeResult interface{}
	fakeErr    error
}

func (fa *fakeRPCAgent) CallMethod(result interface{}, method string, params []string) error {
	if fa.fakeErr == nil {
		val := reflect.ValueOf(result)
		if val.Kind() == reflect.Ptr {
			val.Elem().Set(reflect.ValueOf(fa.fakeResult))
		}
	}
	return fa.fakeErr
}

func (fa *fakeRPCAgent) Close() {
	return
}

func getTestChain(result interface{}, err error) *Blockchain {
	return &Blockchain{RPCAgent: &fakeRPCAgent{fakeResult: result, fakeErr: err}}
}

func TestGetCurrentBlock(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		input  interface{}
		result int
		err    error
	}{
		{
			"0x12",
			18,
			nil,
		},
		{
			"12",
			18,
			nil,
		},
		{
			"",
			-1,
			ErrGetCurrentBlockNum,
		},
		{
			"error",
			-1,
			ErrGetCurrentBlockNum,
		},
		{
			nil,
			-1,
			ErrGetCurrentBlockNum,
		},
	}
	for _, tt := range tests {
		bc := getTestChain(tt.input, tt.err)
		result, err := bc.GetCurrentBlock()
		assert.Equal(tt.result, result)
		assert.Equal(tt.err, err)
	}
}

func TestGetBlock(t *testing.T) {
	assert := assert.New(t)
	blockNum := 123
	testErr := fmt.Errorf("man made error")

	tests := []struct {
		input       map[string]interface{}
		err, retErr error
	}{
		{
			map[string]interface{}{"key": "value"},
			nil,
			nil,
		},
		{
			map[string]interface{}{"key": "value"},
			testErr,
			testErr,
		},
		{
			nil,
			nil,
			ErrGetBlockData,
		},
	}
	for _, tt := range tests {
		bc := getTestChain(tt.input, tt.err)
		_, err := bc.GetBlock(blockNum)
		assert.Equal(tt.retErr, err)
	}
}

func TestGenerateBlockEvent(t *testing.T) {
	assert := assert.New(t)
	blockNum := 123
	testErr := fmt.Errorf("man made error")

	tests := []struct {
		input       map[string]interface{}
		err         error
		emptyFields bool
	}{
		{
			map[string]interface{}{"hash": "value"},
			nil,
			false,
		},
		{
			map[string]interface{}{"key": "value"},
			testErr,
			true,
		},
	}
	for _, tt := range tests {
		bc := getTestChain(tt.input, tt.err)
		event := bc.GenerateBlockEvent(blockNum)
		assert.Equal(tt.emptyFields, event.Fields == nil)
	}
}
