from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column
from sqlalchemy import BigInteger, Float, String, DateTime, Index, PrimaryKeyConstraint
from datetime import datetime

class Base(DeclarativeBase):
    pass

class Candle(Base):
    """Таблица, хранящая свечи в формате OHLCV с интервалом в 1 минуту.
    
    Колонки:
    - ticker: Тикер акции
    - opened_at: Время открытия свечи в UTC
    - open: Открытие свечи в центах
    - high: Максимум свечи в центах
    - low: Минимум свечи в центах
    - close: Закрытие свечи в центах
    - volume: Объем торгов в штуках
    - dividends: Дивиденды на одну акцию в центах
    - splits: Количество новых акций за одну старую
    """
    __tablename__ = "candles"    


    ticker: Mapped[str] = mapped_column(String(10), nullable=False)
    opened_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
    
    open: Mapped[int] = mapped_column(BigInteger, nullable=False)
    high: Mapped[int] = mapped_column(BigInteger, nullable=False)
    low: Mapped[int] = mapped_column(BigInteger, nullable=False)
    close: Mapped[int] = mapped_column(BigInteger, nullable=False)
    volume: Mapped[int] = mapped_column(BigInteger, nullable=False)

    dividends: Mapped[int] = mapped_column(BigInteger, nullable=False)
    splits: Mapped[float] = mapped_column(Float, nullable=False)

    __table_args__ = (
        PrimaryKeyConstraint("ticker", "opened_at", name="pk_candles"),
        Index("idx_ticker_opened_at", "ticker", "opened_at"),
    )