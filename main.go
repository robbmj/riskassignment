
package main

import (
	"fmt"
	"github.com/robbmj/riskassignment/utils"
	"github.com/robbmj/riskassignment/threat"
)

func main() {
	assignment := utils.WriteMyAssignment()
	utils.WriteOut(assignment)
}

func devFunction() {
	threats := utils.ReadFile(3.0)
	//fmt.Println(threats)
	//fmt.Println(threats[0].SingleLossExpectancy())
	//threats.SortByROI()
	threatsWithPositiveROI := threats.FindPositiveROI()
	//html := utils.MakeTable(&threatsWithPositiveROI)

	threatsWithPositiveROI.SortByALE()
	html := utils.MakeTable(&threatsWithPositiveROI)

	bestBuys, cost, savings, roi := threatsWithPositiveROI.DeterminBestPurchases(15000.0, threat.ByReturnOnInvestment{})

	title := "<br><h4>Best Fit Controls, Total Cost of Controls: %.02f, Total Loss over 3 years: %.02f, Total ROI: %.02f</h4><br>"
	html += fmt.Sprintf(title, cost, savings, roi)
	html += "<br />" + utils.MakeTable(&bestBuys)
	utils.WriteOut(html)
}