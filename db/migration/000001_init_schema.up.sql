CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "password" varchar NOT NULL,
  "company" bigserial NOT NULL,
  "access_level" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "companies" (
  "id" bigserial PRIMARY KEY,
  "company_type" varchar NOT NULL,
  "company_name" varchar NOT NULL,
  "owner" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "product_name" varchar NOT NULL,
  "origin" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "entity_id" bigserial,
  "product" bigserial,
  "amount" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "from_company_id" bigserial NOT NULL,
  "to_company_id" bigserial NOT NULL,
  "product" bigserial NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("first_name", "last_name");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "companies" ("company_name");

CREATE INDEX ON "companies" ("owner");

CREATE INDEX ON "products" ("product_name");

CREATE INDEX ON "products" ("product_name", "origin");

CREATE INDEX ON "entries" ("entity_id");

CREATE INDEX ON "entries" ("entity_id", "product");

CREATE INDEX ON "orders" ("from_company_id");

CREATE INDEX ON "orders" ("to_company_id");

CREATE INDEX ON "orders" ("from_company_id", "to_company_id");

ALTER TABLE "users" ADD FOREIGN KEY ("company") REFERENCES "companies" ("id");

ALTER TABLE "companies" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("entity_id") REFERENCES "companies" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("product") REFERENCES "products" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("from_company_id") REFERENCES "companies" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("to_company_id") REFERENCES "companies" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("product") REFERENCES "products" ("id");
