CREATE TABLE IF NOT EXISTS coins_history(
    from_user UUID NOT NULL REFERENCE users(id),
    to_user UUID NOT NULL REFERENCE users(id),
    coins_amount INT,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)