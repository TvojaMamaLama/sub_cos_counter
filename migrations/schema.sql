-- Subscriptions table
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cost DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('USD', 'RUB')),
    period_days INTEGER NOT NULL,
    next_payment DATE NOT NULL,
    category VARCHAR(50) NOT NULL,
    auto_renewal BOOLEAN NOT NULL DEFAULT true,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Payments table
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    subscription_id INTEGER NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('USD', 'RUB')),
    paid_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'completed' CHECK (status IN ('completed', 'pending', 'failed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for better performance
CREATE INDEX idx_subscriptions_active ON subscriptions(active);
CREATE INDEX idx_subscriptions_next_payment ON subscriptions(next_payment);
CREATE INDEX idx_subscriptions_category ON subscriptions(category);
CREATE INDEX idx_payments_subscription_id ON payments(subscription_id);
CREATE INDEX idx_payments_paid_at ON payments(paid_at);