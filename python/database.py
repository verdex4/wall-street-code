import os
from dotenv import load_dotenv
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker
from sqlalchemy.dialects.postgresql import insert
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
        # not recreating existing tables
        await conn.run_sync(Base.metadata.create_all)

async def query_to_df(q, params=None) -> pd.DataFrame:
    async with engine.connect() as conn:
        result = await conn.execute(text(q), params or {})
        rows = result.fetchall()
        columns = list(result.keys())
        df = pd.DataFrame(rows, columns=columns)
        
    return df

async def query_to_scalar(q, params=None):
    async with engine.connect() as conn:
        result = await conn.execute(text(q), params or {})
        return result.scalar()

def insert_on_conflict_nothing(table, conn, keys, data_iter):
    data = [dict(zip(keys, row)) for row in data_iter]
    stmt = insert(table.table).values(data).on_conflict_do_nothing()
    conn.execute(stmt)

async def df_to_db(df: pd.DataFrame, table_name="candles"):
    """Inserts a DataFrame into a database table. No conflicts if data already exists."""
    async with engine.begin() as conn:
        await conn.run_sync(
            lambda sync_conn: df.to_sql(
            name=table_name,
            con=sync_conn,
            if_exists="append",
            index=False,
            method=insert_on_conflict_nothing,
            chunksize=5000,
            )
        )

async def execute_transaction(q, params=None):
    async with engine.begin() as conn:
        await conn.execute(text(q), params or {})

async def print_select_all():
    query = "SELECT * FROM candles;"
    df = await query_to_df(query)
    print(df)