// Package mockleg is a deterministic LegSource: it fabricates a cheapest-price
// quote for any city-pair and date. It lets the combiner run end-to-end before
// a real price feed (Travelpayouts etc.) is wired behind the same interface.
package mockleg

import (
	"context"
	"hash/fnv"
	"time"

	"flightmeta/internal/sources"
)

type Source struct{}

func New() *Source { return &Source{} }

func (s *Source) Name() string { return "mockleg" }

var carriers = []string{"SU", "TK", "EK", "QR", "U6", "CA", "6E", "FZ", "EY", "JU"}

// CheapestLeg returns a deterministic quote. Price depends on the route hash
// and the weekday (so flexible-date calendars show real mid-week-cheaper
// variation); duration depends on the route hash.
func (s *Source) CheapestLeg(ctx context.Context, from, to string, date time.Time, q sources.Query) (sources.LegQuote, bool, error) {
	if err := ctx.Err(); err != nil {
		return sources.LegQuote{}, false, err
	}
	cur := q.Currency
	if cur == "" {
		cur = "RUB"
	}
	h := hash32(from + ">" + to)

	base := int64(900000 + h%2600000)         // 9 000 .. 35 000
	price := base + weekdayAdjust(date.Weekday())
	if price < 500000 {
		price = 500000
	}
	dur := 90 + int(h%420) // 1h30 .. 8h30
	carrier := carriers[h%uint32(len(carriers))]

	return sources.LegQuote{
		Carrier:      carrier,
		FlightNumber: carrier + flightNo(h),
		DurationMin:  dur,
		PriceMinor:   price,
		Currency:     cur,
		DeepLink:     "https://example-partner.test/leg/" + from + "-" + to,
	}, true, nil
}

// weekdayAdjust makes some days cheaper than others (in minor units).
func weekdayAdjust(wd time.Weekday) int64 {
	switch wd {
	case time.Tuesday:
		return -250000
	case time.Wednesday:
		return -200000
	case time.Thursday:
		return -50000
	case time.Friday:
		return 150000
	case time.Saturday:
		return 400000
	case time.Sunday:
		return 300000
	default: // Monday
		return 0
	}
}

func hash32(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

func flightNo(h uint32) string {
	n := 100 + h%8900
	return itoa(int(n))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [6]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
