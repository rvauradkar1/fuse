package fuse

func M1() string {
	return "testing"
}

type OrderController struct {
	s      string
	OrdSvc IOrderService `_fuse:"OrdSvc,value"`
}

func (ordCtrl *OrderController) Order(id string) error {
	return nil
}

type IOrderService interface {
	findOrder() string
}

type OrderService struct {
	t string
}

func (o *OrderService) findOrder() string {
	return "order"
}

type AuthService struct {
}

type OrderDB struct {
}

type CartService struct {
}
