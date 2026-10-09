from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column
from sqlalchemy import BigInteger, Float, String, DateTime, Index, PrimaryKeyConstraint
from datetime import datetime

class Base(DeclarativeBase):
    pass

class Candle(Base):
    """Table, storing candles in OHLCV format with an interval of 1 minute.

    Columns:
    - opened_at: The time the candle opened in UTC
    - ticker: The ticker name of the stock
    - open: The opening price of the candle
    - high: The highest price of the candle
    - low: The lowest price of the candle
    - close: The closing price of the candle
    - volume: The volume of the candle in shares
    - dividends: The dividends paid
    - splits: The number of new shares issued on the candle

    Note, that price and dividends are multiplied by 1000, so it's representing value in one-hundredths of cent (0.0001$).

    """
    __tablename__ = "candles"    

    opened_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
    ticker: Mapped[str] = mapped_column(String(10), nullable=False)
    
    open: Mapped[int] = mapped_column(BigInteger, nullable=False)
    high: Mapped[int] = mapped_column(BigInteger, nullable=False)
    low: Mapped[int] = mapped_column(BigInteger, nullable=False)
    close: Mapped[int] = mapped_column(BigInteger, nullable=False)
    volume: Mapped[int] = mapped_column(BigInteger, nullable=False)

    dividends: Mapped[int] = mapped_column(BigInteger, nullable=True)
    splits: Mapped[float] = mapped_column(Float, nullable=True)

    __table_args__ = (
        PrimaryKeyConstraint("opened_at", "ticker", name="pk_candles"),
        Index("idx_opened_at_ticker", "opened_at", "ticker"),
    )