package culqi_test

import (
	"strings"
	"testing"
	"time"

	culqi "github.com/AJRDRGZ/culqi-go"
)

func TestOrder_Create(t *testing.T) {
	if secretKey == "" {
		t.Skip("No se indicó una llave privada")
	}

	culqi.Key(publicKey, secretKey)
	o := culqi.Order{
		Amount:         3000, // Monto del Order a cobrar recurrentemente. Sin punto decimal Ejemplo: 30.00 serían 3000
		CurrencyCode:   "USD",
		Description:    "Orden de compra de cursos",
		OrderNumber:    "my_internal_orderID3",
		ExpirationDate: int(time.Now().AddDate(0, 0, 2).Unix()),
		ClientDetails: culqi.Customer{
			FirstName:   "Alejandro",
			LastName:    "Rodriguez",
			PhoneNumber: "7777777",
			Email:       "aj.softwaredeveloper@gmail.com",
		},
		Confirm: false,
		Metadata: map[string]string{
			"order_type": "1",
		},
	}

	res, err := o.Create()
	if err != nil {
		t.Fatalf("Order.Create() err = %v; want = %v", err, nil)
	}

	if res == nil {
		t.Fatalf("ResponseOrder = nil; want non-nil value")
	}

	if res.State != "created" {
		t.Errorf("ResponseOrder.State = %s; want = %q", res.State, "created")
	}

	if !strings.HasPrefix(res.ID, "ord_test_") {
		t.Errorf("ResponseOrder.ID = %s; want prefix = %q", res.ID, "ord_test_")
	}

}

func TestOrder_GetByID(t *testing.T) {
	if secretKey == "" {
		t.Skip("No se indicó una llave privada")
	}

	culqi.Key(publicKey, secretKey)

	p := culqi.Order{}
	res, err := p.GetByID("ord_test_R6GwSS8xvmnjsgjo")
	if err != nil {
		t.Fatalf("Order.GetByID() err = %v; want = %v", err, nil)
	}

	if res == nil {
		t.Fatalf("ResponseOrder = nil; want non-nil value")
	}
}
