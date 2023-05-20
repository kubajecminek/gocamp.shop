package models

type Cart []CartItem

type CartItem struct {
	Item     Item `bigquery:"item"`
	Quantity int  `bigquery:"quantity"`
}

func NewCart() Cart {
	return make(Cart, 0)
}

func (c Cart) IsEmpty() bool {
	return len(c) == 0
}

func (c Cart) NumCamps() int {
	var num int
	for _, item := range c {
		if item.Item.IsCamp() {
			num += item.Quantity
		}
	}
	return num
}

func (c *Cart) Add(item Item, quantity int) {
	var newCart Cart
	cart := make(map[Item]int)
	for _, item := range *c {
		cart[item.Item] = item.Quantity
	}
	cart[item] += quantity
	if cart[item] <= 0 {
		delete(cart, item)
	}
	for item, quantity := range cart {
		newCart = append(newCart, CartItem{Item: item, Quantity: quantity})
	}
	*c = newCart
}

func (c Cart) CampsIterable() map[Item][]struct{} {
	items := make(map[Item][]struct{})
	for _, item := range c {
		if item.Item.IsCamp() {
			items[item.Item] = make([]struct{}, item.Quantity)
		}
	}
	return items
}

func (c Cart) CampInside() bool {
	for _, item := range c {
		if item.Item.IsCamp() {
			return true
		}
	}
	return false
}

func (c Cart) BackprintInside() bool {
	for _, item := range c {
		if item.Item.IsBackprint() {
			return true
		}
	}
	return false
}
