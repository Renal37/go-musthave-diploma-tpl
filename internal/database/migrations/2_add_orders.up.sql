CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders (
    id          text PRIMARY KEY,                             -- Уникальный идентификатор заказа
    user_id     uuid REFERENCES users NOT NULL,              -- Идентификатор пользователя (внешний ключ на таблицу users)
    status      order_status NOT NULL DEFAULT 'NEW',         -- Статус заказа
    uploaded_at timestamp NOT NULL DEFAULT current_timestamp  -- Время загрузки заказа
);

CREATE INDEX idx_orders_user_id ON orders (user_id);

CREATE TABLE accrual_flow (
    id       uuid PRIMARY KEY DEFAULT uuid_generate_v4(),  -- Уникальный идентификатор потока начислений
    order_id text REFERENCES orders NOT NULL,                -- Идентификатор заказа (внешний ключ на таблицу orders)
    amount   numeric(15, 2) NOT NULL DEFAULT 0               -- Сумма начисления
);

CREATE INDEX idx_accrual_flow_order_id ON accrual_flow (order_id);
