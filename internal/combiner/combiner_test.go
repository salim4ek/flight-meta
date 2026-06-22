package combiner

import (
	"context"
	"testing"
	"time"

	"flightmeta/internal/offer"
	"flightmeta/internal/sources"
	"flightmeta/internal/sources/mockleg"
)

func TestCombinerProducesDirectAndSelfTransfer(t *testing.T) {
	c := New(mockleg.New(), []Hub{{"IST", 180}})
	q := sources.Query{
		Origin: "MOW", Destination: "BKK",
		DepartDate: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC),
		Currency:   "RUB",
	}
	offers, err := c.Search(context.Background(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(offers) != 2 {
		t.Fatalf("want 2 offers (direct + via IST), got %d", len(offers))
	}

	var direct, self *offer.Offer
	for i := range offers {
		switch offers[i].Connection {
		case offer.Official:
			direct = &offers[i]
		case offer.SelfTransfer:
			self = &offers[i]
		}
	}
	if direct == nil || self == nil {
		t.Fatal("expected one official and one self-transfer offer")
	}
	if len(self.Segments) != 2 || len(self.Layovers) != 1 {
		t.Fatalf("self-transfer should have 2 segments and 1 layover, got %d/%d", len(self.Segments), len(self.Layovers))
	}
	if self.Layovers[0].Airport != "IST" || !self.Layovers[0].SelfTransfer {
		t.Fatal("layover should be a self-transfer at IST")
	}
	if !self.Segments[1].DepartUTC.After(self.Segments[0].ArriveUTC) {
		t.Fatal("leg 2 must depart after leg 1 arrives")
	}
	if !self.Unique {
		t.Fatal("self-transfer offer should be marked unique")
	}
}

func TestCombinerSkipsHubEqualToEndpoint(t *testing.T) {
	c := New(mockleg.New(), []Hub{{"BKK", 180}})
	q := sources.Query{Origin: "MOW", Destination: "BKK", DepartDate: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC)}
	offers, _ := c.Search(context.Background(), q)
	for _, o := range offers {
		if o.Connection == offer.SelfTransfer {
			t.Fatal("should not build a self-transfer via a hub equal to the destination")
		}
	}
}
