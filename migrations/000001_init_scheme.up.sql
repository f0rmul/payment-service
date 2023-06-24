
create table if not exists accounts
(   
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id  UUID DEFAULT NULL,
    balance bigint NOT NULL DEFAULT 0, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS account_idx ON accounts (owner_id);


create table if not exists products
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(255) NOT NULL UNIQUE
);
CREATE INDEX IF NOT EXISTS product_idx ON products (name);


create table if not exists reservations
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    product_id UUID NOT NULL,
    order_id   UUID NOT NULL UNIQUE,
    amount bigint   NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    foreign key (account_id) references accounts (id),
    foreign key (product_id) references products (id)
);
CREATE INDEX IF NOT EXISTS reservation_idx ON reservations (account_id);


create table if not exists operations
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    amount bigint NOT NULL,
    operation_type varchar(255) NOT NULL,
    product_id UUID DEFAULT NULL,
    order_id UUID DEFAULT NULL,
    description varchar(255) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    foreign key (account_id) references accounts (id)
);
CREATE INDEX IF NOT EXISTS operation_idx ON operations (account_id);