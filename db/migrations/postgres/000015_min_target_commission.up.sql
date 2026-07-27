ALTER TABLE stock_session
    ADD COLUMN IF NOT EXISTS min_target_commission numeric(20, 4);

