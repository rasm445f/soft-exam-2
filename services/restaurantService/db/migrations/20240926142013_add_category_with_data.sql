-- +goose Up
-- +goose StatementBegin
-- Add `category` column to the `restaurant` table
ALTER TABLE restaurant
ADD COLUMN category VARCHAR(50);

-- Insert Mock data
INSERT INTO restaurant (name, address, rating, category)
VALUES
    ('Pizza Paradise', '123 Main Street', 4.5, 'pizza'),
    ('Burger Bonanza', '456 High Street', 4.2, 'burger'),
    ('Sushi Spot', '789 Ocean Avenue', 4.8, 'sushi'),
    ('Taco Town', '321 Fiesta Road', 3.7, 'mexican');

INSERT INTO menuitem (restaurantid, name, price, description)
VALUES
    (1, 'Pepperoni Pizza', 12.99, 'Classic pepperoni pizza with mozzarella cheese.'),
    (1, 'Margarita Pizza', 10.99, 'Fresh tomato, basil, and mozzarella cheese.'),
    (1, 'Hawaii Pizza', 12.99, 'The classic pizza with pineapple, ham and mozzarella cheese.'),
    (2, 'Classic Cheeseburger', 9.99, 'Juicy beefy patty with cheddar cheese, and lettuce.'),
    (2, 'Bacon Burger', 11.99, 'Beef patty with crispy bacon, cheese, and BBQ sauce.'),
    (3, 'California Roll', 8.99, 'Crab, avocado, and cucumber rolled with rice.'),
    (3, 'Spicy Tuna Roll', 9.99, 'Fresh tuna with spicy mayo.'),
    (3, 'Flamed Salmon Roll', 14.99, 'Fresh Flamed Salmon Roll, with tempura and rice.'),
    (4, 'Beef Tacos', 7.99, 'Soft tacos filled with seasoned beef and toppings.'),
    (4, 'Chicken Quesadilla', 8.99, 'Grilled chicken with cheese in a flour tortilla.'),
    (4, 'Spicy Chicken Taquitos', 10.99, 'Spicy Taco Roll with cheese and chicken.');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove `category` column from the `restaurant` table
ALTER TABLE restaurant
DROP COLUMN category;

-- Remove mock data
DELETE FROM menuitem WHERE restaurantid IN (1, 2, 3, 4)
DELETE FROM restaurant WHERE id in (1, 2, 3, 4)
-- +goose StatementEnd
