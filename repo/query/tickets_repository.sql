-- name: CreateOrder :one
WITH seat_check AS (
  SELECT id FROM seats
  WHERE id = $3 
)
INSERT INTO tickets (
  consumer_id, event_id, seat_id, ticket_status, order_trade_no
)
SELECT
  $1, $2, $3, $4, $5
FROM seat_check
RETURNING *;

-- name: UpdateSeatStatus :one
UPDATE seats
SET seat_status = $2
WHERE id = $1 AND seat_status = true
RETURNING *;

-- name: SelectSeat :one
SELECT * FROM seats WhERE id = $1 ;

-- name: SelectSeatsListByEventAndSection :many
SELECT * 
FROM seats 
WHERE event_id = $1 
  AND section = $2
  AND seat_number = ANY($3::int[])
  FOR UPDATE;

-- name: SelectTicketListByOrderTradeNo :many
SELECT * FROM tickets WhERE order_trade_no = $1 ;

-- name: UpdateSeatStatusByOrderTradeNo :many
UPDATE seats
SET seat_status = $1
WHERE id IN (
    SELECT seat_id
    FROM tickets
    WHERE order_trade_no = $2
)
RETURNING *;

-- name: DeleteTicketsByOrderTradeNo :exec
DELETE FROM tickets
WHERE order_trade_no = $1
RETURNING *;

-- name: UpdateSeatStatusBySeatIDs :exec
UPDATE seats
SET seat_status = $1
WHERE id = ANY($2::int[]);
-- -- name: UpdateTicketStatus :many
-- UPDATE tickets SET ticket_status = $2
-- WHERE order_trade_no = $1 RETURNING *;

-- name: SelectSectionsAndPricesByEventID :many
SELECT section, price 
FROM (
    SELECT DISTINCT ON (section, price) section, price, id 
    FROM seats 
    WHERE event_id = $1
    ORDER BY section, price, id
) AS distinct_seats
ORDER BY id;



-- name: UpdateTicketStatus :many
WITH updated_tickets AS (
  UPDATE tickets
  SET ticket_status = $2
  WHERE order_trade_no = $1
  RETURNING *
) SELECT  updated_tickets.id, 
  updated_tickets.order_trade_no, 
  updated_tickets.purchase_date, 
  updated_tickets.ticket_status,
  updated_tickets.id AS consumer_id,
  consumers.email AS consumer_email,
  seats.section,
  seats.seat_number,
  seats.seat_status,
  seats.price FROM updated_tickets  
JOIN consumers  ON updated_tickets.consumer_id = consumers.id
JOIN seats  ON updated_tickets.seat_id = seats.id;
-- UPDATE tickets SET ticket_status = $2
-- WHERE order_trade_no = $1 JOIN seats ON tickets.seat_id = seats.id RETURNING *;

-- -- name: UpdateAuthorBios :exec
-- UPDATE authors SET bio = $1;