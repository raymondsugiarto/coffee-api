-- Optional sample mapping. Run AFTER 000004_item_parent migration.
-- These three "Macchiato Matcha" variants roll up under a single "Mawaru Matcha"
-- parent so the close-session UI auto-fills them when a driver loaded Matcha.
--
-- The parent/child relation is by `id`, NOT by `code`, because the catalog
-- may legitimately hold multiple items that share a `code` (e.g. one from
-- seed.sql, others created via the UI). The morning session records the
-- exact `item_id` of the parent it loaded — that's the UUID we must set on
-- the children's `parent_id` column.
--
-- Inspect first to confirm UUIDs:
--
--   SELECT id, code, name, parent_id, created_at FROM item
--    WHERE code IN ('MAM','CMM','SMM','MMM')
--    ORDER BY code, created_at DESC;
--
-- Then update by explicit `id`:

UPDATE item
   SET parent_id = '47f2fe9e-c123-4a06-90ed-de5364a3ec88'   -- Mawaru Matcha (parent)
 WHERE id = '7f99a6c7-99b9-435c-ba2a-58ea77760e8d'         -- Caramel Macchiato Matcha
    OR id = '9733b9e5-97c1-4135-829c-058fc282239c'         -- Strawberry Macchiato Matcha
    OR id = '963bcf80-b3a7-40fe-88c4-ff8c23a1d81d';        -- Mango Macchiato Matcha
