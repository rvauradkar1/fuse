package fuse

func M1() string {
	return "testing"
}




type OrderController struct {
	ordSvc  IOrderService   `_fuse:"OrdSvc,value"`
}

func (ordCtrl *OrderController) Order(id string)  error {
	return nil
}

type IOrderService interface {
	findOrder() string
}

type OrderService struct {

}

type AuthService struct {

}

type OrderDB struct {

}

type CartService struct {

}
