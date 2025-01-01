package ywzb

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64  // Timestamp is the current time when the block is created
	Data          []byte // In Bitcoin specification, transactions (Data in there) are separate data structure
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64 // A cryptographic term from Hashcash description
	Difficulty    int64
}

func (b *Block) ProofOfWork() (int64, []byte) {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-b.Difficulty))

	nonce := 0
	var hashInt big.Int
	var hash [32]byte

	for nonce < math.MaxInt64 {
		data := bytes.Join([][]byte{
			[]byte(strconv.FormatInt(b.Timestamp, 10)),
			b.Data,
			b.PrevBlockHash,
			[]byte(strconv.FormatInt(b.Difficulty, 10)),
			[]byte(strconv.FormatInt(int64(nonce), 10)),
		}, []byte{})
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return int64(nonce), hash[:]
}

func (b *Block) calculateAndSetHash() {
	nonce, h := b.ProofOfWork()
	b.Nonce = nonce
	b.Hash = h
}

type BlockChain struct {
	Blocks []*Block
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	b := newBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, b)
}

func newBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Difficulty:    int64(24),
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
