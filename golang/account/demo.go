package account

import "errors"

// Virtual account for backtesting
type DemoAccount struct {
	// current total assets in demo dollars (if long is opened and price goes up, balance go up)
	// balance is equal to available (including realized profit) + unrealized profit
	balance uint64
	// part of balance that can be used for open new positions (not changing while holding opened pos)
	available uint64
	// {ticker: position}. Assuming only one position per ticker
	openedPositions map[string]*Position
}

func NewDemoAccount() *DemoAccount {
	return &DemoAccount{
		balance:         VIRTUAL_BALANCE,
		available:       VIRTUAL_BALANCE,
		openedPositions: make(map[string]*Position),
	}
}

func (a *DemoAccount) GetBalance() uint64 {
	return a.balance
}

func (a *DemoAccount) GetAvailableBalance() uint64 {
	return a.available
}

func (a *DemoAccount) GetOpenedPositionsCount() int {
	return len(a.openedPositions)
}

func (a *DemoAccount) IsOpened(ticker string) bool {
	p, exists := a.openedPositions[ticker]
	if exists && p.qcur > 0 {
		return true
	}
	return false
}

func (a *DemoAccount) IsLongOpened(ticker string) bool {
	p, exists := a.openedPositions[ticker]
	if exists && p.side == LongSide && p.qcur > 0 {
		return true
	}
	return false
}

func (a *DemoAccount) IsShortOpened(ticker string) bool {
	p, exists := a.openedPositions[ticker]
	if exists && p.side == ShortSide && p.qcur > 0 {
		return true
	}
	return false
}

func (a *DemoAccount) GetSide(ticker string) (PositionSide, error) {
	p, exists := a.openedPositions[ticker]
	if !exists {
		return "", errors.New("no opened position")
	}
	if p == nil {
		return "", errors.New("nil position")
	}
	return p.side, nil
}

func (a *DemoAccount) GetOpenPrice(ticker string) (uint64, error) {
	p, exists := a.openedPositions[ticker]
	if !exists {
		return uint64(0), errors.New("no opened position")
	}
	if p == nil {
		return uint64(0), errors.New("nil position")
	}
	return p.po, nil
}

func (a *DemoAccount) GetOpenQuantity(ticker string) (uint64, error) {
	p, exists := a.openedPositions[ticker]
	if !exists {
		return uint64(0), errors.New("no opened position")
	}
	if p == nil {
		return uint64(0), errors.New("nil position")
	}
	return p.qo, nil
}

func (a *DemoAccount) GetCurrentPrice(ticker string) (uint64, error) {
	p, exists := a.openedPositions[ticker]
	if !exists {
		return uint64(0), errors.New("no opened position")
	}
	if p == nil {
		return uint64(0), errors.New("nil position")
	}
	return p.pcur, nil
}

func (a *DemoAccount) GetCurrentQuantity(ticker string) (uint64, error) {
	p, exists := a.openedPositions[ticker]
	if !exists {
		return uint64(0), errors.New("no opened position")
	}
	if p == nil {
		return uint64(0), errors.New("nil position")
	}
	return p.qcur, nil
}

// price and quantity must be scaled
// returns scaled fee
func (a *DemoAccount) GetFee(price, quantity uint64) uint64 {
	return price * quantity / SCALE / FEE_SCALE
}

func (a *DemoAccount) createPosition(ticker string, side PositionSide, price, quantity uint64) error {
	fee := a.GetFee(price, quantity)
	notional := price * quantity / SCALE
	if notional+fee > a.available {
		return NotEnoughBalanceErr
	}
	p := &Position{
		ticker: ticker,
		side:   side,
		po:     price,
		qo:     quantity,
		pcur:   price,
		qcur:   quantity,
	}
	a.openedPositions[p.ticker] = p

	a.available -= notional
	a.available -= fee
	a.Update(ticker, price)

	return nil
}

// Opens a new position. Note: if long is opened, but got short, it will close long and open short for given quantity
// Example: had long +200, send short for -500 -> close long and open short for 500
// So if available balance is 1000 and long +200, we can't open short for 1100
func (a *DemoAccount) OpenPosition(ticker string, side PositionSide, price, quantity uint64) error {
	pos, exists := a.openedPositions[ticker]

	// swapping side
	if exists && pos.side != side {
		// closing old position
		err := a.CloseFullPosition(ticker, price)
		if err != nil {
			return err
		}
		// create new position
		return a.createPosition(ticker, side, price, quantity)
	}
	// if we adding more, calculating mean price and adding qty
	if exists && pos.side == side {
		fee := a.GetFee(price, quantity)
		notional := price * quantity / SCALE
		if notional+fee > a.available {
			return NotEnoughBalanceErr
		}

		num := pos.pcur*pos.qcur + price*quantity
		totalQty := pos.qcur + quantity

		pos.pcur = roundFraction(num, totalQty)
		pos.qcur = totalQty

		a.available -= notional
		a.available -= fee
		a.Update(ticker, price)

		return nil
	}
	// if no position
	return a.createPosition(ticker, side, price, quantity)
}

func (a *DemoAccount) CloseFullPosition(ticker string, price uint64) error {
	pos, exists := a.openedPositions[ticker]
	if !exists {
		return PositionNotOpenedErr
	}

	fee := a.GetFee(price, pos.qcur)

	if pos.side == LongSide {
		// available += baseCost + profit - fee = po*qc + (pc-po)*qc - fee = pc*qc - fee
		a.available += price*pos.qcur/SCALE - fee
	} else {
		// available += baseCost + profit - fee = po*qc + (po-pc)*qc - fee = (2*po-pc)*qc - fee
		a.available += (2*pos.po-price)*pos.qcur/SCALE - fee
	}

	pos.pcur = price
	pos.qcur = 0
	a.updateBalance()

	return nil
}

// If need to close all position prefer using CloseFullPosition()
func (a *DemoAccount) ClosePosition(ticker string, price, quantity uint64) error {
	pos, exists := a.openedPositions[ticker]
	if !exists {
		return PositionNotOpenedErr
	}
	if quantity > pos.qcur {
		return NotEnoughAssetsErr
	}

	fee := a.GetFee(price, quantity)

	if pos.side == LongSide {
		// available += baseCost + profit - fee = po*qc + (pc-po)*qc - fee = pc*qc - fee
		a.available += price*quantity/SCALE - fee
	} else {
		// available += baseCost + profit - fee = po*qc + (po-pc)*qc - fee = (2*po-pc)*qc - fee
		a.available += (2*pos.po-price)*quantity/SCALE - fee
	}

	pos.pcur = price
	pos.qcur -= quantity
	a.updateBalance()

	return nil
}

// calculates current balance
func (a *DemoAccount) updateBalance() {
	res := a.available
	for _, pos := range a.openedPositions {
		if pos.side == LongSide && pos.qcur > 0 {
			res += pos.pcur * pos.qcur / SCALE // po*qcur + (pcur-po)*qcur; base cost + current profit
		} else if pos.side == ShortSide && pos.qcur > 0 {
			res += (2*pos.po - pos.pcur) * pos.qcur / SCALE // po*qcur + (po-pcur)*qcur; base cost + current profit
		}
	}
	a.balance = res
}

// Updating position and balance depends on current price
// If no position is opened, return error
func (a *DemoAccount) Update(ticker string, price uint64) error {
	pos, exists := a.openedPositions[ticker]
	if !exists {
		return PositionNotOpenedErr
	}
	pos.pcur = price

	a.updateBalance()
	return nil
}
