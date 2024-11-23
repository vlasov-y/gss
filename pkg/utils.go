package pkg

import "fmt"

func PrefixError(prefix string, err error) error {
	if err != nil {
		err = fmt.Errorf("%s: %v", prefix, err)
	}
	return err
}
