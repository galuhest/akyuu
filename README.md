# Akyuu

## Description

Akyuu is a wrapper for gomemcache that accept function with any return value. Akyuu use [dallimin](https://github.com/bukalapak/dallimin/) for it's connectivity.

## Owner

SRE - Library & Service

## Contact and On-Call Information

See [Contact and On-Call Information](https://bukalapak.atlassian.net/wiki/display/INF/Contact+and+On-Call+Information)

## Link

- [dallimin](https://github.com/bukalapak/dallimin/)
- [gomemcache](https://github.com/bradfitz/gomemcache)
- [errors](https://github.com/pkg/errors)

## Usage

Akyuu is composed of 2 function, `New` and `Fetch`.

`New` need a dallimin's Ring struct as parameter.
```golang
func New(*dallimin.Ring) (*Connection, error)
```
`Fetch` need 3 parameter, a Task in the form of `function` or `value`, a string that will be used as `key`, and an integer as `exipration time`.
```golang
func Fetch(interface{}, string, int32)
```

### Example


```golang
import (
	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/bukalapak/dallimin"
)

type Example struct {
	Id int `json:"id"`
}

func main() {
	
	option := dallimin.Option{
	    CheckAlive: true,
	    Failover: true,
	}

	servers := []string{
		"127.0.0.1:11211",
	}

	ss, err := dallimin.New(servers, option)

	ak, err := akyuu.New(ss)

	// function call example

	//as a function
	res, err := ak.Fetch(FunctionOne, "example", 10)
	
	//a value
	res, err = ak.Fetch(FunctionOne(), "example2", 10)

	strt := &Example{Id: 10}
	res, err = ak.Fetch(FunctionFour(strt), "example3", 10)

	res, err = ak.Fetch(FunctionFive(2), "example", 10)

}

// ********* valid function signatures *********

// function that return one value
func FunctionOne() int {
	return 1
}

// function that return a struct
func FunctionTwo() *Example {
	return &Example{Id: 1}
}

// function that return a value and an error
func FunctionThree() (int, error) {
	return 1, errors.New("an error")
}

// function that return a struct and an error
func FunctionFour() (*Example, error) {
	return &Example{Id: 1}, nil
}

// function with parameter(s) that return a value and an error
func FunctionFive(i int) (int, error) {
	return i+1, nil
}

// ********* invalid function signature (will return an error) *********

// function that return more than 2 value
func FunctionSix() (int, string, error) {
	return 1, "invalid", nil
} 

// function that return 2 value but the second value is not an error
// error will be generated if the function does generate an error
// but the second value is not an error
func FunctionSeven() (int, string) {
	return 1, "example"
}

```