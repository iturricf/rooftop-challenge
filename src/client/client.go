package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	getTokenURL    = "https://rooftop-career-switch.herokuapp.com/token?email="
	getBlocksURL   = "https://rooftop-career-switch.herokuapp.com/blocks?token="
	checkBlocksURL = "https://rooftop-career-switch.herokuapp.com/check?token="
)

type Client interface {
	GetToken(login string) (string, error)
	GetBlocks(token string) ([]string, error)
	CheckPair(blocks BlockPair, token string) (bool, error)
	VerifyBlocks(blocks []string, token string) (bool, error)
}

type DefaultClient struct{}

type TokenResponse struct {
	Token string `json:"token"`
}

type BlocksResponse struct {
	Data      []string `json:"data"`
	ChunkSize int      `json:"chunkSize"`
	Length    int      `json:"length"`
}

type BlockPair [2]string

type CheckPayload struct {
	Blocks BlockPair `json:"blocks"`
}

type VerifyPayload struct {
	Encoded string `json:"encoded"`
}

type CheckResponse struct {
	Message bool `json:"message"`
}

func (c *DefaultClient) GetToken(login string) (string, error) {
	resp, err := http.Get(getTokenURL + login)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("response error [err=%v]", resp.Status)
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	token := TokenResponse{}
	err = json.Unmarshal(response, &token)
	if err != nil {
		return "", nil
	}

	return token.Token, nil
}

func (c *DefaultClient) GetBlocks(token string) ([]string, error) {
	resp, err := http.Get(getBlocksURL + token)
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

	blocks := BlocksResponse{}
	err = json.Unmarshal(response, &blocks)
	if err != nil {
		return nil, err
	}

	return blocks.Data, nil
}

func (c *DefaultClient) CheckPair(blocks BlockPair, token string) (bool, error) {
	blockData := CheckPayload{Blocks: blocks}
	data, err := json.Marshal(blockData)
	if err != nil {
		return false, err
	}

	response, err := postBlocksResponse(token, data)
	if err != nil {
		return false, err
	}

	return response.Message, nil
}

func (c *DefaultClient) VerifyBlocks(blocks []string, token string) (bool, error) {
	payload := VerifyPayload{Encoded: strings.Join(blocks[:], "")}
	data, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	response, err := postBlocksResponse(token, data)
	if err != nil {
		return false, err
	}

	return response.Message, nil
}

func postBlocksResponse(token string, data []byte) (*CheckResponse, error) {
	resp, err := http.Post(
		checkBlocksURL+token,
		"application/json",
		bytes.NewBuffer(data),
	)
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

	checkResponse := CheckResponse{}
	err = json.Unmarshal(response, &checkResponse)
	if err != nil {
		return nil, err
	}

	return &checkResponse, nil
}
