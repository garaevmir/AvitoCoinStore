CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    coins INT NOT NULL DEFAULT 1000
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id UUID REFERENCES users(id),
    to_user_id UUID REFERENCES users(id),
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE inventory (
    user_id UUID REFERENCES users(id),
    item_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, item_name)
);