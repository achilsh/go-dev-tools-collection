-- KEYS[1] = token_key
-- KEYS[2] = quota_key
-- ARGV[1] = max_tokens
-- ARGV[2] = refill_rate_per_sec
-- ARGV[3] = max_daily_requests
-- ARGV[4] = current_timestamp (in seconds)

-- return -1: Too many requests as (rate limit)
-- return -2: Daily request quota exceeded.
-- return 1: allow.

local token_key = KEYS[1]
local quota_key = KEYS[2]
local max_tokens = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local max_daily = tonumber(ARGV[3])
local now = tonumber(ARGV[4])

-- 获取上次的 token 数和时间戳
local bucket = redis.call('HMGET', token_key, 'tokens', 'timestamp')
local tokens = tonumber(bucket[1])
local last_refill = tonumber(bucket[2])

if tokens == nil then
    tokens = max_tokens
    last_refill = now
end

-- 计算时间差并补充 token
local delta = math.max(0, now - last_refill)
local refill = delta * refill_rate
tokens = math.min(max_tokens, tokens + refill)
if tokens < 1 then
    return -1 -- 频率限制：无 token 可用
end

-- 总量判断
local current_daily = tonumber(redis.call('GET', quota_key) or '0')
if current_daily + 1 > max_daily then
    return -2 -- 总量限制：超出每日配额
end

-- 更新 token 和 quota
tokens = tokens - 1
redis.call('HMSET', token_key, 'tokens', tokens, 'timestamp', now)
redis.call('EXPIRE', token_key, 3600) -- 设置 token_key TTL
redis.call('INCR', quota_key)
if current_daily == 0 then
    redis.call('EXPIRE', quota_key, 86400 - (now % 86400))
end

return 1