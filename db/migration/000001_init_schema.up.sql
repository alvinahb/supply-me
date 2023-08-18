CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "company_id" bigserial NOT NULL,
  "role" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "companies" (
  "id" bigserial PRIMARY KEY,
  "company_type" varchar NOT NULL,
  "company_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "product_name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "company_id" bigserial NOT NULL,
  "product_id" bigserial NOT NULL,
  "amount" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "from_company_id" bigserial NOT NULL,
  "to_company_id" bigserial NOT NULL,
  "product_id" bigserial NOT NULL,
  "amount" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "inventories" (
  "id" bigserial PRIMARY KEY,
  "company_id" bigserial NOT NULL,
  "product_id" bigserial NOT NULL,
  "amount_available" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("first_name", "last_name");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "companies" ("company_name");

CREATE INDEX ON "products" ("product_name");

CREATE INDEX ON "entries" ("company_id");

CREATE INDEX ON "entries" ("company_id", "product_id");

CREATE INDEX ON "orders" ("from_company_id");

CREATE INDEX ON "orders" ("to_company_id");

CREATE INDEX ON "orders" ("from_company_id", "to_company_id");

CREATE INDEX ON "inventories" ("company_id");

CREATE INDEX ON "inventories" ("company_id", "product_id");

ALTER TABLE "users" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("from_company_id") REFERENCES "companies" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("to_company_id") REFERENCES "companies" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "inventories" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "inventories" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");