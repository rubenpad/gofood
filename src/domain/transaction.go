package domain

type Transaction struct {
	ID       string `json:"id,omitempty"`
	Device   string `json:"device,omitempty"`
	Products []Uid  `json:"products_id,omitempty"`
	From     Ip     `json:"from,omitempty"`
	Owner    Uid    `json:"owner,omitempty"`
}
