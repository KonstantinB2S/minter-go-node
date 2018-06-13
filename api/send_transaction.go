package api

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/MinterTeam/minter-go-node/core/code"
	"net/http"
	"strings"
	//"github.com/tendermint/tendermint/rpc/core/types"
	"github.com/MinterTeam/minter-go-node/core/types"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/common"
)

type SendTransactionRequest struct {
	Transaction string `json:"transaction"`
}

type ResultBroadcastTxCommit struct {
	CheckTx   abci.ResponseCheckTx   `json:"check_tx"`
	DeliverTx abci.ResponseDeliverTx `json:"deliver_tx"`
	Hash      common.HexBytes        `json:"hash"`
	Height    int64                  `json:"height"`
}

func SendTransaction(w http.ResponseWriter, r *http.Request) {

	var req SendTransactionRequest
	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	json.Unmarshal(body, &req)

	result := new(ResultBroadcastTxCommit)
	_, err := client.Call("broadcast_tx_commit", map[string]interface{}{
		"tx": types.Hex2Bytes(strings.TrimLeft(req.Transaction, "Mx")),
	}, result)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if result.CheckTx.Code != code.OK {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(Response{
			Code: result.CheckTx.Code,
			Log:  "Check tx error: " + result.CheckTx.Log,
		})
		return
	}

	if result.DeliverTx.Code != code.OK {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(Response{
			Code: result.CheckTx.Code,
			Log:  "Deliver tx error: " + result.DeliverTx.Log,
		})
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Response{
		Code:   code.OK,
		Result: "Mx" + strings.ToLower(result.Hash.String()),
	})
}
