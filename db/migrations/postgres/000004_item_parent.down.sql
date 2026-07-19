-- 000004 reverse: drop parent_id constraints/index/column on `item`.

DROP INDEX IF EXISTS idx_item_parent;

ALTER TABLE item
    DROP CONSTRAINT IF EXISTS fk_item_parent;

ALTER TABLE item
    DROP COLUMN IF EXISTS parent_id;
