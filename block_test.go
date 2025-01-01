package ywzb

import (
	"bytes"
	"testing"
)

func TestNewBlockChain(t *testing.T) {
	bc := NewBlockChain()
	if len(bc.Blocks) != 1 {
		t.Errorf("Expected 1 block in the chain")
	}
	if string(bc.Blocks[0].Data) != "Genesis Block" {
		t.Errorf("Expected data in the block to be 'Genesis Block'")
	}
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockChain()
	bc.AddBlock("First Block")
	if len(bc.Blocks) != 2 {
		t.Errorf("Expected 2 blocks in the chain")
	}
	if string(bc.Blocks[1].Data) != "First Block" {
		t.Errorf("Expected data in the block to be 'First Block'")
	}
	if !bytes.Equal(bc.Blocks[1].PrevBlockHash, bc.Blocks[0].Hash) {
		t.Errorf("Expected PrevBlockHash to be the hash of the previous block")
	}
}
