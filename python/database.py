import os
from dotenv import load_dotenv
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker
from sqlalchemy import text
import pandas as pd
from pathlib import Path
from models import Base
import logging

logger = logging.getLogger(__name__)

BASE_DIR = Path(__file__).resolve().parent.parent
load_dotenv(dotenv_path=f"{BASE_DIR}/.env")

DB_URL = os.getenv("DB_URL")
if DB_URL is None:
    e = "Environment variable DB_URL is not found"
    logger.error(e)
    raise Exception(e)

engine = create_async_engine(DB_URL)
AsyncSessionLocal = async_sessionmaker(bind=engine)

async def init_db():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)

async def fetch_query_to_df(query_str, params=None) -> pd.DataFrame:
    async with engine.connect() as conn:
        result = await conn.execute(text(query_str), params or {})
        rows = result.fetchall()
        columns = list(result.keys())
        df = pd.DataFrame(rows, columns=columns)
        
    return df

async def fetch_query_to_scalar(query_str, params=None):
    async with engine.connect() as conn:
        result = await conn.execute(text(query_str), params or {})
        return result.scalar()