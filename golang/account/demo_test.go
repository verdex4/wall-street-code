package account

import (
	"errors"
	"testing"
)

func TestInitialization(t *testing.T) {
	acc := NewDemoAccount()
	bal := acc.GetBalance()
	available := acc.GetAvailableBalance()
	l := acc.GetOpenedPositionsCount()

	if bal != VIRTUAL_BALANCE {
		t.Errorf("balance must be %d, got %d", VIRTUAL_BALANCE, bal)
	}
	if available != VIRTUAL_BALANCE {
		t.Errorf("available balance must be %d, got %d", VIRTUAL_BALANCE, available)
	}
	if l != 0 {
		t.Errorf("expected 0 opened positions, got %d", l)
	}
}

// test for creating new position, including correct work with balance
func TestCreatePositionLong(t *testing.T) {
	acc := NewDemoAccount()

	ticker := "NVDA"
	price := uint64(200 * SCALE)
	quantity := uint64(3 * SCALE)

	err := acc.OpenPosition(ticker, LongSide, price, quantity)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	l := acc.GetOpenedPositionsCount()
	if l != 1 {
		t.Errorf("expected one position, got %d", l)
	}

	b := acc.IsLongOpened(ticker)
	if !b {
		t.Errorf("expected opened long position = true, got %v", b)
	}

	po, err := acc.GetOpenPrice(ticker)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	pc, err := acc.GetCurrentPrice(ticker)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if po != price {
		t.Errorf("expected opened price = %d, got %d", price, po)
	}
	if po != pc {
		t.Errorf("expected opened price = current price, opened = %d, current = %d", po, pc)
	}

	qo, err := acc.GetOpenQuantity(ticker)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	qc, err := acc.GetCurrentQuantity(ticker)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if qo != quantity {
		t.Errorf("expected opened quantity = %d, got %d", quantity, qo)
	}
	if qo != qc {
		t.Errorf("expected opened quantity = current quantity, opened = %d, current = %d", qo, qc)
	}

	s, err := acc.GetSide(ticker)
	if s != LongSide {
		t.Errorf("expected side %s, got %s", LongSide, s)
	}

	notional := price * quantity / SCALE
	fee := notional / FEE_SCALE
	expAvail := VIRTUAL_BALANCE - notional - fee
	expBal := VIRTUAL_BALANCE - fee
	actualAvail := acc.GetAvailableBalance()
	actualBal := acc.GetBalance()
	if expAvail != actualAvail {
		t.Errorf("expected available balance = %d, got %d", expAvail, actualAvail)
	}
	if expBal != actualBal {
		t.Errorf("expected balance = %d, got %d", expBal, actualBal)
	}
}

func TestCreatePositionShort(t *testing.T) {
	acc := NewDemoAccount()

	ticker := "NVDA"
	price := uint64(200 * SCALE)
	quantity := uint64(3 * SCALE)

	err := acc.OpenPosition(ticker, ShortSide, price, quantity)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	l := acc.GetOpenedPositionsCount()
	if l != 1 {
		t.Errorf("expected one position, got %d", l)
	}

	b := acc.IsShortOpened(ticker)
	if !b {
		t.Errorf("expected opened short position on %s, but it's not opened", ticker)
	}

	po, err := acc.GetOpenPrice(ticker)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	pc, err := acc.GetCurrentPrice(ticker)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if po != price {
		t.Errorf("expected opened price = %d, got %d", price, po)
	}
	if po != pc {
		t.Errorf("expected opened price = current price, opened = %d, current = %d", po, pc)
	}

	qo, err := acc.GetOpenQuantity(ticker)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	qc, err := acc.GetCurrentQuantity(ticker)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if qo != quantity {
		t.Errorf("expected opened quantity = %d, got %d", quantity, qo)
	}
	if qo != qc {
		t.Errorf("expected opened quantity = current quantity, opened = %d, current = %d", qo, qc)
	}

	s, err := acc.GetSide(ticker)
	if s != ShortSide {
		t.Errorf("expected side %s, got %s", ShortSide, s)
	}

	notional := price * quantity / SCALE
	fee := notional / FEE_SCALE
	expAvail := VIRTUAL_BALANCE - notional - fee
	expBal := VIRTUAL_BALANCE - fee
	actualAvail := acc.GetAvailableBalance()
	actualBal := acc.GetBalance()
	if expAvail != actualAvail {
		t.Errorf("expected available balance = %d, got %d", expAvail, actualAvail)
	}
	if expBal != actualBal {
		t.Errorf("expected balance = %d, got %d", expBal, actualBal)
	}
}

func TestOpenTooMuch(t *testing.T) {
	acc := NewDemoAccount()

	ticker := "NVDA"
	price := uint64(200 * SCALE)
	quantity := uint64((VIRTUAL_BALANCE/price + 10) * SCALE)

	err := acc.OpenPosition(ticker, LongSide, price, quantity)
	if !errors.Is(err, NotEnoughBalanceErr) {
		t.Errorf("expected NotEnoughBalanceErr, got %s", err)
	}
}

func TestCloseTooMuch(t *testing.T) {
	acc := NewDemoAccount()

	ticker := "NVDA"
	price := uint64(200 * SCALE)
	quantity := uint64(10 * SCALE)

	err := acc.OpenPosition(ticker, LongSide, price, quantity)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = acc.ClosePosition(ticker, uint64(150*SCALE), floatMulToUint64(10.001, float64(SCALE)))
	if !errors.Is(err, NotEnoughAssetsErr) {
		t.Errorf("expected NotEnoughAssetsErr, got %s", err)
	}
}

func TestCloseFull(t *testing.T) {
	acc := NewDemoAccount()

	ticker := "NVDA"
	price := uint64(200 * SCALE)
	quantity := uint64(9 * SCALE)

	err := acc.OpenPosition(ticker, LongSide, price, quantity)
	if err != nil {
		t.Errorf("unexpected error at opening long: %s", err)
	}

	err = acc.CloseFullPosition(ticker, uint64(205*SCALE))
	if err != nil {
		t.Errorf("unexpected error at closing long: %s", err)
	}

	// closed all position, asserting it
	if acc.IsLongOpened(ticker) || acc.IsShortOpened(ticker) {
		t.Errorf("closed long: expected no opened positions")
	}
	q, err := acc.GetCurrentQuantity(ticker)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if q != 0 {
		t.Errorf("expected qty = 0, got %d", q)
	}
}

// full test for multiple operations on long position
// without fee:
// p_open=105, q_open=10 - opened long with price 105 and quantity 10
// p_close=110, q_close=2 - closing 2 shares with profit (110-105)*2 = +10, fixed profit total: +10, balance: realized+unrealized = 10+(110-105)*(10-2) = +50
// pc=100, qc=5 - closing 5 shares with loss (100-105)*5 = -25, fixed profit total: +10-25 = -15, balance: -15+(100-105)*(10-2-5) = -30
// pc=120, qc=3 - closing 3 shares with profit (120-105)*3 = +45, fixed profit total: -15+45 = +30, balance: 30+(120-105)*(10-2-5-3) = +30
// with fee (0.1% for example):
// po=105, qo=10, fee=105*10*0.001=1.05, profit: -1.05
// pc=110, qc=2, fee=110*2*0.001=0.22, fixed profit total: -1.05+10-0.22 = +8.73, balance: 8.73+(110-105)*(10-2) = +48.73
// pc=100, qc=5, fee=100*5*0.001=0.5, fixed profit total: +8.73-25-0.5 = -16.77, balance: -16.77+(100-105)*(10-2-5) = -31.77
// pc=120, qc=3, fee=120*3*0.001=0.36, fixed profit total: -16.77+45-0.36 = +27.87, balance: 27.87+(120-105)*(10-2-5-3) = +27.87
// with fee=1/FEE_SCALE
// po=105, qo=10, fee=105*10/FEE_SCALE = 1050/FEE_SCALE, profit: -1050/FEE_SCALE
// pc=110, qc=2, fee=110*2/FEE_SCALE = 220/FEE_SCALE, fixed profit total: -1050/FEE_SCALE+10-220/FEE_SCALE = +10-1270/FEE_SCALE, balance: +10-1270/FEE_SCALE+(110-105)*(10-2) = +50-1270/FEE_SCALE
// pc=100, qc=5, fee=100*5/FEE_SCALE = 500/FEE_SCALE, fixed profit total: +10-1270/FEE_SCALE-25-500/FEE_SCALE = -15-1770/FEE_SCALE, balance: -15-1770/FEE_SCALE+(100-105)*(10-2-5) = -30-1770/FEE_SCALE
// pc=120, qc=3, fee=120*3/FEE_SCALE = 360/FEE_SCALE, fixed profit total: -15-1770/FEE_SCALE+45-360/FEE_SCALE = +30-2130/FEE_SCALE, balance: +30-2130/FEE_SCALE+(120-105)*(10-2-5-3) = +30-2130/FEE_SCALE
func TestLongPosition(t *testing.T) {
	acc := NewDemoAccount()

	ticker := "NVDA"
	start := VIRTUAL_BALANCE
	po := uint64(105 * SCALE)
	qo := uint64(10 * SCALE)

	err := acc.OpenPosition(ticker, LongSide, po, qo)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	b := acc.IsLongOpened(ticker)
	if !b {
		t.Errorf("expected opened long position on %s, but it's not opened", ticker)
	}

	// po=105, qo=10, fee=105*10/FEE_SCALE = 1050/FEE_SCALE
	// available -= (notional + fee)
	expectedAvStart := start - po*qo/SCALE - 1050*SCALE/FEE_SCALE
	actual := acc.GetAvailableBalance()
	if expectedAvStart != actual {
		t.Errorf("expected start available balance = %d, got %d", expectedAvStart, actual)
	}
	// balance -= fee
	expected := start - 1050*SCALE/FEE_SCALE
	actual = acc.GetBalance()
	if expected != actual {
		t.Errorf("expected start balance = %d, got %d", expected, actual)
	}

	// closing 2 shares
	// available += pc*qc - fee
	// pc=110, qc=2, fee = 220/FEE_SCALE, balance: 50-1270/FEE_SCALE
	acc.ClosePosition(ticker, uint64(110*SCALE), uint64(2*SCALE))
	expectedAv1 := expectedAvStart + 110*2*SCALE - 220*SCALE/FEE_SCALE
	actual = acc.GetAvailableBalance()
	if expectedAv1 != actual {
		t.Errorf("closing 2 shares: expected available balance = %d, got %d", expectedAv1, actual)
	}
	expected = start + 50*SCALE - 1270*SCALE/FEE_SCALE
	actual = acc.GetBalance()
	if expected != actual {
		t.Errorf("closing 2 shares: expected balance = %d, got %d", expected, actual)
	}

	// closing 5 shares
	// pc=100, qc=5, fee=500/FEE_SCALE, balance: -30-1770/FEE_SCALE
	acc.ClosePosition(ticker, uint64(100*SCALE), uint64(5*SCALE))
	expectedAv2 := expectedAv1 + 100*5*SCALE - 500*SCALE/FEE_SCALE
	actual = acc.GetAvailableBalance()
	if expectedAv2 != actual {
		t.Errorf("closing 5 shares: expected available balance = %d, got %d", expectedAv2, actual)
	}
	expected = start - 30*SCALE - 1770*SCALE/FEE_SCALE
	actual = acc.GetBalance()
	if expected != actual {
		t.Errorf("closing 5 shares: expected balance = %d, got %d", expected, actual)
	}

	// closing 3 shares
	// pc=120, qc=3, fee=360/FEE_SCALE, balance: 30-2130/FEE_SCALE
	acc.ClosePosition(ticker, uint64(120*SCALE), uint64(3*SCALE))
	expectedAv3 := expectedAv2 + 120*3*SCALE - 360*SCALE/FEE_SCALE
	actual = acc.GetAvailableBalance()
	if expectedAv3 != actual {
		t.Errorf("closing 3 shares: expected available balance = %d, got %d", expectedAv3, actual)
	}
	expected = start + 30*SCALE - 2130*SCALE/FEE_SCALE
	actual = acc.GetBalance()
	if expected != actual {
		t.Errorf("closing 3 shares: expected balance = %d, got %d", expected, actual)
	}

	// closed all position, asserting it
	if acc.IsLongOpened(ticker) || acc.IsShortOpened(ticker) {
		t.Errorf("end: expected no opened positions")
	}
	q, err := acc.GetCurrentQuantity(ticker)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if q != 0 {
		t.Errorf("expected qty = 0, got %d", q)
	}
}

// full test for multiple operations on short position
// without fee:
// po=105, qo=10 - opened short with price 105 and quantity 10
// pc=100, qc=6 - closing 6 shares with profit (105-100)*6 = +30, fixed profit total: +30, balance: 30+(105-100)*(10-6) = +50
// pc=90, qc=2 - closing 2 shares with profit (105-90)*2 = +30, fixed profit total: 30+30 = +60, balance: 60+(105-90)*(10-6-2) = +90
// pc=110, qc=2 - closing last 2 shares with loss (105-110)*2 = -10, fixed profit total: 60-10 = +50, balance: 50+(105-110)*(10-6-2-2) = +50
// with fee=1/FEE_SCALE:
// po=105, qo=10, fee=105*10/FEE_SCALE=1050/FEE_SCALE, balance: -1050/FEE_SCALE
// pc=100, qc=6, fee=100*6/FEE_SCALE=600/FEE_SCALE, fixed profit total: -1050/FEE_SCALE+30-600/FEE_SCALE = 30-1650/FEE_SCALE, balance: 30-1650/FEE_SCALE+(105-100)*(10-6) = +50-1650/FEE_SCALE
// pc=90, qc=2, fee=90*2/FEE_SCALE=180/FEE_SCALE, fixed profit total: 30-1650/FEE_SCALE+30-180/FEE_SCALE = +60-1830/FEE_SCALE, balance: 60-1830/FEE_SCALE+(105-90)*(10-6-2) = +90-1830/FEE_SCALE
// pc=110, qc=2, fee=110*2/FEE_SCALE=220/FEE_SCALE, fixed profit total: 60-1830/FEE_SCALE-10-220/FEE_SCALE = +50-2050/FEE_SCALE, balance: 50-2050/FEE_SCALE+(105-110)*(10-6-2-2) = +50-2050/FEE_SCALE
func TestShortPosition(t *testing.T) {
	acc := NewDemoAccount()

	ticker := "NVDA"
	start := VIRTUAL_BALANCE
	po := uint64(105 * SCALE)
	qo := uint64(10 * SCALE)

	err := acc.OpenPosition(ticker, ShortSide, po, qo)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	b := acc.IsShortOpened(ticker)
	if !b {
		t.Errorf("expected opened short position on %s, but it's not opened", ticker)
	}

	// po=105, qo=10, fee=1050/FEE_SCALE, balance: -1050/FEE_SCALE
	// available -= (notional + fee)
	fee := 1050 * SCALE / FEE_SCALE
	expectedAvStart := start - po*qo/SCALE - fee
	actual := acc.GetAvailableBalance()
	if expectedAvStart != actual {
		t.Errorf("expected start available balance = %d, got %d", expectedAvStart, actual)
	}
	// balance -= fee
	expected := start - fee
	actual = acc.GetBalance()
	if expected != actual {
		t.Errorf("expected start balance = %d, got %d", expected, actual)
	}

	// closing 6 shares
	// pc=100, qc=6, fee=600/FEE_SCALE, balance: +50-1650/FEE_SCALE
	acc.ClosePosition(ticker, uint64(100*SCALE), uint64(6*SCALE))
	// available += baseCostOfClosedShares + profit - fee
	expectedAv1 := expectedAvStart + 105*6*SCALE + (105-100)*6*SCALE - 600*SCALE/FEE_SCALE
	actual = acc.GetAvailableBalance()
	if expectedAv1 != actual {
		t.Errorf("closing 6 shares: expected available balance = %d, got %d", expectedAv1, actual)
	}
	expected = start + 50*SCALE - 1650*SCALE/FEE_SCALE
	actual = acc.GetBalance()
	if expected != actual {
		t.Errorf("closing 6 shares: expected balance = %d, got %d", expected, actual)
	}

	// closing 2 shares
	// pc=90, qc=2, fee=180/FEE_SCALE, balance: +90-1830/FEE_SCALE
	acc.ClosePosition(ticker, uint64(90*SCALE), uint64(2*SCALE))
	expectedAv2 := expectedAv1 + 105*2*SCALE + (105-90)*2*SCALE - 180*SCALE/FEE_SCALE
	actual = acc.GetAvailableBalance()
	if expectedAv2 != actual {
		t.Errorf("closing 2 shares: expected available balance = %d, got %d", expectedAv2, actual)
	}
	expected = start + 90*SCALE - 1830*SCALE/FEE_SCALE
	actual = acc.GetBalance()
	if expected != actual {
		t.Errorf("closing 2 shares: expected balance = %d, got %d", expected, actual)
	}

	// closing last 2 shares
	// pc=110, qc=2, fee=220/FEE_SCALE, balance: +50-2050/FEE_SCALE
	acc.ClosePosition(ticker, uint64(110*SCALE), uint64(2*SCALE))
	expectedAv3 := expectedAv2 + 105*2*SCALE - (110-105)*2*SCALE - 220*SCALE/FEE_SCALE
	actual = acc.GetAvailableBalance()
	if expectedAv3 != actual {
		t.Errorf("closing last 2 shares: expected available balance = %d, got %d", expectedAv3, actual)
	}
	expected = start + 50*SCALE - 2050*SCALE/FEE_SCALE
	actual = acc.GetBalance()
	if expected != actual {
		t.Errorf("closing last 2 shares: expected balance = %d, got %d", expected, actual)
	}

	// closed all position, asserting it
	if acc.IsLongOpened(ticker) || acc.IsShortOpened(ticker) {
		t.Errorf("end: expected no opened positions")
	}
	q, err := acc.GetCurrentQuantity(ticker)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if q != 0 {
		t.Errorf("expected qty = 0, got %d", q)
	}
}
