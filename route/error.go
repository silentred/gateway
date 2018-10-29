package route

import "github.com/silentred/gateway/util"

const (
	DupRoute  = 100001
	NotFound  = 100404
	NoTargets = 100510
)

var (
	ErrDupRoute = util.NewError(DupRoute, "duplicated route")
	ErrNotFound = util.NewError(NotFound, "route not found")
)
