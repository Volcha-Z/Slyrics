package lyrics

// NewMulti queries several providers at once and returns whichever finds
// something first. All three current providers do the same artist
// verification, so there's no correctness reason to prefer one over another
// — only waiting on a fixed "priority" provider that happens to be slow (or
// hung on its own internal fallback/timeout) while a faster one already has
// the answer would just make lyrics take way longer to show up than they
// need to.
func NewMulti(providers ...Provider) Provider {
	return &multiProvider{providers}
}

type multiProvider struct {
	providers []Provider
}

type multiResult struct {
	lines []Line
	err   error
}

func (m *multiProvider) Lyrics(artist, track string) ([]Line, error) {
	ch := make(chan multiResult, len(m.providers))
	for _, p := range m.providers {
		go func(p Provider) {
			lines, err := p.Lyrics(artist, track)
			ch <- multiResult{lines, err}
		}(p)
	}

	var lastErr error
	for range m.providers {
		r := <-ch
		if len(r.lines) > 0 {
			return r.lines, nil
		}
		if r.err != nil {
			lastErr = r.err
		}
	}
	return nil, lastErr
}
