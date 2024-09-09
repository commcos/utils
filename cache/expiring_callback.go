/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package cache

import utilclock "github.com/commcos/utils/clock"

type DeleteHook func(key interface{}, value interface{})

// NewExpiringWithClock is like NewExpiring but allows passing in a custom
// clock for testing.
func NewExpiringWithCallback(hook DeleteHook) *Expiring {
	return &Expiring{
		clock:   utilclock.RealClock{},
		cache:   make(map[interface{}]entry),
		delHook: hook,
	}
}
