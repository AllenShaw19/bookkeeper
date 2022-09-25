package server

import (
	pb "github.com/AllenShaw19/bookkeeper/proto"
	"github.com/AllenShaw19/bookkeeper/server/statemachine"
	"github.com/AllenShaw19/bookkeeper/server/store"
)

type DLogServer interface {
	HandleAppend(req *pb.AppendEntryReq) (pb.AppendEntryResp, error)
	HandleGet(req *pb.GetEntriesReq) (pb.GetEntriesResp, error)
	HandleMetadata(req *pb.MetadataReq) (pb.MetadataResp, error)
	HandleLeadershipTransfer(req *pb.LeadershipTransferReq) (*pb.LeadershipTransferResp, error)
	HandleVote(req *pb.VoteReq) (*pb.VoteResp, error)
	HandleHeartBeat(req *pb.HeartBeatReq) (*pb.HeartBeatResp, error)
	HandlePull(req *pb.PullEntriesReq) (*pb.PullEntriesResp, error)
	HandlePush(req *pb.PushEntriesReq) (*pb.PushEntriesResp, error)
}

func NewDLogServer(conf *DLogConfig) DLogServer {
	return newDLogServer(conf)
}

type dLogServer struct {
	memberState       *MemberState
	dLogConfig        *DLogConfig
	dLogStore         store.DLogStore
	dLogRpcService    DLogRpcService
	dLogEntryPusher   *DLogEntryPusher
	dLogLeaderElector *DLogLeaderElector
	fmsCaller         *statemachine.StateMachineCaller

	isStarted bool // atomic
}

func newDLogServer(dLogConfig *DLogConfig, serverConfig *RpcServerConfig,
	clientConfig *RpcClientConfig) *dLogServer {
	dLogConfig.Init()
	s := &dLogServer{}
	s.dLogConfig = dLogConfig
	s.memberState = NewMemberState(dLogConfig)
	s.dLogStore = createDLogStore(dLogConfig.StoreType, s.dLogConfig, s.memberState)
	s.dLogRpcService = NewLogRpcService(s, serverConfig, clientConfig)
	s.dLogEntryPusher = NewDLogEntryPusher(dLogConfig, s.memberState, s.dLogStore, s.dLogRpcService)
	s.dLogLeaderElector = NewDLogLeaderElector(dLogConfig, s.memberState, s.dLogRpcService)
}

func (s *dLogServer) Startup() {

}

func (s *dLogServer) Shutdown() {

}

func (s *dLogServer) HandleHeartBeat()

//
func createDLogStore(storeType string, config *DLogConfig, state *MemberState) {
	if storeType == StoreTypeMemory {
		return NewDLogMemoryStore(config, state)
	} else {
		return NewDLogMmapFileStore(config, state)
	}
}
