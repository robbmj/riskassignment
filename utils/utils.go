package utils

import (
	"bufio"
	"fmt"
	"github.com/robbmj/riskassignment/threat"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadFile(lifeTime float64) threat.Threats {

	file, _ := os.Open("data.csv")
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')

	threats := make(threat.Threats, 0)

	for i := 0; err != io.EOF; i++ {

		line = strings.Trim(line, "\n")

		values := strings.Split(line, ",")

		t := threat.Threat{}
		t.Name = values[0]
		t.AssetValue, err = strconv.ParseFloat(values[1], 64)
		t.ExposureFactor, err = strconv.ParseFloat(values[2], 64)
		t.RateOfOcurance, err = strconv.ParseFloat(values[3], 64)
		t.OneTimeCost, err = strconv.ParseFloat(values[4], 64)
		t.LifeTimeOfControl = lifeTime

		threats = append(threats, t)

		line, err = reader.ReadString('\n')
	}

	return threats
}

func WriteMyAssignment() string {
	html :=
		`
	<!DOCTYPE html>
	<html>
	<head></head>
	<body>
		<div style="width:97%;padding:10px;">
		<h3>Michael Robb<br>Computer and Network Security: Assignment 1</h3>
		<h3>Part 1</h3>
		<h4>Questions 1 &amp; 2: Calculate SLE, ALE, ROI</h4>`

	threats := ReadFile(3.0)
	threats.SortBy(threat.ByName{})
	html += MakeTable(&threats)
	html += "<br><h4>Question 3: Controls which should be purchased by the company</h4>"
	html += `<p>Controls With a positive ROI will save the company money over a three year period.
			 	The following table lists those controls`

	goodControls := threats.FindPositiveROI()
	goodControls.SortBy(threat.ByReturnOnInvestment{})
	html += MakeTable(&goodControls)

	html += "<br><h4>Question 4: Witch conrols should be implemented on a budget of $15000</h4>"
	html += `<p>The controls that should be implemented are the ones that when combined 
				have the maximum ALU with a total cost of implementation less than $15000.
			</p>
			<p>The follwoing table lists those controls, underneath the table you can find the total ROI, Savings and Cost</p>`

	bestBuys, cost, savings, roi := goodControls.DeterminBestPurchases(15000.0, threat.ByAnnualizedLossExpectancy{})

	html += MakeTable(&bestBuys)

	stats := `<br>Total Cost of Controls: %.02f, 
				Total savings over 3 years: %.02f, Total ROI: %.02f<br>`

	html += fmt.Sprintf(stats, cost, savings, roi)

	html += "<h3>Part 2</h3><h4>Research and describe a leading risk assessment methodology</h4>"

	file, _ := os.Open("octave.html")
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')

	for err != io.EOF {
		html += line
		line, err = reader.ReadString('\n')
	}

	html += "</div></body></html>"
	return html

}

func WriteOut(html string) {
	err := os.Remove("output.html")

	if err != nil {
		panic(err)
	}

	file, err := os.Create("output.html")

	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(file)
	writer.Write([]byte(html))
	writer.Flush()
}

func MakeTable(threats *threat.Threats) string {
	return tableHeader() + tableBody(threats) + closeTable()
}

func tableHeader() string {
	return `<table border='1' cellspaceing='1' cellpadding='1' style="margin-left:auto;margin-right:auto;">
			<tr>
				<th>Threat</th>
				<th>Asset Value</th>
				<th>Exposure Factor</th>
				<th>Rate of Occurance</th>
				<th>One Time Cost</th>
				<th>Single Loss Expectancy</th>
				<th>Annualized Loss Expectancy</th>
				<th>Return On Investment after 3 years</th>
			</tr>`
}

func tableBody(threats *threat.Threats) string {
	html := ""
	for _, threat := range *threats {
		tableRow := `<tr><td>%s</td><td>%.02f</td><td>%.02f</td><td>%.02f</td><td>%.02f</td>
					 <td>%.02f</td><td>%.02f</td><td>%.02f</td></tr>`

		html += fmt.Sprintf(tableRow, threat.Name, threat.AssetValue, threat.ExposureFactor, threat.RateOfOcurance,
			threat.OneTimeCost,
			threat.SingleLossExpectancy(),
			threat.AnnualizedLossExpectancy(),
			threat.ReturnOnInvestment())
	}
	return html
}

func closeTable() string {
	return "</table>"
}
