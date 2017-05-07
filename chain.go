package main

import (
	"sync"
	"time"
)

type Blockchain struct {
	blocks []*Block
	mu     sync.RWMutex
}

func newBlockchain() *Blockchain {
	return &Blockchain{
		blocks: []*Block{genesisBlock},
		mu:     sync.RWMutex{},
	}
}

func (bc *Blockchain) len() int {
	return len(bc.blocks)
}

func (bc *Blockchain) getGenesisBlock() *Block {
	return bc.getBlock(0)
}

func (bc *Blockchain) getLatestBlock() *Block {
	return bc.getBlock(bc.len() - 1)
}

func (bc *Blockchain) getBlock(index int) *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.blocks[index]
}

func (bc *Blockchain) generateBlock(data string) (*Block, error) {
	prevBlock := bc.getLatestBlock()
	prevBlockHash, err := prevBlock.hash()
	if err != nil {
		return nil, err
	}

	return &Block{
		Index:        prevBlock.Index + 1,
		PreviousHash: prevBlockHash,
		Timestamp:    time.Now().Unix(),
		Data:         data,
	}, nil
}

func (bc *Blockchain) addBlock(block *Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.blocks = append(bc.blocks, block)

	return nil
}

func (bc *Blockchain) replaceBlocks(bcNew *Blockchain) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.blocks = bcNew.blocks

	return nil
}

func (bc *Blockchain) isValidGenesisBlock() (bool, error) {
	genesisBlockHash, err := genesisBlock.hash()
	if err != nil {
		return false, err
	}

	firstBlock := bc.getGenesisBlock()
	firstBlockHash, err := firstBlock.hash()
	if err != nil {
		return false, err
	}
	if firstBlockHash != genesisBlockHash {
		return false, nil
	}

	return true, nil
}

func (bc *Blockchain) isValid() (bool, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if bc.len() == 0 {
		return false, nil
	}

	ok, err := bc.isValidGenesisBlock()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	prevBlock := bc.getGenesisBlock()
	for i := 1; i < bc.len(); i++ {
		block := bc.getBlock(i)

		ok, err := isValidBlock(block, prevBlock)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}

		prevBlock = block
	}

	return true, nil
}
