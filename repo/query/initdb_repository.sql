-- name: CreateConsumer :one
INSERT INTO consumers (
  name, email
) VALUES (
  $1, $2
)
RETURNING id, name, email, created_at;

-- name: CreateSeat :one
INSERT INTO seats (
  section, seat_number, seat_status, price, event_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, section, seat_number, seat_status, price, event_id;

-- name: CreateEvent :one
INSERT INTO events (
  event_name, event_date, total_seats
) VALUES (
  $1, $2, $3
)
RETURNING id, event_name, event_date, total_seats;