-- TODO: write two queries.
--
-- 1. LegalEntityByUUID :one  -- SELECT a single row from settlements.legal_entities by legal_entity_uuid.
-- 2. SaveLegalEntity :exec   -- INSERT into settlements.legal_entities, ON CONFLICT update the mutable fields.
--
-- After writing the queries, run `task generate` to regenerate dbmodels.

-- name: LegalEntityByUUID :one
SELECT * FROM settlements.legal_entities
WHERE legal_entity_uuid = $1;

-- name: SaveLegalEntity :exec
INSERT INTO settlements.legal_entities (
    legal_entity_uuid,
    legal_entity_type,
    business_name,
    address,
    tax_id,
    bank_account_number,
    currency
) VALUES (
    sqlc.arg(legal_entity_uuid),
    sqlc.arg(legal_entity_type),
    sqlc.arg(business_name),
    sqlc.arg(address),
    sqlc.arg(tax_id),
    sqlc.arg(bank_account_number),
    sqlc.arg(currency)
) ON CONFLICT (legal_entity_uuid) DO UPDATE SET
    business_name = EXCLUDED.business_name,
    tax_id = EXCLUDED.tax_id,
    address = EXCLUDED.address,
    bank_account_number = EXCLUDED.bank_account_number,
    currency = EXCLUDED.currency,
    updated_at = NOW();
