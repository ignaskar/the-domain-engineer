-- name: SaveBillingCycle :exec
INSERT INTO settlements.billing_cycles (billing_cycle_uuid, partner_uuid, partner_type, billing_cycle_number, closed, settled, start_date, end_date)
VALUES (
        sqlc.arg(billing_cycle_uuid),
        sqlc.arg(partner_uuid),
        sqlc.arg(partner_type),
        sqlc.arg(billing_cycle_number),
        sqlc.arg(closed),
        sqlc.arg(settled),
        sqlc.arg(start_date),
        sqlc.arg(end_date)
       )
ON CONFLICT (billing_cycle_uuid) DO UPDATE SET
   -- Other fields are immutable
    closed = EXCLUDED.closed,
    settled = EXCLUDED.settled,
    end_date = EXCLUDED.end_date;

-- name: OrdersByBillingCycleUUID :many
SELECT * FROM settlements.orders
INNER JOIN settlements.billing_cycle_orders USING (order_uuid)
WHERE billing_cycle_uuid = $1;

-- name: OrderBreakdownsByBillingCycleUUID :many
SELECT ob.*
FROM settlements.order_breakdowns ob
INNER JOIN settlements.billing_cycle_orders bco USING (order_uuid)
WHERE bco.billing_cycle_uuid = $1;

-- name: BillingCyclesByPartnerUUID :many
SELECT *
FROM settlements.billing_cycles
WHERE partner_uuid = $1
ORDER BY billing_cycle_number DESC;
