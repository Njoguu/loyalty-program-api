CREATE TABLE admins (
    id SERIAL NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    username text,
    password text,
    is_admin boolean DEFAULT true
);

CREATE TABLE password_resets (
    id SERIAL NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    username text,
    email_address text,
    password_reset_code text,
    expired_at timestamp with time zone DEFAULT (now() + '00:15:00'::interval)
);

CREATE TABLE products (
    id SERIAL NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text,
    description text,
    price integer,
    quantity integer
);

CREATE TABLE transactions (
    id SERIAL NOT NULL,
    username text,
    points integer,
    state text,
    description text,
    date date,
    "time" time without time zone,
    product text
);

CREATE TABLE users (
    id SERIAL NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    username text UNIQUE,
    first_name text,
    last_name text,
    gender text,
    email_address text,
    password text,
    city text,
    phone_number text,
    redeemable_points integer,
    is_email_verified boolean DEFAULT false NOT NULL,
    virtual_card_number text
);

CREATE TABLE verify_emails (
    id SERIAL NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    username text,
    email_address text,
    secret_code text,
    expired_at timestamp with time zone DEFAULT (now() + '00:15:00'::interval)
);

ALTER TABLE "transactions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "password_resets" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");