package akyuu

import (
	"fmt"
	"reflect"
	"encoding/json"
	"github.com/pkg/errors"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bukalapak/dallimin"
)

type Task interface{}

type Connection struct {
	mem *memcache.Client
}

type Value struct {
	Value interface{} `json:"value"`
}

func New(dalli *dallimin.Ring) (*Connection, error) {
	return &Connection {
		mem: memcache.NewFromSelector(dalli),
	}, nil
}

func (c Connection) Fetch(task Task, key string, expiration int32) (*memcache.Item, error) {
	item, err := c.mem.Get(key)
	// key exists
	if err == nil {
		return item, errors.New("Key already exists")
	}

	res, err := runInterface(task)
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

func runInterface(intf interface{}, ret ...interface{}) (interface{}, error) {
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
			return nil, errors.Wrap(errors.New("Error is not nil"), fmt.Sprintf("%s",errVal))
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
	if retAmt == 2 {
		err := res[1].Interface()

		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s",err))
		}
	}
	return res[0].Interface(), nil
}