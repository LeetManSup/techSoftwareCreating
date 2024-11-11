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
