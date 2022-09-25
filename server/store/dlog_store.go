package store

import (
	pb "github.com/AllenShaw19/bookkeeper/proto"
	"github.com/AllenShaw19/bookkeeper/server"
)

type DLogStore interface {
	GetMemberState() *server.MemberState
	AppendAsLeader(entry *pb.Entry) error
	AppendAsFollower(entry *pb.Entry, leaderTerm int64, leaderId string) error
	Get(index int64) *pb.Entry
	UpdateCommittedIndex(term, committedIndex int64) error
	GetCommittedIndex() int64
	GetLogEndTerm() int64
	GetLogEndIndex() int64
	GetLogBeginIndex() int64
	UpdateLogEndIndexAndTerm() error
	Flush()
	Truncate(entry *pb.Entry, leaderTerm int64, leaderId string) int64
	Startup()
	Shutdown()
}
