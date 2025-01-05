package block

import (
	"bytes"
	"github.com/boltdb/bolt"
	"testing"
)

func TestInitBlockChain(t *testing.T) {
	_, err := InitBlockChain()
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
}

func TestBlockChain_AddBlock(t *testing.T) {
	bc, err := InitBlockChain()
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	err = bc.AddBlock("TestBlockChain_AddBlock")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	err = bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocksBucket"))
		bBlock := b.Get(bc.tip)
		block := Deserialize(bBlock)
		if bytes.Equal(block.Hash, bc.tip) == false {
			t.Errorf("Chain tip is not equal to the last block hash")
		}
		if bytes.Equal(block.Data, []byte("TestBlockChain_AddBlock")) == false {
			t.Errorf("Chain data is not equal to the last block data")
		}
		return nil
	})
}

func TestBlockChainIterator_Next(t *testing.T) {
	bc, err := InitBlockChain()
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	bci := bc.Iterator()
	h := bci.currentHash
	for block := bci.Next(); block != nil; {
		if bytes.Equal(block.Hash, h) == false {
			t.Errorf("Chain tip is not equal to the last block hash. %x != %x", block.Hash, h)
		}
		h = bci.currentHash
		block = bci.Next()
	}
}
