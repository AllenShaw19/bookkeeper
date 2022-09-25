package code

const (
	Unknown int32 = -1
	//Metadata                 = 50000
	//Append                   = 50001
	//Get                      = 50002
	//Vote                     = 51001
	//HeartBeat                = 51002
	//Pull                     = 51003
	//Push                     = 51004
	//LeadershipTransfer       = 51005

	Success                 = 200
	Timeout                 = 300
	MetadataError           = 301
	NetworkError            = 302
	Unsupported             = 303
	UnknownGroup            = 400
	UnknownMember           = 401
	UnexpectedMember        = 402
	ExpiredTerm             = 403
	NotLeader               = 404
	NotFollower             = 405
	InconsistentState       = 406
	InconsistentTerm        = 407
	InconsistentIndex       = 408
	InconsistentLeader      = 409
	IndexOutOfRange         = 410
	UnexpectedArgument      = 411
	RepeatedRequest         = 412
	RepeatedPush            = 413
	DiskError               = 414
	DiskFull                = 415
	TermNotReady            = 416
	FallBehindTooMuch       = 417
	TakeLeadershipFailed    = 418
	IndexLessThanLocalBegin = 419
	RequestWithEmptyBody    = 420
	InternalError           = 500
	TermChanged             = 501
	WaitQuorumAckTimeout    = 502
	LeaderPendingFull       = 503
	IllegalMemberState      = 504
	LeaderNotReady          = 505
	LeaderTransferring      = 506
)
