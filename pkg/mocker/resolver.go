// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mocker

// resolver is a package cache and resolver.
type resolver struct {
	cache []*gopkg
}

// resolve retrieves the package for the given "want" from the cache or finds
// it using gopkg.resolve. It caches the result on a successful lookup and returns
// the package and any error encountered.
func (res *resolver) resolve(want *gopkg) error {
	if want.resolved {
		return nil
	}
	for _, have := range res.cache {
		if have.equal(want) {
			want.from(have)
			return nil
		}
	}
	if err := want.resolve(); err != nil {
		return err
	}
	res.cache = append(res.cache, want)
	return nil
}
