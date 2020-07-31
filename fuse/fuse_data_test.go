package fuse

import "fmt"

func M1() string {
	return "testing"
}

type OrderController struct {
	s       string
	OrdPtr  *OrderService `_fuse:"OrdSvc"`
	OrdSvc  IOrderService `_fuse:"OrdSvc"`
	OrdSvc2 OrderService  `_fuse:"OrdSvc"`
}

func (ordCtrl *OrderController) Order(id string) error {
	return nil
}

type OrderController1 struct {
	s       string `_fuse:"OrdSvc"`
	t       string
	aOrdPtr *OrderService `_fuse:"OrdSvc"`
	OrdPtr  *OrderService `_fuse:"OrdSvc"`
	OrdSvc  IOrderService `_fuse:"OrdSvc"`
}

func (ordCtrl *OrderController1) Order(id string) error {
	return nil
}

type OrderController2 struct {
	s      string
	OrdPtr *OrderService `_fuse:"OrdCtrl"`
	OrdSvc IOrderService `_fuse:"OrdCtrl"`
}

func (ordCtrl *OrderController2) Order(id string) error {
	return nil
}

type IOrderService interface {
	findOrder() string
	saveOrder() string
}

type OrderService struct {
	T      string
	Status string
}

func (o OrderService) findOrder() string {
	return o.T
}

func (o OrderService) saveOrder() string {
	o.Status = "Saved"
	return "saved"
}

type IAuthService interface {
	auth() bool
}

type AuthService struct {
	Status string
}

func (a *AuthService) auth() bool {
	a.Status = "authenticated"
	return true
}

type OrderDB struct {
}

type CartService struct {
}

type Isvc1 interface {
	M1()
}

type Svc1 struct {
	S2 Isvc2 `_fuse:"svc2"`
	S3 *Svc3 `_fuse:"svc3"`
	s  string
}

func (i Svc1) M1() {
	fmt.Println("Inside svc1 M1")
}

type Isvc2 interface {
	M2()
}

type Svc2 struct {
	s string
}

func (i Svc2) M2() {
	fmt.Println("Inside svc2 M2")
}

type Svc3 struct {
	s string
}

func M3() {
	fmt.Println("Inside svc3 M3")
}
