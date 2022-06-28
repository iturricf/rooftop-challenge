package challenge

import (
	"errors"
	"fmt"

	"github.com/iturricf/rooftop-challenge/client"
)

var apiClient client.Client

func Solve(login string) error {
	apiClient = &client.DefaultClient{}
	token, err := apiClient.GetToken(login)
	if err != nil {
		return err
	}

	blocks, err := apiClient.GetBlocks(token)
	if err != nil {
		return err
	}

	orderedBlocks, err := Check(blocks, token)
	if err != nil {
		return err
	}

	fmt.Println("Should be ordered. Checking...")
	checked, err := apiClient.VerifyBlocks(orderedBlocks, token)
	if err != nil {
		return err
	}

	if !checked {
		return errors.New("failed to find a solution")
	}

	fmt.Printf("Blocks are ordered\n\n")
	for i, block := range orderedBlocks {
		fmt.Printf("Block #%v: %v\n", i, block)
	}

	return nil
}

func Check(blocks []string, token string) ([]string, error) {
	sortedIndex := 0
	scanIndex := sortedIndex + 1
	// Assuming it will always be possible to find a solution given a list of blocks
	// Then, there is no need to check the last element as it will always be in order.
	for sortedIndex < len(blocks)-2 {
		pair := client.BlockPair{blocks[sortedIndex], blocks[scanIndex]}
		checked, err := apiClient.CheckPair(pair, token)
		if err != nil {
			return nil, err
		}
		if checked && scanIndex != sortedIndex+1 {
			swap := blocks[sortedIndex+1]
			blocks[sortedIndex+1] = blocks[scanIndex]
			blocks[scanIndex] = swap
		}
		if checked {
			sortedIndex += 1
			scanIndex = sortedIndex + 1
		} else {
			scanIndex += 1
		}

	}

	return blocks, nil
}
