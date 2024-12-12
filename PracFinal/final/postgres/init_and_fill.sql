-- Таблица Роль в системе
CREATE TABLE UserRole (
    RoleID INT PRIMARY KEY NOT NULL,
    RoleName VARCHAR(20) NOT NULL
);

-- Таблица Пользовательская учетная запись
CREATE TABLE Account (
    AccountLogin VARCHAR(30) PRIMARY KEY NOT NULL,
    AccountPassword VARCHAR(30) NOT NULL,
    AccountRoleID INT NOT NULL,
    FOREIGN KEY (AccountRoleID) REFERENCES UserRole(RoleID)
);

-- Таблица Корзина
CREATE TABLE Cart (
    CartID INT PRIMARY KEY NOT NULL,
    CartDateOfLast TIMESTAMPTZ NOT NULL
);

-- Таблица Покупатель
CREATE TABLE Customer (
    CustomerID INT PRIMARY KEY NOT NULL,
    CustomerName VARCHAR(30) NOT NULL,
    CustomerSurname VARCHAR(30) NOT NULL,
    CustomerPatronymic VARCHAR(30) NOT NULL,
    CustomerPhone VARCHAR(12) NOT NULL,
    CustomerEmail VARCHAR(254),
    CustomerCart INT NOT NULL,
    CustomerAccountLogin VARCHAR(30) NOT NULL,
    FOREIGN KEY (CustomerCart) REFERENCES Cart(CartID),
    FOREIGN KEY (CustomerAccountLogin) REFERENCES Account(AccountLogin)
);

-- Таблица Адрес
CREATE TABLE Address (
    AddressID INT PRIMARY KEY NOT NULL,
    AddressCity VARCHAR(30) NOT NULL,
    AddressStreet VARCHAR(50) NOT NULL,
    AddressHouse SMALLINT NOT NULL,
    AddressCorps SMALLINT,
    AddressBuilding SMALLINT,
    AddressCustomerID INT NOT NULL,
    FOREIGN KEY (AddressCustomerID) REFERENCES Customer(CustomerID)
);

-- Таблица Производитель
CREATE TABLE Manufacturer (
    ManufacturerID INT PRIMARY KEY NOT NULL,
    ManufacturerName VARCHAR(150) NOT NULL,
    ManufacturerCountry VARCHAR(58) NOT NULL
);

-- Таблица Товар
CREATE TABLE Product (
    ProductID INT PRIMARY KEY NOT NULL,
    ProductName VARCHAR(30) NOT NULL,
    ProductDescription TEXT NOT NULL,
    ProductPrice MONEY NOT NULL,
    ProductCount INT NOT NULL,
    ProductManufacturerID INT NOT NULL,
    FOREIGN KEY (ProductManufacturerID) REFERENCES Manufacturer(ManufacturerID)
);

-- Таблица Поставщик
CREATE TABLE Supplier (
    SupplierID INT PRIMARY KEY NOT NULL,
    SupplierCompany VARCHAR(254) NOT NULL,
    SupplierPhone VARCHAR(12) NOT NULL,
    SupplierBankAccount VARCHAR(9) NOT NULL,
    SupplierBIC INT NOT NULL,
    SupplierNN INT8 NOT NULL,
    SupplierKPP INT NOT NULL
);

-- Таблица Статус заказа
CREATE TABLE Status (
    StatusID INT PRIMARY KEY NOT NULL,
    StatusName VARCHAR(15) NOT NULL
);

-- Таблица Поставщик товара
CREATE TABLE ProductSupplier (
    ProductID INT NOT NULL,
    SupplierID INT NOT NULL,
    PRIMARY KEY (ProductID, SupplierID),
    FOREIGN KEY (ProductID) REFERENCES Product(ProductID),
    FOREIGN KEY (SupplierID) REFERENCES Supplier(SupplierID)
);

-- Таблица Заказ
CREATE TABLE CustomerOrder (
    OrderID INT PRIMARY KEY,
    OrderTime TIMESTAMPTZ NOT NULL,
    OrderAddressID INT NOT NULL,
    OrderStatusID INT NOT NULL,
    OrderSeller VARCHAR(30) NOT NULL,
    OrderCourier VARCHAR(30) NOT NULL,
    AccountLogin VARCHAR(30) NOT NULL,
    FOREIGN KEY (OrderAddressID) REFERENCES Address(AddressID),
    FOREIGN KEY (OrderStatusID) REFERENCES Status(StatusID),
    FOREIGN KEY (OrderSeller) REFERENCES Account(AccountLogin),
    FOREIGN KEY (OrderCourier) REFERENCES Account(AccountLogin),
    FOREIGN KEY (AccountLogin) REFERENCES Account(AccountLogin)
);

-- Таблица Товар из корзины
CREATE TABLE CartProduct (
    CartProductProductID INT NOT NULL,
    CartProductCartID INT NOT NULL,
    PRIMARY KEY (CartProductProductID, CartProductCartID),
    FOREIGN KEY (CartProductProductID) REFERENCES Product(ProductID),
    FOREIGN KEY (CartProductCartID) REFERENCES Cart(CartID)
);

-- Таблица Продукт в заказе
CREATE TABLE ProductOrder (
    ProductID INT NOT NULL,
    OrderID INT NOT NULL,
    PRIMARY KEY (ProductID, OrderID),
    FOREIGN KEY (ProductID) REFERENCES Product(ProductID),
    FOREIGN KEY (OrderID) REFERENCES CustomerOrder(OrderID)
);

-- Заполнение таблиц
-- Вставка данных в таблицу Роль в системе
INSERT INTO UserRole (RoleID, RoleName) VALUES 
    (1, 'Admin'),
    (2, 'Seller'),
    (3, 'Courier');

-- Вставка данных в таблицу Пользовательская учетная запись
INSERT INTO Account (AccountLogin, AccountPassword, AccountRoleID) VALUES 
    ('admin1', 'password123', 1),
    ('seller1', 'password456', 2),
    ('courier1', 'password789', 3);

-- Вставка данных в таблицу Корзина
INSERT INTO Cart (CartID, CartDateOfLast) VALUES 
    (1, '2024-11-11 10:30:00'),
    (2, '2024-11-12 12:00:00');

-- Вставка данных в таблицу Покупатель
INSERT INTO Customer (CustomerID, CustomerName, CustomerSurname, CustomerPatronymic, CustomerPhone, CustomerEmail, CustomerCart, CustomerAccountLogin) VALUES 
    (1, 'Иван', 'Иванов', 'Иванович', '79991234567', 'ivanov@example.com', 1, 'admin1'),
    (2, 'Петр', 'Петров', 'Петрович', '79991234568', 'petrov@example.com', 2, 'seller1');

-- Вставка данных в таблицу Адрес
INSERT INTO Address (AddressID, AddressCity, AddressStreet, AddressHouse, AddressCorps, AddressBuilding, AddressCustomerID) VALUES 
    (1, 'Москва', 'Ленина', 10, NULL, 1, 1),
    (2, 'Санкт-Петербург', 'Невский пр.', 15, 2, NULL, 2);

-- Вставка данных в таблицу Производитель
INSERT INTO Manufacturer (ManufacturerID, ManufacturerName, ManufacturerCountry) VALUES 
    (1, 'Samsung', 'Южная Корея'),
    (2, 'Apple', 'США');

-- Вставка данных в таблицу Товар
INSERT INTO Product (ProductID, ProductName, ProductDescription, ProductPrice, ProductCount, ProductManufacturerID) VALUES 
    (1, 'Телефон', 'Смартфон с 64 ГБ памяти', 30000, 50, 1),
    (2, 'Ноутбук', 'Ноутбук с 256 ГБ SSD', 70000, 30, 2);

-- Вставка данных в таблицу Поставщик
INSERT INTO Supplier (SupplierID, SupplierCompany, SupplierPhone, SupplierBankAccount, SupplierBIC, SupplierNN, SupplierKPP) VALUES 
    (1, 'TechSupply', '79990001234', '123456789', 123456, 123456789, 123456789),
    (2, 'GadgetWholesale', '79990004321', '987654321', 654321, 987654321, 987654321);

-- Вставка данных в таблицу Статус заказа
INSERT INTO Status (StatusID, StatusName) VALUES 
    (1, 'В обработке'),
    (2, 'Доставлен'),
    (3, 'Отменен');

-- Вставка данных в таблицу Поставщик товара
INSERT INTO ProductSupplier (ProductID, SupplierID) VALUES 
    (1, 1),
    (2, 2);

-- Вставка данных в таблицу Заказ
INSERT INTO CustomerOrder (OrderID, OrderTime, OrderAddressID, OrderStatusID, OrderSeller, OrderCourier, AccountLogin) VALUES 
    (1, '2024-11-11 14:00:00', 1, 1, 'seller1', 'courier1', 'admin1'),
    (2, '2024-11-12 15:00:00', 2, 2, 'seller1', 'courier1', 'admin1');

-- Вставка данных в таблицу Товар из корзины
INSERT INTO CartProduct (CartProductProductID, CartProductCartID) VALUES 
    (1, 1),
    (2, 2);

-- Вставка данных в таблицу Продукт в заказе
INSERT INTO ProductOrder (ProductID, OrderID) VALUES 
    (1, 1),
    (2, 2);
