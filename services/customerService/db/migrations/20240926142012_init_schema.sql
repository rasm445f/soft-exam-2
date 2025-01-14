-- +goose Up
-- +goose StatementBegin

-- Create ZipCode table first
CREATE TABLE IF NOT EXISTS ZipCode (
    Zip_Code INT PRIMARY KEY,
    City varchar(255),
    CONSTRAINT zipcode_not_null CHECK (Zip_Code IS NOT NULL)
);

-- Create Address table second
CREATE TABLE IF NOT EXISTS Address (
    ID serial PRIMARY KEY,
    Street_Address text,
    Zip_Code int REFERENCES ZipCode(Zip_Code),
    CONSTRAINT address_not_null CHECK (Address IS NOT NULL)
);

-- Create Customer table last
CREATE TABLE IF NOT EXISTS Customer (
    ID serial PRIMARY KEY,
    Name varchar(255),
    Email varchar(255) UNIQUE,
    PASSWORD varchar(255),
    PhoneNumber varchar(15),
    AddressID int REFERENCES Address(ID),
    CONSTRAINT email_not_null CHECK (Email IS NOT NULL),
    CONSTRAINT password_not_null CHECK (Password IS NOT NULL)
);

-- Insert all Danish zip codes
INSERT INTO ZipCode (Zip_Code, City) VALUES
(1000, 'København K'),
(1050, 'København K'),
(1100, 'København K'),
(1150, 'København K'),
(1200, 'København K'),
(1250, 'København K'),
(1300, 'København K'),
(1350, 'København K'),
(1400, 'København K'),
(1450, 'København K'),
(1500, 'København V'),
(1550, 'København V'),
(1600, 'København V'),
(1650, 'København V'),
(1700, 'København V'),
(1750, 'København V'),
(1800, 'Frederiksberg C'),
(1850, 'Frederiksberg C'),
(1900, 'Frederiksberg C'),
(1950, 'Frederiksberg C'),
(2000, 'Frederiksberg'),
(2100, 'København Ø'),
(2150, 'Nordhavn'),
(2200, 'København N'),
(2300, 'København S'),
(2400, 'København NV'),
(2450, 'København SV'),
(2500, 'Valby'),
(2600, 'Glostrup'),
(2605, 'Brøndby'),
(2610, 'Rødovre'),
(2620, 'Albertslund'),
(2625, 'Vallensbæk'),
(2630, 'Taastrup'),
(2635, 'Ishøj'),
(2640, 'Hedehusene'),
(2650, 'Hvidovre'),
(2660, 'Brøndby Strand'),
(2665, 'Vallensbæk Strand'),
(2670, 'Greve'),
(2680, 'Solrød Strand'),
(2690, 'Karlslunde'),
(2700, 'Brønshøj'),
(2720, 'Vanløse'),
(2730, 'Herlev'),
(2740, 'Skovlunde'),
(2750, 'Ballerup'),
(2760, 'Måløv'),
(2765, 'Smørum'),
(2770, 'Kastrup'),
(2791, 'Dragør'),
(2800, 'Kongens Lyngby'),
(2820, 'Gentofte'),
(2830, 'Virum'),
(2840, 'Holte'),
(2850, 'Nærum'),
(2860, 'Søborg'),
(2870, 'Dyssegård'),
(2880, 'Bagsværd'),
(2900, 'Hellerup'),
(2920, 'Charlottenlund'),
(2930, 'Klampenborg'),
(2942, 'Skodsborg'),
(2950, 'Vedbæk'),
(2960, 'Rungsted Kyst'),
(2970, 'Hørsholm'),
(2980, 'Kokkedal'),
(2990, 'Nivå'),
(3000, 'Helsingør'),
(3050, 'Humlebæk'),
(3060, 'Espergærde'),
(3070, 'Snekkersten'),
(3080, 'Tikøb'),
(3100, 'Hornbæk'),
(3120, 'Dronningmølle'),
(3140, 'Ålsgårde'),
(3150, 'Hellebæk'),
(3200, 'Helsinge'),
(3210, 'Vejby'),
(3220, 'Tisvildeleje'),
(3230, 'Græsted'),
(3250, 'Gilleleje'),
(3300, 'Frederiksværk'),
(3310, 'Ølsted'),
(3320, 'Skævinge'),
(3330, 'Gørløse'),
(3360, 'Liseleje'),
(3370, 'Melby'),
(3390, 'Hundested'),
(3400, 'Hillerød'),
(3450, 'Allerød'),
(3460, 'Birkerød'),
(3480, 'Fredensborg'),
(3490, 'Kvistgård'),
(3500, 'Værløse'),
(3520, 'Farum'),
(3540, 'Lynge'),
(3550, 'Slangerup'),
(3600, 'Frederikssund'),
(3630, 'Jægerspris'),
(3650, 'Ølstykke'),
(3660, 'Stenløse'),
(3670, 'Veksø Sjælland'),
(3700, 'Rønne'),
(3720, 'Aakirkeby'),
(3730, 'Nexø'),
(3740, 'Svaneke'),
(3751, 'Østermarie'),
(3760, 'Gudhjem'),
(3770, 'Allinge'),
(3782, 'Klemensker'),
(3790, 'Hasle'),
(4000, 'Roskilde'),
(4040, 'Jyllinge'),
(4050, 'Skibby'),
(4060, 'Kirke Såby'),
(4070, 'Kirke Hyllinge'),
(4100, 'Ringsted'),
(4130, 'Viby Sjælland'),
(4140, 'Borup'),
(4160, 'Herlufmagle'),
(4171, 'Glumsø'),
(4173, 'Fjenneslev'),
(4174, 'Jystrup Midtsjælland'),
(4180, 'Sorø'),
(4190, 'Munke Bjergby'),
(4200, 'Slagelse'),
(4220, 'Korsør'),
(4230, 'Skælskør'),
(4241, 'Vemmelev'),
(4242, 'Boeslunde'),
(4243, 'Rude'),
(4250, 'Fuglebjerg'),
(4261, 'Dalmose'),
(4262, 'Sandved'),
(4270, 'Høng'),
(4281, 'Gørlev'),
(4291, 'Ruds Vedby'),
(4293, 'Dianalund'),
(4295, 'Stenlille'),
(4296, 'Nyrup'),
(4300, 'Holbæk'),
(4320, 'Lejre'),
(4330, 'Hvalsø'),
(4340, 'Tølløse'),
(4350, 'Ugerløse'),
(4360, 'Kirke Eskilstrup'),
(4370, 'Store Merløse'),
(4390, 'Vipperød'),
(4400, 'Kalundborg'),
(9000, 'Aalborg')
ON CONFLICT (Zip_Code) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop tables in reverse order
DROP TABLE IF EXISTS Customer;
DROP TABLE IF EXISTS Address;
DROP TABLE IF EXISTS ZipCode;

-- +goose StatementEnd
