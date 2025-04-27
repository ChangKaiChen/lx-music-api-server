package consts

import "time"

const (
	SuccessCode     = 200
	ForbiddenCode   = 403
	NotFoundCode    = 404
	ServerErrorCode = 500
	TooManyCode     = 429
)

const (
	SuccessStatusCode     = 200
	ForbiddenStatusCode   = 403
	TooManyStatusCode     = 429
	ServerErrorStatusCode = 500
)

const (
	SentinelThreshold        = 100
	SentinelStatIntervalInMs = 1000
)
const (
	SingleKeyRequestRate   = 10
	SingleKeyRequestWindow = 5 * time.Second
)
const (
	MuxConnectionNum = 2
)
