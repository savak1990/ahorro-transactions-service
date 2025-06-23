-- migrate:up
-- Initial schema for Ahorro Transactions Service

CREATE TABLE IF NOT EXISTS merchant (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS category (
    id UUID PRIMARY KEY,
    category_name TEXT NOT NULL,
    category_group TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    priority INTEGER
);

CREATE TABLE IF NOT EXISTS transaction (
    id UUID PRIMARY KEY,
    group_id UUID NOT NULL,
    user_id UUID NOT NULL,
    balance_id UUID NOT NULL,
    merchant_id UUID REFERENCES merchant(id),
    type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense', 'movement')),
    operation_id UUID,
    approved_at TIMESTAMPTZ NOT NULL,
    transacted_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS transaction_entry (
    id UUID PRIMARY KEY,
    transaction_id UUID NOT NULL REFERENCES transaction(id),
    description TEXT,
    amount NUMERIC(18,2) NOT NULL,
    category_id UUID REFERENCES category(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    budget_id UUID
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_transaction_user_id ON transaction(user_id);
CREATE INDEX IF NOT EXISTS idx_transaction_type ON transaction(type);
CREATE INDEX IF NOT EXISTS idx_transaction_transacted_at ON transaction(transacted_at);
CREATE INDEX IF NOT EXISTS idx_transaction_entry_transaction_id ON transaction_entry(transaction_id);
CREATE INDEX IF NOT EXISTS idx_transaction_entry_category_id ON transaction_entry(category_id);

-- migrate:down
-- Drop indexes first
DROP INDEX IF EXISTS idx_transaction_entry_category_id;
DROP INDEX IF EXISTS idx_transaction_entry_transaction_id;
DROP INDEX IF EXISTS idx_transaction_transacted_at;
DROP INDEX IF EXISTS idx_transaction_type;
DROP INDEX IF EXISTS idx_transaction_user_id;

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS transaction_entry;
DROP TABLE IF EXISTS transaction;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS merchant;
