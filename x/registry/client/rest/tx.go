package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"kyve/x/registry/types"
)

type CreatePoolRequest struct {
	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title          string       `json:"title" yaml:"title"`
	Description    string       `json:"description" yaml:"description"`
	Deposit        sdk.Coins    `json:"deposit" yaml:"deposit"`
	Name           string       `json:"name" yaml:"name"`
	Runtime        string       `json:"runtime" yaml:"runtime"`
	Logo           string       `json:"logo" yaml:"logo"`
	Versions       string       `json:"versions" yaml:"versions"`
	Config         string       `json:"config" yaml:"config"`
	StartHeight    uint64       `json:"startHeight" yaml:"startHeight"`
	UploadInterval uint64       `json:"uploadInterval" yaml:"uploadInterval"`
	OperatingCost  uint64       `json:"operatingCost" yaml:"operatingCost"`
}

func ProposalCreatePoolRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "create-pool",
		Handler:  newCreatePoolHandler(clientCtx),
	}
}

func newCreatePoolHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePoolRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewCreatePoolProposal(req.Title, req.Description, req.Name, req.Runtime, req.Logo, req.Versions, req.Config, req.StartHeight, req.UploadInterval, req.OperatingCost)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type UpdatePoolRequest struct {
	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title          string       `json:"title" yaml:"title"`
	Description    string       `json:"description" yaml:"description"`
	Deposit        sdk.Coins    `json:"deposit" yaml:"deposit"`
	Id             uint64       `json:"id" yaml:"id"`
	Name           string       `json:"name" yaml:"name"`
	Runtime        string       `json:"runtime" yaml:"runtime"`
	Logo           string       `json:"logo" yaml:"logo"`
	Versions       string       `json:"versions" yaml:"versions"`
	Config         string       `json:"config" yaml:"config"`
	UploadInterval uint64       `json:"uploadInterval" yaml:"uploadInterval"`
	OperatingCost  uint64       `json:"operatingCost" yaml:"operatingCost"`
}

func ProposalUpdatePoolRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "update-pool",
		Handler:  newUpdatePoolHandler(clientCtx),
	}
}

func newUpdatePoolHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdatePoolRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewUpdatePoolProposal(req.Title, req.Description, req.Id, req.Name, req.Runtime, req.Logo, req.Versions, req.Config, req.UploadInterval, req.OperatingCost)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type PausePoolRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`
	Id          uint64       `json:"id" yaml:"id"`
}

func ProposalPausePoolRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "pause-pool",
		Handler:  newPausePoolHandler(clientCtx),
	}
}

func newPausePoolHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PausePoolRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewPausePoolProposal(req.Title, req.Description, req.Id)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type UnpausePoolRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`
	Id          uint64       `json:"id" yaml:"id"`
}

func ProposalUnpausePoolRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "unpause-pool",
		Handler:  newUnpausePoolHandler(clientCtx),
	}
}

func newUnpausePoolHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UnpausePoolRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewUnpausePoolProposal(req.Title, req.Description, req.Id)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
