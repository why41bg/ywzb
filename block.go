package ywzb

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64  // Timestamp is the current time when the block is created
	Data          []byte // In Bitcoin specification, transactions (Data in there) are separate data structure
	PrevBlockHash []byte
	Hash          []byte
}

func (b *Block) calculateAndSetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

type BlockChain struct {
	Blocks []*Block
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := newBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func newBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
	}
	block.calculateAndSetHash()
	return block
}

func newGenesisBlock() *Block {
	return newBlock("Genesis Block", []byte{})
}

func NewBlockChain() *BlockChain {
	return &BlockChain{Blocks: []*Block{newGenesisBlock()}}
}
