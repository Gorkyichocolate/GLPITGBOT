DROP INDEX IF EXISTS idx_tickets_external_ticket_id;

ALTER TABLE tickets
DROP COLUMN IF EXISTS external_ticket_id;
