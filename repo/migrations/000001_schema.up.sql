CREATE TABLE "consumers" (
  "id" SERIAL  PRIMARY KEY,
  "name" VARCHAR(50) NOT NULL,
  "email" VARCHAR(100) UNIQUE NOT NULL,
  "created_at" TIMESTAMP DEFAULT (now()) NOT NULL
);

CREATE TABLE "seats" (
  "id" SERIAL  PRIMARY KEY,
  "section" VARCHAR(50)  NOT NULL,
  "seat_number" INT NOT NULL,
  "seat_status" BOOLEAN DEFAULT (true) NOT NULL,
  "price" Float NOT NULL ,
  "event_id" INT NOT NULL
);

CREATE TABLE "events" (
  "id" SERIAL  PRIMARY KEY,
  "event_name" VARCHAR(100) NOT NULL,
  "event_date" timestamp NOT NULL,
  "total_seats" INT NOT NULL
);

CREATE TABLE "tickets" (
  "id" SERIAL  PRIMARY KEY,
  "consumer_id" INT NOT NULL,
  "event_id" INT NOT NULL,
  "seat_id" INT NOT NULL,
  "order_trade_no" VARCHAR(100) NOT NULL,
  "purchase_date" TIMESTAMP DEFAULT (now()) NOT NULL,
  "ticket_status" INT NOT NULL
);

ALTER TABLE "tickets" ADD FOREIGN KEY ("consumer_id") REFERENCES "consumers" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("seat_id") REFERENCES "seats" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");

ALTER TABLE "seats" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");

-- 重置序列
TRUNCATE tickets RESTART IDENTITY;
-- ALTER SEQUENCE tickets_id_seq RESTART WITH 1;