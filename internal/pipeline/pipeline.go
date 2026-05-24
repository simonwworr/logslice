// Package pipeline wires together all processing stages into a single
// cohesive execution unit. It connects the reader, slicer, filter, sampler,
// dedup, limiter, highlight, and formatter components so callers only need
// to configure a Config and call Run.
package pipeline

import (
	"context"
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/dedup"
	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/formatter"
	"github.com/yourorg/logslice/internal/highlight"
	"github.com/yourorg/logslice/internal/limiter"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/sampler"
	"github.com/yourorg/logslice/internal/slicer"
	"github.com/yourorg/logslice/internal/stats"
)

// Config holds all user-facing options that control pipeline behaviour.
type Config struct {
	// Time range boundaries (RFC3339 or space-separated).
	From string
	To   string

	// Filter options.
	Keywords []string
	Patterns []string
	Invert   bool

	// Sampling: keep every Nth line (0 = disabled).
	SampleN int

	// Dedup mode: "" (none), "consecutive", or "global".
	DedupMode string

	// Output limits.
	MaxLines int64
	MaxBytes int64

	// Highlight terms in output.
	HighlightKeywords []string
	HighlightPatterns []string

	// Output format: "raw", "numbered", or "json".
	Format string

	// Destination writer.
	Out io.Writer
}

// Result summarises what happened after Run returns.
type Result struct {
	Stats *stats.Stats
}

// Run executes the full processing pipeline, reading lines from src and
// writing matching output to cfg.Out. It returns a Result with collected
// statistics, or an error if any stage could not be initialised.
func Run(ctx context.Context, src io.Reader, cfg Config) (*Result, error) {
	// --- time range ---
	var from, to interface{ IsZero() bool }
	_ = from
	_ = to
	sliceFrom, err := parser.ParseTimestamp(cfg.From)
	if err != nil && cfg.From != "" {
		return nil, fmt.Errorf("pipeline: invalid --from: %w", err)
	}
	sliceTo, err := parser.ParseTimestamp(cfg.To)
	if err != nil && cfg.To != "" {
		return nil, fmt.Errorf("pipeline: invalid --to: %w", err)
	}

	// --- filter ---
	filterOpts := []filter.Option{}
	if len(cfg.Keywords) > 0 {
		filterOpts = append(filterOpts, filter.WithKeywords(cfg.Keywords...))
	}
	if len(cfg.Patterns) > 0 {
		filterOpts = append(filterOpts, filter.WithPatterns(cfg.Patterns...))
	}
	if cfg.Invert {
		filterOpts = append(filterOpts, filter.WithInvert())
	}
	f, err := filter.New(filterOpts...)
	if err != nil {
		return nil, fmt.Errorf("pipeline: filter: %w", err)
	}

	// --- sampler ---
	var smp *sampler.Sampler
	if cfg.SampleN > 1 {
		smp, err = sampler.New(cfg.SampleN)
		if err != nil {
			return nil, fmt.Errorf("pipeline: sampler: %w", err)
		}
	}

	// --- dedup ---
	var dd *dedup.Deduplicator
	if cfg.DedupMode != "" {
		dd, err = dedup.New(cfg.DedupMode)
		if err != nil {
			return nil, fmt.Errorf("pipeline: dedup: %w", err)
		}
	}

	// --- limiter ---
	lim, err := limiter.New(cfg.MaxLines, cfg.MaxBytes)
	if err != nil {
		return nil, fmt.Errorf("pipeline: limiter: %w", err)
	}

	// --- highlighter ---
	highlightOpts := []highlight.Option{}
	for _, kw := range cfg.HighlightKeywords {
		highlightOpts = append(highlightOpts, highlight.WithKeyword(kw))
	}
	for _, pat := range cfg.HighlightPatterns {
		highlightOpts = append(highlightOpts, highlight.WithPattern(pat))
	}
	hl, err := highlight.New(highlightOpts...)
	if err != nil {
		return nil, fmt.Errorf("pipeline: highlight: %w", err)
	}

	// --- formatter ---
	fmt_ := cfg.Format
	if fmt_ == "" {
		fmt_ = "raw"
	}
	fw, err := formatter.New(cfg.Out, fmt_)
	if err != nil {
		return nil, fmt.Errorf("pipeline: formatter: %w", err)
	}

	// --- stats ---
	st := stats.New()

	// --- main loop ---
	sc := slicer.Slice(src, sliceFrom, sliceTo)
	for {
		select {
		case <-ctx.Done():
			st.Finish()
			return &Result{Stats: st}, ctx.Err()
		default:
		}

		line, ok := sc()
		if !ok {
			break
		}

		st.RecordLine(int64(len(line)))

		// filter
		matched, ferr := f.Match(line)
		if ferr != nil {
			return nil, fmt.Errorf("pipeline: filter match: %w", ferr)
		}
		if !matched {
			continue
		}

		// sampler
		if smp != nil && !smp.Keep() {
			continue
		}

		// dedup
		if dd != nil && dd.IsDuplicate(line) {
			continue
		}

		// highlight
		outLine := hl.Apply(line)

		// limiter
		if lerr := lim.Add(outLine); lerr != nil {
			break
		}

		// write
		if werr := fw.WriteLine(outLine); werr != nil {
			return nil, fmt.Errorf("pipeline: write: %w", werr)
		}
	}

	st.Finish()
	return &Result{Stats: st}, nil
}
