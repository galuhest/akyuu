package akyuu

import (
	"fmt"
	"reflect"
	"encoding/json"
	"github.com/pkg/errors"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bukalapak/dallimin"
)

// Task holds type alias for empty interface (interface{})
// We use type alias instead of simply use string for clean code purpose.
type Task interface{}

// Connection is a struct that holds the information
// about the memcached used by user
type Connection struct {
	mem *memcache.Client
}

// Value struct will hold the result of Task
type Value struct {
	Value interface{} `json:"value"`
}

// Create new Connection
func New(dalli *dallimin.Ring) (*Connection, error) {
	return &Connection {
		mem: memcache.NewFromSelector(dalli),
	}, nil
}

// Fetch has 2 behavior:
// if the given key already exists in memcached, it returns the oject inside the memcached
// if it isn't, insert the result of task into memcached as value and key as key
func (c Connection) Fetch(task Task, key string, expiration int32) (*memcache.Item, error) {
	item, err := c.mem.Get(key)
	// key exists
	if err == nil {
		return item, errors.New("Key already exists")
	}

	res, err := runTask(task)
	// task returns an error
	if err != nil {
		return nil, errors.Wrap(err, "Task failed to complete")
	}

	val := &Value{Value: res}

	jval, err := json.Marshal(val)

	if err != nil {
		return nil, errors.Wrap(err, "Can't encapsulate return value")
	}

	mitem := &memcache.Item{
		Key: key,
		Value: []byte(jval),
		Expiration: expiration,
	}
	err = c.mem.Set(mitem)
	if err != nil {
		return nil, errors.Wrap(err, "Can't insert data to memcache")
	}

	return item, nil
}

// runTask will run Task while inferring it's signature
// ret holds many value for the sake of compatibility with any function
func runTask(intf interface{}, ret ...interface{}) (interface{}, error) {
	// too many return value from function
	if len(ret) > 1 {
		return nil, errors.New("Too many return value (max 2)")
	}

	fn := reflect.ValueOf(intf)
	fnType := fn.Type()

	// if intf is not a function but a value
	if fnType.Kind() != reflect.Func {
		err := ret[0]
		// TODO: check if err is an error !
		if err != nil {
			errVal := reflect.ValueOf(err).Interface()
			return nil, errors.New(fmt.Sprintf("%s",errVal))
		}

		return fn, nil
	}

	// intf is a function, execute intf
	var res []reflect.Value
	retAmt := fnType.NumOut()
	// return value is more than 2
	if retAmt > 2 {
		return nil, errors.New("Too many return value (max 2)")
	} 

	res = fn.Call([]reflect.Value{})
	// return value has error value
	if retAmt == 2 {
		err := res[1].Interface()
		// TODO : check if err is an error
		// error is not nil
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s",err))
		}
	}
	return res[0].Interface(), nil
}