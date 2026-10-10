// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mocker

// resolver is a package cache and resolver.
type resolver struct {
	cache []*gopkg
}

// resolve fills "want" from the cache or resolves it using gopkg.resolve,
// caching it on success. It returns the resolution error, if any.
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
