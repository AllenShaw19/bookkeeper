package proto

const (
	PushTypeAppend = iota
	PushTypeCommit
	PushTypeCompare
	PushTypeTruncate
)

const (
	VoteResultUnknown = iota
	VoteResultAccept
	VoteResultRejectUnknownLeader
	VoteResultRejectUnexpectedLeader
	VoteResultRejectExpiredVoteTerm
	VoteResultRejectAlreadyVoted
	VoteResultRejectAlreadyHasLeader
	VoteResultRejectTermNotReady
	VoteResultRejectTermSmallThanLeader
	VoteResultRejectExpiredLeaderTerm
	VoteResultRejectSmallLogEndIndex
	VoteResultRejectTakingLeadership
)

const (
	VoteParseResultWaitToReVote = iota
	VoteParseResultReVoteImmediately
	VoteParseResultPassed
	VoteParseResultWaitToVoteNext
)
