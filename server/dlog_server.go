package server

import "github.com/AllenShaw19/bookkeeper/server/store"

type DLogServer interface {
}

type dLogServer struct {
	memberState    *MemberState
	dLogConfig     *DLogConfig
	dLogStore      store.DLogStore
	dLogRpcService DLogRpcService
}
