import yfinance as yf
from database import init_db, AsyncSessionLocal
from sqlalchemy import text
import logging
import asyncio

logging.basicConfig(level=logging.INFO, format='%(filename)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

def test1():
    t = yf.Ticker("AAPL")
    hist = t.history(period="5d", interval="1d", start="2020-08-29")
    print(hist)

async def test_db():
    async with AsyncSessionLocal() as async_session:
        res = await async_session.execute(text("SELECT COUNT(*) FROM candles"))
        print(f"result: {res.scalar()}")

async def main():
    await init_db()
    await test_db()
    
if __name__ == "__main__":
    asyncio.run(main())