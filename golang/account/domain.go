package account

import "github.com/verdex4/wall-street-code/config"

type PositionSide string

const (
	SCALE                        = config.SCALE
	FEE_SCALE                    = config.FEE_SCALE
	VIRTUAL_BALANCE              = config.VIRTUAL_BALANCE
	LongSide        PositionSide = "long"
	ShortSide       PositionSide = "short"
)

var PositionNotOpenedErr = config.PositionNotOpenedErr
var NotEnoughBalanceErr = config.NotEnoughBalanceErr
var NotEnoughAssetsErr = config.NotEnoughAssetsErr

type Position struct {
	ticker string
	side   PositionSide // "long" or "short"
	po     uint64       // Open price
	qo     uint64       // Open quantity
	pcur   uint64       // Current price
	qcur   uint64       // Current quantity
}
