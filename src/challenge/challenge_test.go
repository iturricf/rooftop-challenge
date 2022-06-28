package challenge

import (
	"reflect"
	"testing"

	"github.com/iturricf/rooftop-challenge/client"
	"golang.org/x/exp/slices"
)

type clientMock struct{}

var checkPairMock func(blocks client.BlockPair, token string) (bool, error)

func (c clientMock) GetToken(login string) (string, error) {
	return "token", nil
}

func (c clientMock) GetBlocks(token string) ([]string, error) {
	return []string{"a", "b", "c"}, nil
}

func (c clientMock) CheckPair(blocks client.BlockPair, token string) (bool, error) {
	return checkPairMock(blocks, token)
}

func (c clientMock) VerifyBlocks(blocks []string, token string) (bool, error) {
	return true, nil
}

func TestCheck(t *testing.T) {
	expected := []string{
		"f319",
		"46ec",
		"c1c7",
		"3720",
		"c7df",
		"c4ea",
		"4e3e",
		"80fd",
	}

	apiClient = &clientMock{}
	checkPairMock = func(blocks client.BlockPair, token string) (bool, error) {
		ordered := []string{
			"f319",
			"46ec",
			"c1c7",
			"3720",
			"c7df",
			"c4ea",
			"4e3e",
			"80fd",
		}

		idx := slices.IndexFunc(ordered, func(el string) bool { return el == blocks[0] })

		if idx == len(ordered)-1 || idx == -1 {
			return false, nil
		}

		if ordered[idx+1] == blocks[1] {
			return true, nil
		}

		return false, nil
	}

	blocks := []string{
		"f319",
		"3720",
		"4e3e",
		"46ec",
		"c7df",
		"c1c7",
		"80fd",
		"c4ea",
	}

	result, err := Check(blocks, "token")
	if err != nil {
		t.Errorf("failed while checking blocks")
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, expected %v", result, expected)
	}
}
