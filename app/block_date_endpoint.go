package app

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type DateByBlockEndpoint struct {
	Context *Context
}

func NewDateByBlock(context *Context) *DateByBlockEndpoint {
	return &DateByBlockEndpoint{
		Context: context,
	}
}

func (it *DateByBlockEndpoint) GetDateByBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	query, err := parseGetDateByBlockQuery(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := it.Context.TableBlockDate.Get(query.Block, it.Context.DB)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	max, err := it.Context.TableBlockDate.Max(it.Context.DB)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	current, err := it.Context.Client.Eth.Syncing()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(GetDateByBlockResponse{
		Block:  query.Block,
		Date:   date.Date,
		Synced: max,
		Max:    current.HighestBlock,
	})
}

func parseGetDateByBlockQuery(r *http.Request) (*GetDateByBlock, error) {
	decoder := json.NewDecoder(r.Body)

	query := GetDateByBlock{}

	err := decoder.Decode(&query)

	if err != nil {
		fmt.Println("Error in parsing query: ", err)
		return nil, err
	}

	return &query, nil
}

type GetDateByBlock struct {
	Block int64 `json:"block"`
}

type GetDateByBlockResponse struct {
	Block  int64 `json:"block"`
	Date   int64 `json:"date"`
	Synced int64 `json:"synced"`
	Max    int64 `json:"max"`
}
