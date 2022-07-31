/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package flowcontrol

// IntMax returns the maximum of the params
func IntMax(a, b int) int {
	if b > a {
		return b
	}
	return a
}

// Int64Max returns the maximum of the params
func Int64Max(a, b int64) int64 {
	if b > a {
		return b
	}
	return a
}

// Int64Min returns the minimum of the params
func Int64Min(a, b int64) int64 {
	if b < a {
		return b
	}
	return a
}
