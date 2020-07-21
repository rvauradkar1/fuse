package fuse

func M1() string {
	return "testing"
}



type OrderController struct {

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
