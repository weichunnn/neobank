create table "users" (
  "username" varchar primary key,
  "hashed_password" varchar not null,
  "full_name" varchar not null, 
  "email" varchar UNIQUE not null, 
  "password_created_at" timestamptz not null default '0001-01-01 00:00:00Z',
  "created_at" timestamptz not null default (now())
);

alter table "accounts" add constraint "account_owner_fkey" foreign key ("owner") references "users" ("username");

-- index automatically created by postgres
alter table "accounts" add constraint "owner_currency_key" unique ("owner", "currency");