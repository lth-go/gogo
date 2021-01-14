package vm

//
// ConstantPool
//
type ConstantPool struct {
	pool []interface{}
}

func NewConstantPool() ConstantPool {
	return ConstantPool{
		pool: []interface{}{},
	}
}

func (c *ConstantPool) SetPool(pool []interface{}) {
	c.pool = pool
}

func (c *ConstantPool) GetInt(index int) int {
	return c.pool[index].(int)
}

func (c *ConstantPool) GetFloat(index int) float64 {
	return c.pool[index].(float64)
}

func (c *ConstantPool) GetString(index int) string {
	return c.pool[index].(string)
}
