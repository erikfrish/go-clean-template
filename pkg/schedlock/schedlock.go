package schedlock

import (
	"fmt"
)

type Repository interface {
	IsFirst(jobName string) (bool, error)
	Release(jobName string) error
}

func DoOnce(jobName string, f func(), r Repository) error {
	first, err := r.IsFirst(jobName)
	if err != nil {
		return fmt.Errorf("r.IsFirst: %w", err)
	}
	if first {
		f()
		err = r.Release(jobName)
		if err != nil {
			return fmt.Errorf("r.Release: %s: %w", jobName, err)
		}
	}
	return nil
}
