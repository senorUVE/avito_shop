-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user"
(
    id       uuid         PRIMARY KEY,
    username    varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    token    text
);

COMMENT ON TABLE "user" IS 'Таблица пользователей';
COMMENT ON COLUMN "user".token IS 'Refresh токен пользователя';

CREATE TABLE "info"
(
    user_id uuid PRIMARY KEY REFERENCES "user"(id) ON DELETE CASCADE,
    coins   int NOT NULL DEFAULT 0 CHECK (coins >= 0)
);
COMMENT ON TABLE "info" IS 'Баланс монеток пользователей';

CREATE TABLE "inventory"
(
    user_id  UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    type     TEXT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, type)
);
COMMENT ON TABLE "inventory" IS 'Инвентарь пользователя';
COMMENT ON COLUMN "inventory".type IS 'Тип товара в инвентаре';
COMMENT ON COLUMN "inventory".quantity IS 'Количество товара';

CREATE TABLE "coin_transactions"
(
    id        SERIAL PRIMARY KEY,
    from_user UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    to_user   UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    amount    INTEGER NOT NULL CHECK (amount > 0)
);
COMMENT ON TABLE "coin_transactions" IS 'Транзакции перевода монет';
COMMENT ON COLUMN "coin_transactions".from_user IS 'ID отправителя';
COMMENT ON COLUMN "coin_transactions".to_user IS 'ID получателя';
COMMENT ON COLUMN "coin_transactions".amount IS 'Количество монет';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "coin_transactions" CASCADE;
DROP TABLE IF EXISTS "inventory" CASCADE;
DROP TABLE IF EXISTS "info" CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
-- +goose StatementEnd