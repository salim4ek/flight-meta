// Package search fans out a query to all configured sources concurrently,
// enforces a per-source timeout, then merges, filters and ranks the results.
package search

import (
	"context"
	"log/slog"
	"time"

	"flightmeta/internal/connection"
	"flightmeta/internal/offer"
	"flightmeta/internal/rank"
	"flightmeta/internal/sources"
	"flightmeta/internal/visa"
)

// Orchestrator coordinates the concurrent fan-out across data sources.
type Orchestrator struct {
	sources       []sources.Adapter
	sourceTimeout time.Duration
	visa          *visa.Resolver // optional; enriches layovers with transit-visa info
	log           *slog.Logger
}

// New builds an orchestrator over the given source adapters. resolver may be
// nil (no transit-visa enrichment).
func New(log *slog.Logger, sourceTimeout time.Duration, resolver *visa.Resolver, srcs ...sources.Adapter) *Orchestrator {
	return &Orchestrator{sources: srcs, sourceTimeout: sourceTimeout, visa: resolver, log: log}
}

// Result is the merged, ranked search response plus per-source diagnostics.
type Result struct {
	Offers         []offer.Offer `json:"offers"`
	Sources        []SourceStat  `json:"sources"`
	VisaDisclaimer string        `json:"visaDisclaimer,omitempty"`
}

// SourceStat reports how each source performed (count, error, latency).
type SourceStat struct {
	Name   string `json:"name"`
	Count  int    `json:"count"`
	Err    string `json:"error,omitempty"`
	TookMs int64  `json:"tookMs"`
}

// Search queries every source concurrently. A slow or failing source is
// isolated by its own timeout and never blocks the others; its error is
// surfaced in SourceStat rather than failing the whole request.
func (o *Orchestrator) Search(ctx context.Context, q sources.Query) Result {
	type partial struct {
		stat   SourceStat
		offers []offer.Offer
	}
	ch := make(chan partial, len(o.sources))

	for _, src := range o.sources {
		go func(src sources.Adapter) {
			sctx, cancel := context.WithTimeout(ctx, o.sourceTimeout)
			defer cancel()

			start := time.Now()
			offers, err := src.Search(sctx, q)
			stat := SourceStat{
				Name:   src.Name(),
				Count:  len(offers),
				TookMs: time.Since(start).Milliseconds(),
			}
			if err != nil {
				stat.Err = err.Error()
				o.log.Warn("source failed", "source", src.Name(), "err", err)
			}
			ch <- partial{stat: stat, offers: offers}
		}(src)
	}

	var all []offer.Offer
	stats := make([]SourceStat, 0, len(o.sources))
	for range o.sources {
		p := <-ch
		stats = append(stats, p.stat)
		all = append(all, p.offers...)
	}

	// Enrichment must run before ranking: the OnlyVisaFreeTransit filter reads
	// VisaStatus, and HideInfeasible reads layover Risk (which itself depends on
	// VisaStatus for self-transfer). So: visa first, then connection.
	res := Result{Sources: stats}
	if o.visa != nil {
		o.visa.Enrich(all, q.Passport)
		res.VisaDisclaimer = o.visa.Disclaimer()
	}
	connection.Enrich(all)
	res.Offers = rank.Apply(all, q.Filters)
	return res
}

// CalendarDay is the cheapest price found for a single date.
type CalendarDay struct {
	Date       string `json:"date"` // YYYY-MM-DD
	PriceMinor int64  `json:"priceMinor"`
	Currency   string `json:"currency"`
	HasOffers  bool   `json:"hasOffers"`
	Cheapest   bool   `json:"cheapest"`
}

// CalendarResult is a price-per-date window for flexible-date search.
type CalendarResult struct {
	Days []CalendarDay `json:"days"`
}

// Calendar runs a search for each day in [depart-window, depart+window] and
// returns the cheapest price per day, flagging the overall minimum. window is
// clamped to [1, 14]. Filters from q apply to each day's cheapest.
func (o *Orchestrator) Calendar(ctx context.Context, q sources.Query, window int) CalendarResult {
	if window < 1 {
		window = 1
	}
	if window > 14 {
		window = 14
	}
	dates := make([]time.Time, 0, window*2+1)
	for d := -window; d <= window; d++ {
		dates = append(dates, q.DepartDate.AddDate(0, 0, d))
	}

	type dayResult struct {
		idx int
		day CalendarDay
	}
	ch := make(chan dayResult, len(dates))
	for i, date := range dates {
		go func(i int, date time.Time) {
			dq := q
			dq.DepartDate = date
			r := o.Search(ctx, dq)
			cd := CalendarDay{Date: date.Format("2006-01-02")}
			if len(r.Offers) > 0 {
				cd.HasOffers = true
				cd.PriceMinor = r.Offers[0].PriceMinor // Offers are cheapest-first
				cd.Currency = r.Offers[0].Currency
			}
			ch <- dayResult{idx: i, day: cd}
		}(i, date)
	}

	days := make([]CalendarDay, len(dates))
	for range dates {
		dr := <-ch
		days[dr.idx] = dr.day
	}

	// Flag the cheapest day with offers.
	minIdx, minPrice := -1, int64(0)
	for i, d := range days {
		if !d.HasOffers {
			continue
		}
		if minIdx == -1 || d.PriceMinor < minPrice {
			minIdx, minPrice = i, d.PriceMinor
		}
	}
	if minIdx >= 0 {
		days[minIdx].Cheapest = true
	}
	return CalendarResult{Days: days}
}
