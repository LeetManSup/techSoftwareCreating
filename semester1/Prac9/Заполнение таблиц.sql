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
