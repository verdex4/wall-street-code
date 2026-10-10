import asyncio
import yfinance as yf

# define your message callback
def message_handler(message):
    print("Received message:", message)

async def main():
    # =======================
    # With Context Manager
    # =======================
    async with yf.AsyncWebSocket() as ws:
        await ws.subscribe(["AAPL"])
        await ws.listen(message_handler)

asyncio.run(main())