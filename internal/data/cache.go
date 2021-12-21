package data

import "sync"

type Cache struct {
	sync.RWMutex
	orderSet map[string]ReceivedOrder
}

func NewCache() *Cache {
	return &Cache{
        orderSet: make(map[string]ReceivedOrder),
    }
}

func (c *Cache) AddOrder(order ReceivedOrder) {
	c.Lock()

	c.orderSet[*order.OrderUID] = order

	c.Unlock()
}

func (c *Cache) GetOrder(orderUID string) (ReceivedOrder, bool) {
	// Если метод нуждается только в чтении — он использует RLock(),
	// который не заблокирует другие операции чтения,
	// но заблокирует операцию записи и наоборот. 
	c.RLock()
	defer c.RUnlock()

	order, ok := c.orderSet[orderUID]
	return order, ok
}

func (c *Cache) GetOrders() []ReceivedOrder {
	res := make([]ReceivedOrder, 0, len(c.orderSet))

	c.RLock()
	defer c.RUnlock()

	for _, ro := range c.orderSet {
		res = append(res, ro)
	}

	return res
}