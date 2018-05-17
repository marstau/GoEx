package bibox

import (
	"github.com/marstau/GoEx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var hb2 = NewBibox(http.DefaultClient, "", "")

func TestBibox_GetTicker(t *testing.T) {
	ticker, err := hb2.GetTicker(goex.BTS_CNY)
	assert.Nil(t, err)
	t.Log(ticker)
}

func TestBibox_GetDepth(t *testing.T) {
	depth, err := hb2.GetDepth(2, goex.BCC_CNY)
	assert.Nil(t, err)
	t.Log("asks: ", depth.AskList)
	t.Log("bids: ", depth.BidList)
}
