package fuse

func M1() string {
	return "testing"
}

type OrderController struct {
	s       string
	OrdPtr  *OrderService `_fuse:"OrdSvc,ptr"`
	OrdSvc  IOrderService `_fuse:"OrdSvc,val"`
	OrdSvc2 OrderService  `_fuse:"OrdSvc,val"`
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

func (o OrderService) findOrder() string {
	return "order"
}

type AuthService struct {
}

type OrderDB struct {
}

type CartService struct {
}
