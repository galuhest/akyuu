package akyuu_test

import (
	"testing"
	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/bukalapak/dallimin"

	"github.com/galuhest/akyuu"
)

type Example struct {
	Id int `json:id`
}

func TestCreateNewObject(t *testing.T) {
	option := dallimin.Option{
	    CheckAlive: true,
	    Failover: true,
	}

	servers := []string{
		"127.0.0.1:11211",
	}

	ss, err := dallimin.New(servers, option)

	_, err = akyuu.New(ss)
	assert.Equal(t, nil, err)
}

func TestFetchItem(t * testing.T) {
	option := dallimin.Option{
	    CheckAlive: true,
	    Failover: true,
	}

	servers := []string{
		"127.0.0.1:11211",
	}

	ss, err := dallimin.New(servers, option)

	ak, err := akyuu.New(ss)

	_, err = ak.Fetch(returnStruct, "example", 100)
	assert.Equal(t, nil, err)
	
	_, err = ak.Fetch(returnStruct, "example", 100)
	assert.Equal(t, nil, err)

	_, err = ak.Fetch(returnError, "error", 100)
	assert.Equal(t, nil, err)
}

func returnError() (int, error) {
	return 1, errors.New("an error")
}

func returnStruct() (int, error) {
	return 1, nil
}