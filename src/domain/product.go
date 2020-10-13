package domain

type Product struct {
	ID    string `json:"id,omitempty"`
	UID   string `json:"uid,omitempty"`
	Name  string `json:"name,omitempty"`
	Price string `json:"price,omitempty"`
}
