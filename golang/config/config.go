package config

import "errors"

type PositionSide string

const (
	SCALE           uint64       = 10000
	VIRTUAL_BALANCE uint64       = 100000 * SCALE // 100000$
	FEE_SCALE       uint64       = 1000           // 0.1% = 0.001 -> 1/0.001 = 1000
	LongSide        PositionSide = "long"
	ShortSide       PositionSide = "short"
)

var MarketClosedErr = errors.New("markets are closed")
var PositionNotOpenedErr = errors.New("position is not opened")
var NotEnoughBalanceErr = errors.New("not enough balance")
var NotEnoughAssetsErr = errors.New("not enough assets")
