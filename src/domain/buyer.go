package domain

type Buyer struct {
	ID   string `json:"id,omitempty"`
	UID  string `json:"uid,omitempty"`
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}
