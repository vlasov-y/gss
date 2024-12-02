package pkg

import "fmt"

func PrefixError(prefix string, err error) error {
	if err != nil {
		err = fmt.Errorf("%s: %w", prefix, err)
	}
	return err
}
