package async

// Note: These constants have specific values defined by the canonical ABI, so we use explicit values instead of iota.

const (
	EVENT_NONE         uint32 = 0
	EVENT_SUBTASK      uint32 = 1
	EVENT_STREAM_READ  uint32 = 2
	EVENT_STREAM_WRITE uint32 = 3
	EVENT_FUTURE_READ  uint32 = 4
	EVENT_FUTURE_WRITE uint32 = 5
)

const (
	STATUS_STARTING uint32 = 0
	STATUS_STARTED  uint32 = 1
	STATUS_RETURNED uint32 = 2
)

const (
	CALLBACK_CODE_EXIT  uint32 = 0
	CALLBACK_CODE_YIELD uint32 = 1
	CALLBACK_CODE_WAIT  uint32 = 2
)

const (
	RETURN_CODE_BLOCKED   uint32 = 0xFFFFFFFF
	RETURN_CODE_COMPLETED uint32 = 0
	RETURN_CODE_DROPPED   uint32 = 1
)
