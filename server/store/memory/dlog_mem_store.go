package memory

import (
	pb "github.com/AllenShaw19/bookkeeper/proto"
	"github.com/AllenShaw19/bookkeeper/server"
)

type DLogMemoryStore struct {
	beginIndex     int64
	endIndex       int64
	committedIndex int64
	endTerm        int64
	cachedEntries  map[int64]*pb.Entry
	dLogConfig     *server.DLogConfig
	memberState    *server.MemberState
}

func NewDLogMemoryStore(config *server.DLogConfig, state *server.MemberState) *DLogMemoryStore {
	s := &DLogMemoryStore{}
	s.dLogConfig = config
	s.memberState = state
	return s
}

func (s *DLogMemoryStore) AppendAsLeader(entry *pb.Entry) (*pb.Entry, error) {
	// TOOD: check member is leaer
	if !s.memberState.IsLeader() {
		return
	}
	// lock memberstate
	s.memberState.Lock()
	defer s.memberState.Unlock()

	// check leader again

	s.endIndex++
	s.committedIndex++
	s.endTerm = s.memberState.CurrTerm()

	entry.Index = s.endIndex
	entry.Term = s.memberState.CurrTerm()

	s.cachedEntries[entry.Index] = entry
	if s.beginIndex == -1 {
		s.beginIndex = s.endIndex
	}
	s.UpdateLogEndIndexAndTerm()
	return entry, nil
}

func (s *DLogMemoryStore) AppendAsFollower(entry *pb.Entry, leaderTerm int64, leaderId string) (*pb.Entry, error) {
	if !s.memberState.IsFollower() {
		return
	}
	// lock memberstate
	s.memberState.Lock()
	defer s.memberState.Unlock()

	// check is follower again
	// check leader term
	// check leader id
	s.endTerm = entry.Term
	s.endIndex = entry.Index
	s.committedIndex = entry.Index
	s.cachedEntries[entry.Index] = entry
	if s.beginIndex == -1 {
		s.beginIndex = s.endIndex
	}

	s.UpdateLogEndIndexAndTerm()
	return entry, nil
}

func (s *DLogMemoryStore) GetMemberState() *server.MemberState {
	return s.memberState
}

func (s *DLogMemoryStore) UpdateLogEndIndexAndTerm() error {
	if s.GetMemberState() != nil {
		return s.GetMemberState().UpdateLogEndIndexAndTerm(s.GetLogEndIndex(), s.GetLogEndTerm())
	}
	return nil
}

func (s *DLogMemoryStore) GetLogEndTerm() int64 {
	return s.endTerm
}
func (s *DLogMemoryStore) GetLogEndIndex() int64 {
	return s.endIndex
}
func (s *DLogMemoryStore) GetLogBeginIndex() int64 {
	return s.beginIndex
}
func (s *DLogMemoryStore) GetCommittedIndex() int64 {
	return s.committedIndex
}
func (s *DLogMemoryStore) Get(index int64) *pb.Entry {
	return s.cachedEntries[index]
}

func (s *DLogMemoryStore) Flush() {}

func (s *DLogMemoryStore) Truncate(entry *pb.Entry, leaderTerm int64, leaderId string) int64 {
	entry, err := s.AppendAsFollower(entry, leaderTerm, leaderId)
	if err != nil {

	}
	return entry.Index
}

func (s *DLogMemoryStore) Startup() {}

func (s *DLogMemoryStore) Shutdown() {}

func (s *DLogMemoryStore) UpdateCommittedIndex(term, committedIndex int64) error { return nil }
