package challenge

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Token struct {
	Token string `json:"token"`
}

type Blocks struct {
	Data      []string `json:"data"`
	ChunkSize int      `json:"chunkSize"`
	Length    int      `json:"length"`
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

	for i, block := range blocks {
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
