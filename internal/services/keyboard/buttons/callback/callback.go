package callback

import (
	"strings"
)

type Callback struct {
	data map[string]string
}

func New() *Callback {
	return &Callback{
		data: make(map[string]string),
	}
}

func (c *Callback) Endpoint() string {
	return c.data["endpoint"]
}

func (c *Callback) SetEndpoint(endpoint string) *Callback {
	c.data["endpoint"] = endpoint

	return c
}

func (c *Callback) Value(key string) string {
	return c.data[key]
}

const (
	expSeparator  = "&"
	pairSeparator = "="
)

func FromString(data string) *Callback {
	rows := strings.Split(data, expSeparator)
	callback := New()

	for _, row := range rows {
		item := strings.Split(row, pairSeparator)
		if len(item) == 0 || item[0] == "" {
			continue
		}

		if len(item) == 1 {
			callback.data[item[0]] = ""

			continue
		}

		callback.data[item[0]] = item[1]
	}

	return callback
}

func (c *Callback) With(k, v string) *Callback {
	cp := c.Clone()

	cp.data[k] = v

	return cp
}

func (c *Callback) WithCallback(o *Callback) *Callback {
	c = c.Clone()

	for k, v := range o.data {
		c.data[k] = v
	}

	return c
}

func (c *Callback) Clone() *Callback {
	cp := New()
	cp.data = make(map[string]string, len(c.data))

	for k, v := range c.data {
		cp.data[k] = v
	}

	return cp
}

func (c *Callback) Data() string {
	if len(c.data) == 0 {
		return ""
	}

	i := 0

	sb := &strings.Builder{}

	for k, v := range c.data {
		if i != 0 {
			sb.WriteString(expSeparator)
		}

		i++

		sb.WriteString(k)

		if v != "" {
			sb.WriteString(pairSeparator + v)
		}
	}

	return sb.String()
}
