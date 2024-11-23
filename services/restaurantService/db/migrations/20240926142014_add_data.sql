-- +goose up
-- +goose StatementBegin

ALTER TABLE Restaurant
DROP COLUMN Address,
ADD COLUMN street VARCHAR(255),
ADD COLUMN zip_code INT REFERENCES Address(zip_code);

INSERT INTO Address (zip_code, city) VALUES
    (2800, 'Kgs. Lyngby'),
    (2970, 'Hørsholm'),
    (2980, 'Kokkedal');

INSERT INTO Restaurant (name, street, zip_code, rating, category) VALUES
    ('Pizza Paradise', 'Lyngby Hovedgade 25', 2800, 4.5, 'pizza'),
    ('Taco Town', 'Lyngbyvej 18', 2800, 3.7, 'mexican'),
    ('Café Aroma', 'Kongensgade 15', 2800, 4.5, 'coffee'),
    ('Lyngby Sushi', 'Gammel Torv 7', 2800, 4.8, 'sushi'),
    ('Burger Bonanza', 'Hovedgaden 10', 2970, 4.2, 'burger'),
    ('Hørsholm Grill', 'Vedbækvej 10', 2970, 4.2, 'grill'),
    ('Pizza Palace', 'Teglvej 20', 2970, 4.0, 'pizza'),
    ('Italian Feast', 'Strandvejen 15', 2970, 4.3, 'italian'),
    ('Sushi Spot', 'Stationsvej 25', 2980, 4.8, 'sushi'),
    ('Kokkedal Bites', 'Enghavevej 33', 2980, 3.9, 'burger');

INSERT INTO Menuitem (restaurantid, name, price, description) VALUES
    (1, 'Pepperoni Pizza', 12.99, 'Classic pepperoni pizza with mozzarella cheese.'),
    (1, 'Margarita Pizza', 10.99, 'Fresh tomato, basil, and mozzarella cheese.'),
    (1, 'Hawaiian Pizza', 13.99, 'Pineapple, ham, and mozzarella cheese.'),
    (1, 'BBQ Chicken Pizza', 14.99, 'Grilled chicken with BBQ sauce and onions.'),
    (1, 'Veggie Delight Pizza', 11.99, 'Mushrooms, peppers, olives, and onions.'),
    (1, 'Four Cheese Pizza', 15.99, 'Mozzarella, cheddar, parmesan, and gorgonzola.'),
    (1, 'Meat Lovers Pizza', 16.99, 'Pepperoni, sausage, bacon, and ham.'),
    (1, 'Calzone', 9.99, 'Folded pizza with ricotta, mozzarella, and ham.'),
    (1, 'Garlic Bread', 4.99, 'Classic garlic bread with melted cheese.'),
    (1, 'Tiramisu', 6.99, 'Italian dessert with mascarpone and espresso.'),

    (2, 'Beef Tacos', 7.99, 'Soft tacos filled with seasoned beef and toppings.'),
    (2, 'Chicken Quesadilla', 8.99, 'Grilled chicken with cheese in a flour tortilla.'),
    (2, 'Spicy Chicken Taquitos', 9.99, 'Spicy chicken rolled in crispy tortillas.'),
    (2, 'Loaded Nachos', 10.99, 'Tortilla chips with cheese, jalapeños, and guacamole.'),
    (2, 'Burrito Bowl', 12.99, 'Rice, beans, chicken, and toppings in a bowl.'),
    (2, 'Shrimp Tacos', 11.99, 'Grilled shrimp in soft tortillas with slaw.'),
    (2, 'Vegetarian Burrito', 9.99, 'Beans, rice, veggies, and salsa in a tortilla.'),
    (2, 'Guacamole & Chips', 5.99, 'Fresh guacamole served with tortilla chips.'),
    (2, 'Mexican Street Corn', 4.99, 'Corn on the cob with cheese, lime, and chili.'),
    (2, 'Churros', 6.99, 'Fried dough sticks with cinnamon sugar.'),

    (3, 'Cappuccino', 3.99, 'Espresso with steamed milk and foam.'),
    (3, 'Latte', 4.49, 'Smooth espresso with steamed milk.'),
    (3, 'Espresso', 2.99, 'Rich and bold shot of coffee.'),
    (3, 'Americano', 3.49, 'Espresso diluted with hot water.'),
    (3, 'Mocha', 4.99, 'Espresso, chocolate, and steamed milk.'),
    (3, 'Flat White', 4.49, 'Espresso with velvety steamed milk.'),
    (3, 'Matcha Latte', 4.99, 'Green tea matcha with steamed milk.'),
    (3, 'Cold Brew Coffee', 4.29, 'Smooth, cold-brewed coffee.'),
    (3, 'Blueberry Muffin', 2.99, 'Freshly baked blueberry muffin.'),
    (3, 'Croissant', 3.49, 'Buttery and flaky croissant.'),

    (4, 'Eel Roll', 11.99, 'Eel and cucumber with sweet eel sauce.'),
    (4, 'Avocado Roll', 7.99, 'Avocado wrapped in seasoned rice and seaweed.'),
    (4, 'Spicy Salmon Roll', 10.99, 'Salmon with spicy mayo and cucumber.'),
    (4, 'Tempura Roll', 12.99, 'Fried shrimp with avocado and cucumber.'),
    (4, 'Sake Sashimi', 13.99, 'Thinly sliced fresh salmon.'),
    (4, 'Unagi Don', 15.99, 'Grilled eel over rice with sauce.'),
    (4, 'Miso Soup', 3.99, 'Traditional Japanese soup with tofu.'),
    (4, 'Seaweed Salad', 4.99, 'Seasoned seaweed with sesame dressing.'),
    (4, 'Tuna Tartare', 12.99, 'Fresh tuna with soy and avocado.'),
    (4, 'Green Tea Ice Cream', 4.99, 'Creamy matcha-flavored ice cream.'),

    (5, 'Classic Cheeseburger', 9.99, 'Juicy beef patty with cheddar cheese.'),
    (5, 'Bacon Burger', 11.99, 'Beef patty with crispy bacon and BBQ sauce.'),
    (5, 'Veggie Burger', 8.99, 'Plant-based patty with lettuce and tomato.'),
    (5, 'Double Cheeseburger', 12.99, 'Two beef patties with double cheese.'),
    (5, 'Mushroom Swiss Burger', 11.49, 'Beef patty with mushrooms and Swiss cheese.'),
    (5, 'Chicken Sandwich', 9.49, 'Grilled chicken breast with lettuce and mayo.'),
    (5, 'BBQ Pulled Pork Sandwich', 10.99, 'Slow-cooked pulled pork with BBQ sauce.'),
    (5, 'Loaded Fries', 6.99, 'Fries topped with cheese, bacon, and ranch.'),
    (5, 'Onion Rings', 4.99, 'Crispy battered onion rings.'),
    (5, 'Milkshake', 4.99, 'Vanilla, chocolate, or strawberry milkshake.'),

    (6, 'Grilled Chicken', 14.99, 'Tender grilled chicken with sides.'),
    (6, 'Pork Ribs', 18.99, 'Slow-cooked ribs with BBQ sauce.'),
    (6, 'Steak Sandwich', 12.99, 'Juicy steak in a toasted bun.'),
    (6, 'Loaded Baked Potato', 8.99, 'Potato topped with cheese, bacon, and sour cream.'),
    (6, 'BBQ Wings', 9.99, 'Chicken wings glazed with BBQ sauce.'),
    (6, 'Grilled Salmon', 19.99, 'Fresh salmon fillet with lemon butter.'),
    (6, 'Coleslaw', 3.99, 'Classic creamy coleslaw.'),
    (6, 'Cornbread', 4.99, 'Sweet and moist cornbread.'),
    (6, 'Pulled Pork Sandwich', 11.99, 'Slow-cooked pulled pork with BBQ sauce.'),
    (6, 'Banana Pudding', 5.99, 'Classic Southern dessert with bananas and cream.'),

    (7, 'Margherita Pizza', 10.99, 'Classic pizza with tomato, mozzarella, and basil.'),
    (7, 'Pepperoni Pizza', 12.99, 'Mozzarella and pepperoni on tomato sauce.'),
    (7, 'BBQ Chicken Pizza', 14.99, 'Grilled chicken with BBQ sauce and red onions.'),
    (7, 'Veggie Supreme Pizza', 11.99, 'Bell peppers, onions, olives, and mushrooms.'),
    (7, 'Meat Lovers Pizza', 15.99, 'Pepperoni, sausage, bacon, and ham.'),
    (7, 'Hawaiian Pizza', 12.99, 'Pineapple, ham, and mozzarella cheese.'),
    (7, 'Garlic Bread', 4.99, 'Freshly baked garlic breadsticks.'),
    (7, 'Cheese Sticks', 6.99, 'Mozzarella sticks served with marinara sauce.'),
    (7, 'Buffalo Wings', 9.99, 'Spicy buffalo-style chicken wings.'),
    (7, 'Chocolate Lava Cake', 6.99, 'Warm chocolate cake with a gooey center.'),
    
    (8, 'Spaghetti Carbonara', 12.99, 'Creamy pasta with pancetta and Parmesan.'),
    (8, 'Lasagna', 14.99, 'Layers of pasta, meat, and cheese.'),
    (8, 'Margherita Pizza', 11.99, 'Tomato, mozzarella, and fresh basil.'),
    (8, 'Chicken Alfredo', 13.99, 'Grilled chicken over fettuccine with Alfredo sauce.'),
    (8, 'Bruschetta', 7.99, 'Toasted bread with tomatoes and basil.'),
    (8, 'Caesar Salad', 8.99, 'Crisp romaine with Caesar dressing and croutons.'),
    (8, 'Tiramisu', 6.99, 'Layered dessert with mascarpone and espresso.'),
    (8, 'Panna Cotta', 5.99, 'Creamy Italian dessert with berry sauce.'),
    (8, 'Garlic Bread', 3.99, 'Freshly baked garlic breadsticks.'),
    (8, 'Gelato', 4.99, 'Authentic Italian ice cream in various flavors.'),

    (9, 'California Roll', 8.99, 'Crab, avocado, and cucumber rolled with rice.'),
    (9, 'Spicy Tuna Roll', 9.99, 'Fresh tuna with spicy mayo.'),
    (9, 'Salmon Sashimi', 12.99, 'Fresh sliced salmon.'),
    (9, 'Dragon Roll', 13.99, 'Eel, avocado, and cucumber with eel sauce.'),
    (9, 'Rainbow Roll', 14.99, 'Variety of fish and avocado over a California roll.'),
    (9, 'Shrimp Tempura Roll', 10.99, 'Crispy shrimp tempura with avocado.'),
    (9, 'Dynamite Roll', 11.99, 'Spicy seafood mix with cucumber.'),
    (9, 'Miso Soup', 3.99, 'Traditional Japanese soup with tofu and seaweed.'),
    (9, 'Edamame', 4.99, 'Steamed soybeans with sea salt.'),
    (9, 'Mochi Ice Cream', 5.99, 'Ice cream wrapped in sweet rice dough.'),

    (10, 'Classic Cheeseburger', 8.99, 'Juicy beef patty with cheddar cheese.'),
    (10, 'Bacon Cheeseburger', 10.99, 'Beef patty with crispy bacon.'),
    (10, 'Mushroom Burger', 9.99, 'Beef patty with sautéed mushrooms and Swiss cheese.'),
    (10, 'Spicy Chicken Burger', 11.99, 'Crispy chicken with spicy mayo.'),
    (10, 'Veggie Burger', 8.99, 'Plant-based patty with lettuce and tomato.'),
    (10, 'Fries', 3.99, 'Crispy golden fries.'),
    (10, 'Onion Rings', 4.99, 'Battered and fried onion rings.'),
    (10, 'Milkshake', 4.99, 'Vanilla, chocolate, or strawberry shake.'),
    (10, 'BBQ Pulled Pork Sandwich', 9.99, 'Pulled pork with tangy BBQ sauce.'),
    (10, 'Apple Pie', 5.99, 'Warm apple pie with cinnamon.');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE Restaurant
DROP COLUMN street,
DROP COLUMN zip_code,
ADD COLUMN Address TEXT;

DROP TABLE Address;

DELETE FROM Menuitem WHERE restaurantid IN(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

-- +goose StatementEnd