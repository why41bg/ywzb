package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"math"
	"math/big"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64  // Timestamp is the current time when the block is created
	Data          []byte // In Bitcoin specification, transactions (Data in there) are separate data structure
	PrevBlockHash []byte
	Hash          []byte // SHA256 as the default algorithm
	Nonce         int64  // A cryptographic term from Hashcash description
	Difficulty    int64
}

// ProofOfWork trying to find a nonce that makes the hash of the block lower than a target value, which called Hashcash
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

// Serialize the block to a byte array
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	_ = encoder.Encode(b)
	return result.Bytes()
}

type Chain struct {
	tip []byte
	Db  *bolt.DB
}

func (bc *Chain) AddBlock(data string) error {
	var lastHash []byte

	_ = bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	nb := newBlock(data, lastHash)

	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		err := b.Put(nb.Hash, nb.Serialize())
		if err != nil {
			return err
		}

		err = b.Put([]byte("l"), nb.Hash)
		if err != nil {
			return err
		}

		bc.tip = nb.Hash
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (bc *Chain) Iterator() *ChainIterator {
	return &ChainIterator{bc.tip, bc.Db}
}

type ChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bci *ChainIterator) Next() *Block {
	var block *Block
	if bytes.Equal(bci.currentHash, []byte{}) {
		return block
	}
	_ = bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		encodedBlock := b.Get(bci.currentHash)
		block = Deserialize(encodedBlock)
		return nil
	})
	bci.currentHash = block.PrevBlockHash
	return block
}

func newBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          nil,
		Difficulty:    int64(12),
	}
	block.calculateAndSetHash()
	return block
}

func newGenesisBlock() *Block {
	return newBlock("Genesis Block", nil)
}

func InitBlockChain() (*Chain, error) {
	var tip []byte
	db, err := bolt.Open("blocks", 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))

		if b == nil {
			genesis := newGenesisBlock()
			b, err := tx.CreateBucket([]byte("blocksBucket"))
			if err != nil {
				return err
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				return err
			}
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				return err
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	bc := Chain{tip, db}
	return &bc, nil
}

func Deserialize(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	_ = decoder.Decode(&block)
	return &block
}
