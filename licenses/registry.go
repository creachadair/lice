// Copyright (C) 2018, Michael J. Fromberger
// All Rights Reserved.

package licenses

import "log"

type registry struct {
	known []License // ordered by slug
}

func (r *registry) fetch(slug string) *License {
	if i, ok := r.lookup(slug); ok {
		out := r.known[i]
		return &out
	}
	return nil
}

func (r *registry) visit(f func(License)) {
	for _, lic := range r.known {
		f(lic)
	}
}

func (r *registry) lookup(slug string) (int, bool) {
	i, j := 0, len(r.known)-1
	for i <= j {
		m := (i + j) / 2
		cur := r.known[m].Slug
		if slug == cur {
			return m, true
		} else if slug < cur {
			j = m
		} else {
			i = m + 1
		}
	}
	return i, false
}

func (r *registry) insert(lic License) bool {
	i, ok := r.lookup(lic.Slug)
	if ok {
		return false
	}
	r.known = append(append(r.known[:i], lic), r.known[i:]...)
	return true
}

var global = new(registry)

// Register records a new license in the registry, using its slug as a
// key. This function will panic if the license slug is empty, or if the slug
// is already registered to a different license.
func Register(lic License) {
	if lic.Slug == "" {
		log.Panic("empty license slug")
	} else if !global.insert(lic) {
		log.Panicf("duplicate registrations for slug %q", lic.Slug)
	}
}

// Lookup returns the license information for the specified slug, or nil if no
// such license is registered.
func Lookup(slug string) *License { return global.fetch(slug) }

// List calls f for each registered license.  Licenses are visited in
// lexicographic order by slug.
func List(f func(License)) { global.visit(f) }
