package challenge

import (
	"fmt"
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
	return []string{}, nil
}

func (c clientMock) CheckPair(blocks client.BlockPair, token string) (bool, error) {
	return checkPairMock(blocks, token)
}

func (c clientMock) VerifyBlocks(blocks []string, token string) (bool, error) {
	return true, nil
}

func TestCheck(t *testing.T) {
	var testMap = []struct {
		blocks   []string
		expected []string
	}{
		{
			[]string{
				"f319",
				"3720",
				"4e3e",
				"46ec",
				"c7df",
				"c1c7",
				"80fd",
				"c4ea",
			},
			[]string{
				"f319",
				"46ec",
				"c1c7",
				"3720",
				"c7df",
				"c4ea",
				"4e3e",
				"80fd",
			},
		},
		{
			[]string{
				"a",
				"d",
				"c",
				"e",
				"g",
				"f",
				"b",
			},
			[]string{
				"a",
				"b",
				"c",
				"d",
				"e",
				"f",
				"g",
			},
		},
	}

	apiClient = &clientMock{}

	for i, test := range testMap {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			checkPairMock = func(blocks client.BlockPair, token string) (bool, error) {
				idx := slices.IndexFunc(test.expected, func(el string) bool { return el == blocks[0] })

				if idx == len(test.expected)-1 || idx == -1 {
					return false, nil
				}

				if test.expected[idx+1] == blocks[1] {
					return true, nil
				}

				return false, nil
			}

			result, err := Check(test.blocks, "token")
			if err != nil {
				t.Errorf("failed while checking blocks")
			}

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("got %v, expected %v", result, test.expected)
			}
		})
	}
}
