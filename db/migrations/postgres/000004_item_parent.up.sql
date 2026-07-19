-- ============================================================
-- Coffee Cart Driver Inventory Management Module
-- Migration 000004
-- Adds self-referencing parent_id to the `item` table so drivers can
-- group product variants under a parent product (e.g. a generic
-- "Matcha" parent with children "Matcha Mango" / "Matcha Strawberry").
-- ============================================================

ALTER TABLE item
    ADD COLUMN IF NOT EXISTS parent_id varchar(255) NULL;

ALTER TABLE item
    ADD CONSTRAINT fk_item_parent
        FOREIGN KEY (parent_id) REFERENCES item(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_item_parent ON item(parent_id);
