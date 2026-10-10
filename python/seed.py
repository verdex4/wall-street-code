import yfinance as yf
from database import df_to_db, query_to_scalar
from datetime import datetime, timedelta
import pandas as pd
import asyncio
from zoneinfo import ZoneInfo
import logging

logger = logging.getLogger(__name__)

async def download_realtime(tickers: list[str]):
    """Fetches last 1 min candles from REST for selected tickers, saving them to database and working while session is open.

    For now, websocket is omitted due to many tickers and hard implementation.
    """
    cur_time = datetime.now(tz=ZoneInfo("America/New_York"))
    logger.info(f"Starting downloading realtime data. Current time: {cur_time}")
    if cur_time.weekday() == 5 or cur_time.weekday() == 6:
        open_datetime = cur_time.replace(hour=9, minute=30, second=0, microsecond=0)
        while open_datetime.weekday != 0:
            open_datetime += timedelta(days=1)
        delta = open_datetime - cur_time
        sleep_time = delta.total_seconds()
        logger.info(f"It's weekend, markets are closed, sleeping for {delta} until {open_datetime}")
        await asyncio.sleep(sleep_time)
    if cur_time < datetime(cur_time.year, cur_time.month, cur_time.day, 9, 30, 0):
        open_datetime = cur_time.replace(hour=9, minute=30, second=0, microsecond=0)
        delta = open_datetime - cur_time
        sleep_time = delta.total_seconds()
        logger.info(f"It's {cur_time}, markets are not opened, waiting {delta} until open at {open_datetime}")
        await asyncio.sleep(sleep_time)
    if cur_time >= datetime(cur_time.year, cur_time.month, cur_time.day, 16, 0, 0):
        open_datetime = cur_time.replace(day=cur_time.day + 1, hour=9, minute=30, second=0, microsecond=0)
        while open_datetime.weekday == 5 or open_datetime.weekday == 6:
            open_datetime += timedelta(days=1)
        delta = open_datetime - cur_time
        sleep_time = delta.total_seconds()
        logger.info(f"Markets are closed, waiting {delta} until open at {open_datetime}")
        await asyncio.sleep(sleep_time)
    if cur_time.replace(second=0, microsecond=0, tzinfo=None) == datetime(cur_time.year, cur_time.month, cur_time.day, 9, 30, 0):
        logger.info("Session is opened, but need to wait for first candle. Waiting...")
        await asyncio.sleep(60 - cur_time.second)

    cur_time = datetime.now(tz=ZoneInfo("America/New_York"))
    while cur_time < datetime(cur_time.year, cur_time.month, cur_time.day, 16, 0, 0):
        logger.debug(f"Downloading realtime data. Current time: {cur_time}")
        await fetch_history(tickers, fast_download=True)
        cur_time = datetime.now(tz=ZoneInfo("America/New_York"))

    logger.info("Session closed. Realtime downloading is finished")
    open_datetime = cur_time.replace(day=cur_time.day + 1, hour=9, minute=30, second=0, microsecond=0)
    while open_datetime.weekday == 5 or open_datetime.weekday == 6:
        open_datetime += timedelta(days=1)
    delta = open_datetime - cur_time
    sleep_time = delta.total_seconds()
    logger.info(f"Waiting {delta} until open at {open_datetime}")
    await asyncio.sleep(sleep_time)
    await download_realtime(tickers)

async def fetch_history(tickers: list[str], fast_download: bool = False):
    """Fetches history data (1 min candles) by 1 day for selected tickers that not in database yet.
    
    NOTE: if database is empty, FILLING 10 TICKERS WILL TAKE AROUND 30 MINUTES. It will save from API rate limits.
    There is no official rate limits, but algorithm making around 700 requests per hour if loading history.

    If downloading data for last minutes, prefer using fast_download=True to make requests faster (small requests will not cause rate limits as fast as history loading).
    With fast_download=True number of requests per hour will be around 1400.
    """
    logger.info("Starting filling database")

    for ticker in tickers:
        cur_time = datetime.now(tz=ZoneInfo("America/New_York"))
        last_candle = await query_to_scalar(
            "SELECT MAX(opened_at) FROM candles WHERE ticker = :ticker;",
            params={"ticker": ticker}
        )
        if last_candle is not None:
            last_candle = last_candle.astimezone(ZoneInfo("America/New_York"))
            if last_candle >= cur_time.replace(second=0, microsecond=0) - timedelta(minutes=1):
                logger.info(f"Ticker {ticker} is up to date, skipping")
                continue
       
        # intraday data available for 30 days
        # finding last available day or continue to load existing data
        # session opening at 9:30, closing at 16:00
        if cur_time.replace(tzinfo=None) < datetime(cur_time.year, cur_time.month, cur_time.day, 9, 30, 0):
            available = cur_time - timedelta(days=30)
            if last_candle is not None and last_candle >= available:
                if last_candle.hour == 15 and last_candle.minute == 59:
                    start = datetime(last_candle.year, last_candle.month, last_candle.day + 1, 9, 30, 0)
                else:
                    start = (last_candle + timedelta(minutes=1)).replace(tzinfo=None)
            else:
                start = datetime(available.year, available.month, available.day, 9, 30, 0)
            # end is exclusive from yf.Ticker.history() method
            end = datetime(cur_time.year, cur_time.month, cur_time.day - 1, 16, 0, 0)
        elif datetime(cur_time.year, cur_time.month, cur_time.day, 9, 30, 0) < cur_time.replace(tzinfo=None) < datetime(cur_time.year, cur_time.month, cur_time.day, 16, 0, 0):
            available = cur_time - timedelta(days=30)
            if last_candle is not None and last_candle >= available:
                if last_candle.hour == 15 and last_candle.minute == 59:
                    start = datetime(last_candle.year, last_candle.month, last_candle.day + 1, 9, 30, 0)
                else:
                    start = (last_candle + timedelta(minutes=1)).replace(tzinfo=None)
            else:
                start = (available + timedelta(minutes=1)).replace(second=0, microsecond=0, tzinfo=None)
            end = cur_time.replace(second=0, microsecond=0, tzinfo=None)
        else:
            available = (cur_time - timedelta(days=29)).replace(hour=9, minute=30, second=0, microsecond=0)
            if last_candle is not None and last_candle >= available:
                if last_candle.hour == 15 and last_candle.minute == 59:
                    start = datetime(last_candle.year, last_candle.month, last_candle.day + 1, 9, 30, 0)
                else:
                    start = (last_candle + timedelta(minutes=1)).replace(tzinfo=None)
            else:
                start = available.replace(tzinfo=None)
            end = datetime(cur_time.year, cur_time.month, cur_time.day, 16, 0, 0)

        if start >= end:
            logger.info(f"Ticker {ticker} is up to date, skipping")
            continue

        t = yf.Ticker(ticker)
        last_time = start
        delta = end - start

        logger.info(f"Starting download history data for ticker {ticker}. Start: {start}, end: {end}")
        while start < end:
            try:
                # if it is weekend -> skip
                if start.weekday() == 5 or start.weekday() == 6:
                    logger.debug(f"Skipping weekend day {start}")
                    start += timedelta(days=1)
                    continue
                if delta <= timedelta(days=1):
                    # requesting minutes we need, not all day
                    df = t.history(interval="1m", start=start, end=end)
                else:
                    df = t.history(period="1d", interval="1m", start=start)
                if df.empty:
                    # it can be holiday so continue
                    logger.debug(f"No data loaded for ticker {ticker} with start {start}. Continuing to next day")
                    start += timedelta(days=1)
                    if fast_download:
                        logger.debug("Using fast download, skipping 2.5 seconds")
                        await asyncio.sleep(2.5)
                    else:
                        logger.debug(f"Slow mode enabed, sleeping for 5 seconds")
                        await asyncio.sleep(5)
                    continue
                last_time = df.index[-1].tz_localize(None)
                logger.debug(f"Downloaded {df.shape[0]} candles for ticker {ticker} from {start} until {last_time}")

                # parsing df
                df = df.reset_index()
                df = df.rename(
                    columns={
                        "Datetime": "opened_at",
                        "Open": "open", 
                        "High": "high", 
                        "Low": "low", 
                        "Close": "close", 
                        "Volume": "volume",
                        "Dividends": "dividends",
                        "Stock Splits": "splits",
                    }
                )
                df["ticker"] = ticker
                # convert to int (price are in one-hundredths of cent, see models.py for details)
                df["open"] = (round(df["open"] * 10**4)).astype("int64")
                df["high"] = (round(df["high"] * 10**4)).astype("int64")
                df["low"] = (round(df["low"] * 10**4)).astype("int64")
                df["close"] = (round(df["close"] * 10**4)).astype("int64")
                df["dividends"] = (round(df["dividends"] * 10**4)).astype("int64")

                logger.debug("Saving data to database")
                await df_to_db(df, "candles")
                logger.debug("Data saved to database")

                start += timedelta(days=1)
                if fast_download:
                    logger.debug("Using fast download, skipping 2.5 seconds")
                    await asyncio.sleep(2.5)
                else:
                    logger.debug(f"Slow mode enabed, sleeping for 5 seconds")
                    await asyncio.sleep(5)
            except Exception as e:
                if "Too Many Requests" in str(e):
                    logger.critical(f"CRITICAL ERROR: API RATE LIMITED. Text: {e}. Please try again later after few hours.")
                    return
                else:
                    logger.error(f"Error while downloading data: {e}")
                    return

        logger.info(f"Finished downloading ticker {ticker} at {last_time}.")

    logger.info("All data downloaded successfully")
    