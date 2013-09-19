
package threat

import (
	"sort"
	//"fmt"
)

type Threat struct {
	Name           string
	AssetValue     float64
	ExposureFactor float64
	RateOfOcurance float64
	OneTimeCost    float64
	LifeTimeOfControl float64
}

type Threats []Threat

type ByReturnOnInvestment struct {
	Threats
}

type ByAnnualizedLossExpectancy struct {
	Threats
}

func (t Threat) SingleLossExpectancy() float64 {
	return float64(t.AssetValue * t.ExposureFactor)
}

func (t Threat) AnnualizedLossExpectancy() float64 {
	return t.SingleLossExpectancy() * t.RateOfOcurance
}

func (t Threat) ReturnOnInvestment() float64 {
	return 100 * ((t.LifeTimeOfControl * t.AnnualizedLossExpectancy() - t.OneTimeCost) / t.OneTimeCost)
}

func (ts *Threats) Sort() {
	sort.Sort(*ts)
}

func (ts *Threats) SortBy(by interface{}) {
	switch by.(type) {
	case ByAnnualizedLossExpectancy:
		sort.Sort(ByAnnualizedLossExpectancy{*ts})
	case ByReturnOnInvestment:
		sort.Sort(ByReturnOnInvestment{*ts})
	}
}

func (ts *Threats) SortByALE() {
	sort.Sort(ByAnnualizedLossExpectancy{*ts})
}

func (ts Threats) Len() int {
	return len(ts)
}

func (ts Threats) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts Threats) Less(i, j int) bool {
	return ts[i].Name < ts[j].Name
}

func (t ByReturnOnInvestment) Less(i, j int) bool {
	return t.Threats[i].ReturnOnInvestment() > t.Threats[j].ReturnOnInvestment()
}

func (t ByAnnualizedLossExpectancy) Less(i, j int) bool {
	return t.Threats[i].AnnualizedLossExpectancy() > t.Threats[j].AnnualizedLossExpectancy()
}

func (ts *Threats) FindPositiveROI() Threats {

	subSet := make(Threats, 0)
	for _, threat := range *ts {
		if threat.ReturnOnInvestment() > 0 {
			subSet = append(subSet, threat)
		}
	}
	return subSet
}

// This is a very simple algorithum to find the optimal Purchases
// it just happens to work with the data set
// TODO: implement this properly
func (ts *Threats) DeterminBestPurchases(budget float64, by interface{}) (Threats, float64, float64, float64) {

	ts.SortBy(by)

	subSet := make(Threats, 0)

	totalCost, totalROI, totalSavings, lifeTime := 0.0, 0.0, 0.0, 0.0

	for _, threat := range *ts {
		if totalCost + threat.OneTimeCost <= budget {

			totalCost += threat.OneTimeCost
			totalROI += threat.ReturnOnInvestment()
			totalSavings += threat.AnnualizedLossExpectancy()
			lifeTime = threat.LifeTimeOfControl

			subSet = append(subSet, threat)
		}
	}

	totalSavings *= lifeTime

	return subSet, totalCost, totalSavings, totalROI
}
