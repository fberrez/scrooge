CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS account;
DROP TABLE IF EXISTS transaction;

CREATE TABLE account (
    uuid UUID NOT NULL,
    balance INT NOT NULL DEFAULT 0,
    PRIMARY KEY (uuid)
);

CREATE TABLE transaction (
    uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    description VARCHAR,
    currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    PRIMARY KEY (uuid),
    FOREIGN KEY (account_id) REFERENCES account(uuid)
);

INSERT INTO account(uuid, balance) VALUES 
    ('a84f1a1b-d6eb-4819-be29-2055b8862094', 100),
    ('7f87b08e-760b-43df-9af4-5354da34e7b4', 500);



