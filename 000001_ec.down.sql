CREATE TABLE "consumers" (
  "id" INT PRIMARY KEY,
  "Name" VARCHAR(50) NOT NULL,
  "Email" VARCHAR(100) UNIQUE,
  "CreatedAt" TIMESTAMP DEFAULT (now())
);

CREATE TABLE "seats" (
  "id" INT PRIMARY KEY,
  "Section" VARCHAR(50) UNIQUE,
  "SeatNumber" VARCHAR(10),
  "SeatType" VARCHAR(20),
  "SeatStatus" BOOLEAN DEFAULT (false)
);

CREATE TABLE "events" (
  "id" INT PRIMARY KEY,
  "EventName" VARCHAR(100),
  "EventDate" datetime,
  "TotalSeats" INT
);

CREATE TABLE "tickets" (
  "id" INT PRIMARY KEY,
  "consumer_id" INT,
  "event_id" INT,
  "seat_id" INT,
  "PurchaseDate" TIMESTAMP DEFAULT (now()),
  "Price" DECIMAL(10,2),
  "TicketStatus" VARCHAR(20)
);

ALTER TABLE "tickets" ADD FOREIGN KEY ("consumer_id") REFERENCES "consumers" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("seat_id") REFERENCES "seats" ("id");
