CREATE TABLE IF NOT EXISTS summary (
                                       id SERIAL PRIMARY KEY,
                                       total_input DOUBLE PRECISION,
                                       total_output DOUBLE PRECISION,
                                       turnover DOUBLE PRECISION,
                                       total_buys DOUBLE PRECISION,
                                       total_sells DOUBLE PRECISION,
                                       portfolio_value DOUBLE PRECISION,
                                       commissions DOUBLE PRECISION,
                                       taxes DOUBLE PRECISION,
                                       net_stock_profit DOUBLE PRECISION,
                                       created_at TIMESTAMP DEFAULT now()
    );
