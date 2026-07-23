package lyrics

import (
	"testing"
	"time"
)

type fakeProvider struct {
	delay time.Duration
	lines []Line
	err   error
}

func (f fakeProvider) Lyrics(artist, track string) ([]Line, error) {
	time.Sleep(f.delay)
	return f.lines, f.err
}

func TestMultiTakesFastestSuccess(t *testing.T) {
	slow := fakeProvider{delay: 50 * time.Millisecond, lines: []Line{{Words: "from slow"}}}
	fast := fakeProvider{delay: 5 * time.Millisecond, lines: []Line{{Words: "from fast"}}}

	p := NewMulti(slow, fast)
	start := time.Now()
	lines, err := p.Lyrics("a", "b")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0].Words != "from fast" {
		t.Fatalf("expected the faster provider's result regardless of order, got %+v", lines)
	}
	// should return as soon as the fast one answers (~5ms), not wait for
	// the slow one just because it's listed first
	if elapsed >= 20*time.Millisecond {
		t.Fatalf("didn't return early on the fastest result: %v", elapsed)
	}
}

func TestMultiFallsBackWhenFirstEmpty(t *testing.T) {
	empty := fakeProvider{delay: 5 * time.Millisecond}
	backup := fakeProvider{delay: 10 * time.Millisecond, lines: []Line{{Words: "backup"}}}

	p := NewMulti(empty, backup)
	lines, err := p.Lyrics("a", "b")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0].Words != "backup" {
		t.Fatalf("expected fallback to backup provider, got %+v", lines)
	}
}

func TestMultiWaitsPastEmptyForLateSuccess(t *testing.T) {
	fastEmpty := fakeProvider{delay: 5 * time.Millisecond}
	slowSuccess := fakeProvider{delay: 20 * time.Millisecond, lines: []Line{{Words: "eventually"}}}

	p := NewMulti(fastEmpty, slowSuccess)
	lines, err := p.Lyrics("a", "b")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0].Words != "eventually" {
		t.Fatalf("expected to keep waiting past an empty result for a later success, got %+v", lines)
	}
}
