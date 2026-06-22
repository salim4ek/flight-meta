package sources

import (
	"context"
	"time"
)

// LegQuote is the cheapest price for a single city-pair leg on a date. The
// combiner schedules concrete times itself, so a quote only needs duration +
// price + a deep link to buy that leg.
type LegQuote struct {
	Carrier      string
	FlightNumber string
	DurationMin  int
	PriceMinor   int64
	Currency     string
	DeepLink     string
}

// LegSource provides per-leg cheapest prices. The own combiner builds
// self-transfer itineraries by stitching legs from a LegSource across hubs —
// this is how we surface routes that single-ticket aggregators don't show,
// without depending on a virtual-interlining vendor. Real implementations wrap
// a price feed (e.g. Travelpayouts cheapest-prices) behind this interface.
type LegSource interface {
	Name() string
	CheapestLeg(ctx context.Context, from, to string, date time.Time, q Query) (LegQuote, bool, error)
}
