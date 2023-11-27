//go:build !gotip

package main

import (
	"context"
	"fmt"
	"strings"
)

// StreamPages streams pages and filter dog content.
func (s *Service) StreamPages(ctx context.Context) func(func(Page, error) bool) {
	return func(yield func(Page, error) bool) {
		var zero Page
		s.db.StreamPages(ctx)(func(p Page, err error) bool {
			if err != nil {
				yield(zero, fmt.Errorf("db: %v", err))
				return false
			}

			// We don't want to know about dogs.
			title := strings.ToLower(p.Title)
			if strings.Contains(title, "dog") {
				return true
			}

			return yield(p, err)
		})
	}
}
