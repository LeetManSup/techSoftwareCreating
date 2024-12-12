package models

type Product struct {
	ProductID             int    `gorm:"column:productid;primaryKey" json:"product_id"`
	ProductName           string `gorm:"column:productname;type:varchar(30);not null" json:"product_name"`
	ProductDescription    string `gorm:"column:productdescription;type:text;not null" json:"product_description"`
	ProductPrice          string `gorm:"column:productprice;not null" json:"product_price"`
	ProductCount          int    `gorm:"column:productcount;not null" json:"product_count"`
	ProductManufacturerID int    `gorm:"column:productmanufacturerid;not null" json:"product_manufacturer_id"`
}

func (Product) TableName() string {
	return "product"
}

