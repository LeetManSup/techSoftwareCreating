package models

type CustomerOrder struct {
	OrderID        int    `gorm:"column:orderid;primaryKey" json:"order_id"`
	OrderTime      string `gorm:"column:ordertime;type:timestamp;not null" json:"order_time"`
	OrderAddressID int    `gorm:"column:orderaddressid;not null" json:"order_address_id"`
	OrderStatusID  int    `gorm:"column:orderstatusid;not null" json:"order_status_id"`
	OrderSeller    string `gorm:"column:orderseller;type:varchar(30);not null" json:"order_seller"`
	OrderCourier   string `gorm:"column:ordercourier;type:varchar(30);not null" json:"order_courier"`
	AccountLogin   string `gorm:"column:accountlogin;type:varchar(30);not null" json:"account_login"`
}

func (CustomerOrder) TableName() string {
	return "customerorder"
}

