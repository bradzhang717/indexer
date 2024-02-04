package jsonrpc

import (
	"errors"
	"github.com/uxuycom/indexer/xylog"
)

var rpcHandlersBeforeInitV2 = map[string]commandHandler{
	"inds_getInscriptionsStats":      indsGetInscriptionsStats,
	"inds_getTxs":                    indsGetTxs,
	"inds_getTicks":                  indsGetTicks, //handleFindAllInscriptions,
	"inds_getTransactionByAddress":   handleFindAddressTransactions,
	"inds_getBalanceByAddress":       indsGetBalanceByAddress,
	"inds_getHoldersByTick":          indsGetHoldersByTick,
	"inds_getLastBlockNumberIndexed": handleGetLastBlockNumber,
	"inds_getTickByCallData":         handleGetTxOperate,
	"inds_getTransactionByHash":      handleGetTxByHash,
	//"inscription.Tick":          handleFindInscriptionTick,
	//"address.Balance": handleFindAddressBalance,
}

func indsGetInscriptionsStats(s *RpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {

	req, ok := cmd.(*IndsGetInscriptionsStatsCmd)
	if !ok {
		return ErrRPCInvalidParams, errors.New("invalid params")
	}
	xylog.Logger.Infof("find all txs cmd params:%v", req)
	return findInscriptionsStats(s, req.Limit, req.Offset, req.SortMode)
}

func indsGetTxs(s *RpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {

	req, ok := cmd.(*IndsGetTxsCmd)
	if !ok {
		return ErrRPCInvalidParams, errors.New("invalid params")
	}
	xylog.Logger.Infof("find all txs cmd params:%v", req)
	return findTxs(s, req.Limit, req.Offset, req.SortMode)

}

func indsGetTicks(s *RpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	req, ok := cmd.(*IndsGetTicksCmd)
	if !ok {
		return ErrRPCInvalidParams, errors.New("invalid params")
	}

	xylog.Logger.Infof("find all Inscriptions cmd params:%v", req)
	return findInsciptions(s, req.Limit, req.Offset, req.Chain, req.Protocol, req.Tick, req.DeployBy, req.Sort, req.SortMode)
}

func indsGetBalanceByAddress(s *RpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	req, ok := cmd.(*IndsGetBalanceByAddressCmd)
	if !ok {
		return ErrRPCInvalidParams, errors.New("invalid params")
	}
	xylog.Logger.Infof("find user balances cmd params:%v", req)

	return findAddressBalances(s, req.Limit, req.Offset, req.Address, req.Chain, req.Protocol, req.Tick, req.Sort)
}

func indsGetHoldersByTick(s *RpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	req, ok := cmd.(*IndsGetHoldersByTickCmd)
	if !ok {
		return ErrRPCInvalidParams, errors.New("invalid params")
	}
	xylog.Logger.Infof("find user balances cmd params:%v", req)

	return findTickHolders(s, req.Limit, req.Offset, req.Chain, req.Protocol, req.Tick, req.SortMode)
}
