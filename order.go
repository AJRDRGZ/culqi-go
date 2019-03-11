package culqi

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	orderURL = baseURL + "/orders"
)

// Order objeto request de Order
type Order struct {
	Amount         int               `json:"amount"`
	CurrencyCode   string            `json:"currency_code"`
	Description    string            `json:"description"`
	OrderNumber    string            `json:"order_number"`
	ExpirationDate int               `json:"expiration_date"`
	ClientDetails  Customer          `json:"client_details"`
	Confirm        bool              `json:"confirm"`
	Metadata       map[string]string `json:"metadata"`
}

// ResponseOrder objeto respuesta de Order
type ResponseOrder struct {
	Order
	Object       string `json:"object"`
	ID           string `json:"id"`
	PaymentCode  string `json:"payment_code"`
	State        string `json:"state"`
	TotalFee     int    `json:"total_fee"`
	NetAmount    int    `json:"net_amount"`
	UpdatedAt    int    `json:"updated_at"`
	PaidAt       int    `json:"paid_at"`
	CreationDate int    `json:"creation_date"`
}

// ResponseOrderAll respuesta de tarjeta para GetAll
type ResponseOrderAll struct {
	Data []ResponseOrder `json:"data"`
	WrapperResponse
}

// Create método para crear un Order
func (ord *Order) Create() (*ResponseOrder, error) {
	j, err := json.Marshal(ord)
	if err != nil {
		return nil, err
	}

	res, err := do("POST", orderURL, nil, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	fmt.Println("Create:", string(res))
	rord := &ResponseOrder{}
	if err = json.Unmarshal(res, rord); err != nil {
		return nil, err
	}

	return rord, nil
}

// GetByID método para obtener un Order por id
func (ord *Order) GetByID(id string) (*ResponseOrder, error) {
	if id == "" {
		return nil, ErrParameter
	}

	res, err := do("GET", orderURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("GetByID:", string(res))
	rord := &ResponseOrder{}
	if err = json.Unmarshal(res, rord); err != nil {
		return nil, err
	}

	return rord, nil
}
