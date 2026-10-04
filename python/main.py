import yfinance as yf

def test():
    dat = yf.Ticker("AAPL")
    hist = dat.history(period="1d", interval="1m")

    print(hist)
    #hist["Close"].plot()
    #plt.show()

def main():
    print("Hello, World!")
    
if __name__ == "__main__":
    main()