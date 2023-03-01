package model

type Product struct {
	ID          int64  `gorm:"index" json:"id"`
	ProductName string `gorm:"column:productname;not null" json:"productname"`
	Count       int64  `gorm:"column:count;not null" json:"count"`
	IsDelete 	int		`gorm:"column:isdelete;not null" json:"isdelete"`
}

func (u *Product) TableName() string {
	return "product"
}

type ProductInfo struct {
	Product string `json:"product"`
	Status  int    `json:"status"`
}


type Order struct {
	ID          int64  `gorm:"index" json:"id"`
	ProductName string `gorm:"column:productname;not null" json:"productname"`
	Count       int64  `gorm:"column:count;not null" json:"count"`

}

func (u *Order) TableName() string {
	return "order"
}