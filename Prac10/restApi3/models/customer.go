package models

type Customer struct {
	CustomerID           int    `gorm:"column:customerid;primaryKey" json:"customer_id"`
	CustomerName         string `gorm:"column:customername;type:varchar(30);not null" json:"customer_name"`
	CustomerSurname      string `gorm:"column:customersurname;type:varchar(30);not null" json:"customer_surname"`
	CustomerPatronymic   string `gorm:"column:customerpatronymic;type:varchar(30);not null" json:"customer_patronymic"`
	CustomerPhone        string `gorm:"column:customerphone;type:varchar(12);not null" json:"customer_phone"`
	CustomerEmail        string `gorm:"column:customeremail;type:varchar(254)" json:"customer_email"`
	CustomerCart         int    `gorm:"column:customercart;not null" json:"customer_cart"`
	CustomerAccountLogin string `gorm:"column:customeraccountlogin;type:varchar(30);not null" json:"customer_account_login"`
}

func (Customer) TableName() string {
	return "customer"
}
