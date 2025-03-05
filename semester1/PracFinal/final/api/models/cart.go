package models

type Cart struct {
	CartID         int    `gorm:"column:cartid;primaryKey" json:"cart_id"`
	CartDateOfLast string `gorm:"column:cartdateoflast;type:timestamp;not null" json:"cart_date_of_last"`
}

func (Cart) TableName() string {
	return "cart"
}

