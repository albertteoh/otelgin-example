package data

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	ID    string  `json:"id"`
}

type Inventory struct {
	Products []Product `json:"products"`
}

type Response struct {
	Inventory Inventory `json:"inventory"`
}
