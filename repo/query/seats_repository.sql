-- name: SelectSeatsList :many
SELECT * FROM seats WhERE event_id = $1 ;

-- name: 