-- 这里的顺序要和 up.sql 刚好相反
-- 先删外键约束（解除关系）
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";

-- 再删外键
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

-- 最后删表
DROP TABLE IF EXISTS "users";