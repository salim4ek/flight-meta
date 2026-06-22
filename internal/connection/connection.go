// Package connection assesses whether a layover gives enough time to make the
// next flight (Minimum Connection Time). Defaults are indicative/typical values
// averaged across carriers; a self-transfer (separate tickets) needs more time
// because the traveller must collect and re-check baggage and re-clear security.
package connection

import "flightmeta/internal/offer"

// Indicative MCT defaults (minutes), international transfer.
const (
	OfficialMinMinutes     = 90  // single ticket, protected connection
	SelfTransferMinMinutes = 150 // separate tickets: re-check bags + security
)

// Risk levels written onto offer.Layover.Risk.
const (
	RiskSafe       = "safe"
	RiskRisky      = "risky"
	RiskInfeasible = "infeasible"
)

// RequiredMinutes returns the minimum connection time for a layover type.
func RequiredMinutes(selfTransfer bool) int {
	if selfTransfer {
		return SelfTransferMinMinutes
	}
	return OfficialMinMinutes
}

// Enrich sets Risk on each layover by comparing its duration to the required
// MCT. A self-transfer through a country that requires a transit visa forces a
// landside exit the traveller can't legally make, so it is infeasible.
func Enrich(offers []offer.Offer) {
	for i := range offers {
		for j := range offers[i].Layovers {
			lay := &offers[i].Layovers[j]
			required := RequiredMinutes(lay.SelfTransfer)
			mins := int(lay.Duration.Minutes())

			switch {
			case lay.SelfTransfer && lay.VisaStatus == "visa_required":
				lay.Risk = RiskInfeasible
			case mins >= required:
				lay.Risk = RiskSafe
			case mins*10 >= required*7: // >= 70% of required
				lay.Risk = RiskRisky
			default:
				lay.Risk = RiskInfeasible
			}
		}
	}
}

// WorstRisk returns the most severe layover risk in an offer ("" if none).
func WorstRisk(o offer.Offer) string {
	worst := ""
	rank := map[string]int{"": 0, RiskSafe: 1, RiskRisky: 2, RiskInfeasible: 3}
	for _, l := range o.Layovers {
		if rank[l.Risk] > rank[worst] {
			worst = l.Risk
		}
	}
	return worst
}
