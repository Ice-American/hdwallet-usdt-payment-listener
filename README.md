# hdwallet-usdt-payment-listener

-- 原有的 addresses 表只用于记录 HD 派生根信息，不再存子地址
CREATE TABLE addresses (
  id SERIAL PRIMARY KEY,
  user_id TEXT NOT NULL,
  mnemonic_hash TEXT NOT NULL,  -- 可选，用于多根管理
  created_at TIMESTAMP DEFAULT NOW()
);

-- 新增 deposits 表，每次请求都插一条
CREATE TABLE deposits (
  id SERIAL PRIMARY KEY,
  user_id TEXT NOT NULL,
  derivation_index INTEGER NOT NULL,
  sub_address TEXT NOT NULL,
  request_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_id, sub_address)
);

-- 监听到真正的 on-chain 充值后，往 payments 表写流水
CREATE TABLE payments (
  id SERIAL PRIMARY KEY,
  deposit_id INT REFERENCES deposits(id),
  tx_hash TEXT NOT NULL,
  amount NUMERIC NOT NULL,
  confirmed_at TIMESTAMP DEFAULT NOW()
);
