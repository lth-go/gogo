package vm

type ConstantPool struct {
	pool []Constant
}

func NewConstantPool() ConstantPool {
	return ConstantPool{
		pool: []Constant{},
	}
}

func (c *ConstantPool) SetPool(pool []Constant) {
	c.pool = pool
}

func (c *ConstantPool) Append(value Constant) {
	c.pool = append(c.pool, value)
}

func (c *ConstantPool) Length() int {
	return len(c.pool)
}

func (c *ConstantPool) getInt(index int) int {
	return c.pool[index].getInt()
}

func (c *ConstantPool) getDouble(index int) float64 {
	return c.pool[index].getDouble()
}

func (c *ConstantPool) getString(index int) string {
	return c.pool[index].getString()
}
