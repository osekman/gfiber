package model

type Customers struct {
	Id               int    `orm:"id" json:"id"`
	CustomerId       int    `orm:"CustomerId" json:"CustomerId"`
	CustomerCode     string `orm:"CustomerCode" json:"CustomerCode"`
	CustomerUsername string `orm:"CustomerUsername" json:"CustomerUsername"`
	CustomerGroupId  int    `orm:"CustomerGroupId" json:"CustomerGroupId"`
	CustomerName     string `orm:"CustomerName" json:"CustomerName"`
	CustomerPhone    string `orm:"CustomerPhone" json:"CustomerPhone"`
	InvoiceMobile    string `orm:"InvoiceMobile" json:"InvoiceMobile"`
	DeliveryMobile   string `orm:"DeliveryMobile" json:"DeliveryMobile"`
}

func (*Customers) customers() string {
	return "customers"
}
