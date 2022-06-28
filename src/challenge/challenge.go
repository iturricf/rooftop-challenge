package challenge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Token struct {
	Token string `json:"token"`
}

type Blocks struct {
	Data      []string `json:"data"`
	ChunkSize int      `json:"chunkSize"`
	Length    int      `json:"length"`
}

type BlockPair [2]string

type BlockPairRequest struct {
	Blocks BlockPair `json:"blocks"`
}

type EncodedBlock struct {
	Encoded string `json:"encoded"`
}

type CheckResponse struct {
	Message bool `json:"message"`
}

func SolveWithLogin(login string) error {
	token, err := GetTokenWithLogin(login)
	if err != nil {
		return err
	}

	blocks, err := FetchBlocks(*token)
	if err != nil {
		return err
	}

	orderedBlocks, err := Check(blocks, *token)
	if err != nil {
		return err
	}

	for i, block := range orderedBlocks {
		fmt.Printf("Block #%v: %v\n", i, block)
	}

	return nil
}

func GetTokenWithLogin(login string) (*Token, error) {
	resp, err := http.Get("https://rooftop-career-switch.herokuapp.com/token?email=" + login)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response error [err=%v]", resp.Status)
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	token := Token{}
	err = json.Unmarshal(response, &token)
	if err != nil {
		return nil, nil
	}

	return &token, nil
}

func FetchBlocks(token Token) ([]string, error) {
	resp, err := http.Get("https://rooftop-career-switch.herokuapp.com/blocks?token=" + token.Token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response error [err=%v]", resp.Status)
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	blocks := Blocks{}
	err = json.Unmarshal(response, &blocks)
	if err != nil {
		return nil, err
	}

	return blocks.Data, nil
}

func Check(blocks []string, token Token) ([]string, error) {
	sortedIndex := 0
	scanIndex := sortedIndex + 1
	for sortedIndex < len(blocks)-1 {
		pair := BlockPair{blocks[sortedIndex], blocks[scanIndex]}
		checked, err := checkPair(pair, token)
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

	fmt.Println("Should be ordered. Checking ...")
	checked, err := checkAll(blocks, token)
	if err != nil {
		return nil, err
	}

	if checked {
		fmt.Println("Checked")
	}

	fmt.Println("NOT Checked")
	return blocks, nil
}

func checkPair(blocks BlockPair, token Token) (bool, error) {
	blockData := BlockPairRequest{Blocks: blocks}
	data, err := json.Marshal(blockData)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(
		"https://rooftop-career-switch.herokuapp.com/check?token="+token.Token,
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("response error [err=%v]", resp.Status)
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	checkResponse := CheckResponse{}
	err = json.Unmarshal(response, &checkResponse)
	if err != nil {
		return false, err
	}

	return checkResponse.Message, nil
}

func checkAll(blocks []string, token Token) (bool, error) {
	encodedBlocks := EncodedBlock{Encoded: strings.Join(blocks[:], "")}
	jsonBlocks, err := json.Marshal(encodedBlocks)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(
		"https://rooftop-career-switch.herokuapp.com/check?token="+token.Token,
		"application/json",
		bytes.NewBuffer(jsonBlocks),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("response error [err=%v]", resp.Status)
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	checkResponse := CheckResponse{}
	err = json.Unmarshal(response, &checkResponse)
	if err != nil {
		return false, err
	}

	return checkResponse.Message, nil
}
