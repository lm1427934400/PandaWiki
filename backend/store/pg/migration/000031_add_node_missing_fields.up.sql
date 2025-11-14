-- 添加Node模型中缺失的字段
ALTER TABLE nodes ADD COLUMN IF NOT EXISTS permissions jsonb DEFAULT '{}';
ALTER TABLE nodes ADD COLUMN IF NOT EXISTS creator_id text;
ALTER TABLE nodes ADD COLUMN IF NOT EXISTS editor_id text;
ALTER TABLE nodes ADD COLUMN IF NOT EXISTS edit_time timestamptz;
