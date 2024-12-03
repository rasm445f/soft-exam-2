-- +goose Up
-- +goose StatementBegin
ALTER TABLE "Order"
    ADD CONSTRAINT fk_delivery_agent
    FOREIGN KEY (DeliveryAgentID)
    REFERENCES DeliveryAgent (ID)
    ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "Order" DROP CONSTRAINT fk_delivery_agent;
-- +goose StatementEnd
