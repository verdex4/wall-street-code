import yfinance as yf
from database import init_db, print_select_all, execute_transaction, query_to_scalar, query_to_df
from seed import fetch_history, download_realtime
from sqlalchemy import text
import logging
import asyncio
import matplotlib.pyplot as plt
import pandas as pd

logging.basicConfig(level=logging.INFO, format='%(filename)s -> %(funcName)s() - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# top 60 tickers by liquidity on Nasdaq and NYSE at 2026-10-08
# at competition available only stocks, without futures
# it is a little looking ahead bias, but there is not much data (only month+)
tickers = [
    "NVDA", "AAPL", "GOOG", "GOOGL", "MSFT", 
    "AMZN", "SPCX", "META", "AVGO", "TSLA", 
    "MU", "LLY", "AMD", "BRK-B", "BRK-A", 
    "WMT", "JPM", "V", "XOM", "JNJ",
    "INTC", "MA", "ABBV", "PLTR", "CSCO",
    "COST", "CVX", "ORCL", "AMAT", "LRCX",
    "KO", "BAC", "CAT", "DELL", "PG",
    "MRK", "UNH", "PANW", "PM", "GE",
    "NFLX", "MS", "HD", "CRWD", "ANET",
    "GEV", "TXN", "GS", "KLAC", "RTX",
    "WFC", "TMO", "SNDK", "MRVL", "AMGN",
    "C", "IBM", "APH", "AXP", "ADI"
]

async def main():
    logger.info("Python app is starting")
    
    await init_db()
    await fetch_history(tickers)
    await download_realtime(tickers)

    df = await query_to_df("SELECT * FROM candles WHERE ticker = 'NVDA' ORDER BY opened_at;")
    plt.plot(df["close"] / 10000)
    plt.xlabel("Candles")
    plt.ylabel("Price")
    plt.title("NVDA")
    plt.savefig('/app/plots/nvda_price_chart_26_09_09-26_10_09.png')

    logger.info("End of Python app")
    
if __name__ == "__main__":
    asyncio.run(main())