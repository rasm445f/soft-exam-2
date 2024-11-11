# soft-exam-2

## How to run the project

1. Clone the repository `gh repo clone rasm445f/soft-exam-2`
2. Run in your terminal `chmod +x install-tools.sh`
3. Run the install-tools.sh scripts like so `./install-tools.sh` if you don't already have the tools. (Ensure your PATH is set up correctly).
4. Rename the `.example.env` file to `.env` and populate the fields as needed.
5. Make sure you are in the ROOT of the project, Run `docker compose up` to start the PostgreSQL Docker container or `docker compose up -d` to run it in the background as a daemon on every boot.
6. Make a new terminal in the ROOT of the project
7. Run `make migrate-up` to setup the database with the tables etc. specified in the `db/migration/` folder.
8. Run `make run` to start the server.
9. Check the server is running by visiting `http://localhost:8080/` in your browser.
10. you can now test the endpoints using the swagger documentation at `http://localhost:8080/swagger/index.html`

## Technology Stack

### Version Control Platform:

- Git - Github

### Text Editing and Development Environment:

- VSCode / Neovim
- DBeaver
- Swagger

### General Online Research Tools:

- Stack Overflow
- MDN Web Docs
- Golang Docs

## Development Stack

### Backend Development:

- Golang

### Database Management:

- PostgreSQL

### Development Tools:

- Docker
- Docker Compose

### CI/CD Pipeline:

- GitHub Actions








-- Customer Table
CREATE TABLE Customer (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Email VARCHAR(255) UNIQUE NOT NULL,
    PhoneNumber VARCHAR(15),
    Address TEXT,
    OrderHistory TEXT,
    FeedbackHistory TEXT
);

-- Restaurant Table
CREATE TABLE Restaurant (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Address TEXT NOT NULL,
    Rating DECIMAL(2, 1)
);

-- DeliveryAgent Table
CREATE TABLE DeliveryAgent (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    ContactInfo TEXT,
    Availability BOOLEAN,
    DeliveryHistory TEXT
);

-- Order Table
CREATE TABLE "Order" (
    ID SERIAL PRIMARY KEY,
    TotalAmount DECIMAL(10, 2) NOT NULL,
    VATAmount DECIMAL(10, 2),
    Status VARCHAR(50) NOT NULL,
    Timestamp TIMESTAMP DEFAULT NOW(),
    Comment TEXT,
    CustomerID INT NOT NULL REFERENCES Customer(ID),
    RestaurantID INT NOT NULL REFERENCES Restaurant(ID),
    DeliveryAgentID INT REFERENCES DeliveryAgent(ID),
    PaymentID INT REFERENCES Payment(ID),
    BonusID INT REFERENCES Bonus(ID),
    FeedbackID INT REFERENCES Feedback(ID)
);

-- Payment Table
CREATE TABLE Payment (
    ID SERIAL PRIMARY KEY,
    OrderID INT NOT NULL REFERENCES "Order"(ID),
    PaymentStatus VARCHAR(50) NOT NULL,
    PaymentMethod VARCHAR(50) NOT NULL
);

-- Bonus Table
CREATE TABLE Bonus (
    ID SERIAL PRIMARY KEY,
    Description TEXT,
    EarlyLateAmount DECIMAL(10, 2),
    Percentage DECIMAL(5, 2)
);

-- Feedback Table
CREATE TABLE Feedback (
    ID SERIAL PRIMARY KEY,
    OrderID INT NOT NULL REFERENCES "Order"(ID),
    CustomerID INT NOT NULL REFERENCES Customer(ID),
    FoodRating INT CHECK (FoodRating BETWEEN 1 AND 5),
    DeliveryAgentRating INT CHECK (DeliveryAgentRating BETWEEN 1 AND 5),
    Comment TEXT
);

-- Complaint Table
CREATE TABLE Complaint (
    ID SERIAL PRIMARY KEY,
    CustomerID INT NOT NULL REFERENCES Customer(ID),
    RestaurantID INT NOT NULL REFERENCES Restaurant(ID),
    ComplaintType VARCHAR(50),
    Description TEXT,
    Status VARCHAR(50)
);

-- OrderItem Table
CREATE TABLE OrderItem (
    ID SERIAL PRIMARY KEY,
    OrderID INT NOT NULL REFERENCES "Order"(ID),
    Name VARCHAR(255) NOT NULL,
    Price DECIMAL(10, 2) NOT NULL,
    Amount DECIMAL(10, 2) NOT NULL
);

-- Fee Table
CREATE TABLE Fee (
    ID SERIAL PRIMARY KEY,
    OrderID INT NOT NULL REFERENCES "Order"(ID),
    Amount DECIMAL(10, 2) NOT NULL,
    Description TEXT
);
