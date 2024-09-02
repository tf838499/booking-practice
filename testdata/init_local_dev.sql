begin;

INSERT INTO events (id, event_name, event_date, total_seats) VALUES
(1, 'Concert A', '2024-07-01 20:00:00', 100),
(2, 'Concert B', '2024-08-15 19:30:00', 150);

INSERT INTO consumers (id, name, email, created_at) VALUES
(1, 'John Doe', 'john.doe@example.com', now()),
(2, 'Jane Smith', 'jane.smith@example.com', now());

INSERT INTO seats (id, section, seat_number, seat_type, seat_status, price, event_id) VALUES
(1, 'A', 1, 'Regular', true, 100.00, 1),
(2, 'A', 2, 'Regular', true, 100.00, 1),
(3, 'A', 3, 'Regular', true, 100.00, 1),
(4, 'A', 4, 'Regular', true, 100.00, 1),
(5, 'B', 1, 'VIP', true, 200.00, 1),
(6, 'B', 2, 'VIP', true, 200.00, 1),
(7, 'B', 3, 'VIP', false, 200.00, 1),
(8, 'B', 4, 'VIP', true, 200.00, 1),
(9, 'C', 1, 'Regular', true, 100.00, 2),
(10, 'C', 2, 'Regular', false, 100.00, 2),
(11, 'C', 3, 'Regular', true, 100.00, 2),
(12, 'C', 4, 'Regular', true, 100.00, 2),
(13, 'D', 1, 'VIP', true, 200.00, 2),
(14, 'D', 2, 'VIP', false, 200.00, 2),
(15, 'D', 3, 'VIP', true, 200.00, 2),
(16, 'D', 4, 'VIP', true, 200.00, 2);


INSERT INTO tickets (consumer_id, event_id, seat_id, order_trade_no, ticket_status)
VALUES
(1, 1, 1, 'ORDER_TRADE_NO_1', 1),
(1, 1, 2, 'ORDER_TRADE_NO_1', 1),
(2, 2, 3, 'ORDER_TRADE_NO_2', 2),
(2, 2, 4, 'ORDER_TRADE_NO_2', 2);

end;