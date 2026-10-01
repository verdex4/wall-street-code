import yfinance as yf
import matplotlib.pyplot as plt

dat = yf.Ticker("NVDA")
hist = dat.history(period="1d", interval="1m")

hist["Close"].plot()
plt.show()