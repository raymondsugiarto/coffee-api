-- One-time fix for the morning session that loaded a non-seed Mawaru Matcha.
-- Run with your client of choice (psql, pgAdmin, TablePlus, etc).
-- Connection details: see config/resources/database.yml
--   host: localhost, port: 5432, user: postgres, pass: postgres, db: coffee
--
-- The morning session's item id is  47f2fe9e-c123-4a06-90ed-de5364a3ec88
-- The three children are the well-known CMM / SMM / MMM rows (also seeded).
-- This script flips their parent_id to point at the active parent.

UPDATE item
   SET parent_id = '47f2fe9e-c123-4a06-90ed-de5364a3ec88'   -- Mawaru Matcha (active)
 WHERE id = '7f99a6c7-99b9-435c-ba2a-58ea77760e8d'         -- Caramel Macchiato Matcha
    OR id = '9733b9e5-97c1-4135-829c-058fc282239c'         -- Strawberry Macchiato Matcha
    OR id = '963bcf80-b3a7-40fe-88c4-ff8c23a1d81d';        -- Mango Macchiato Matcha

-- Verify:
--   SELECT id, code, parent_id FROM item
--    WHERE id IN ('47f2fe9e-c123-4a06-90ed-de5364a3ec88',
--                 '7f99a6c7-99b9-435c-ba2a-58ea77760e8d',
--                 '9733b9e5-97c1-4135-829c-058fc282239c',
--                 '963bcf80-b3a7-40fe-88c4-ff8c23a1d81d');
