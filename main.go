package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

func formatUSD(n float64) string {
	return fmt.Sprintf("$%s", commaFmt(n))
}

func commaFmt(n float64) string {
	s := fmt.Sprintf("%.0f", n)
	out := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out += ","
		}
		out += string(c)
	}
	return out
}

var indexTmpl = template.Must(template.New("index.html").Funcs(template.FuncMap{
	"formatUSD": formatUSD,
}).ParseFiles("templates/index.html"))

var mortgageTmpl = template.Must(template.New("mortgage.html").Funcs(template.FuncMap{
	"formatUSD": formatUSD,
}).ParseFiles("templates/mortgage.html"))

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/mortgage", mortgageHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func mortgageHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err == nil {
			amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
			rate, _ := strconv.ParseFloat(r.FormValue("rate"), 64)
			years, _ := strconv.Atoi(r.FormValue("years"))

			monthlyRate := rate / 100 / 12
			numPayments := years * 12

			monthlyPayment := 0.0
			if monthlyRate > 0 {
				monthlyPayment = amount * monthlyRate / (1 - math.Pow(1+monthlyRate, float64(-numPayments)))
			} else if numPayments > 0 {
				monthlyPayment = amount / float64(numPayments)
			}
			annualPayment := monthlyPayment * 12
			totalPayment := monthlyPayment * float64(numPayments)

			data = map[string]interface{}{
				"Amount":         amount,
				"Rate":           fmt.Sprintf("%.3f", rate),
				"Years":          years,
				"MonthlyPayment": monthlyPayment,
				"AnnualPayment":  annualPayment,
				"TotalPayment":   totalPayment,
			}
		}
	}
	if err := mortgageTmpl.Execute(w, data); err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err == nil {
			income, _ := strconv.ParseFloat(r.FormValue("income"), 64)
			ownerSalary, _ := strconv.ParseFloat(r.FormValue("ownerSalary"), 64)
			price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
			years, _ := strconv.Atoi(r.FormValue("years"))
			rate, _ := strconv.ParseFloat(r.FormValue("rate"), 64)

			downPayment, _ := strconv.ParseFloat(r.FormValue("downPayment"), 64)
			loanAmount := price - downPayment

			monthlyRate := rate / 100 / 12
			numPayments := years * 12

			monthlyPayment := 0.0
			if monthlyRate > 0 {
				monthlyPayment = loanAmount * monthlyRate / (1 - math.Pow(1+monthlyRate, float64(-numPayments)))
			} else if numPayments > 0 {
				monthlyPayment = loanAmount / float64(numPayments)
			}
			totalLoanPayments := monthlyPayment * float64(numPayments)
			totalPaidPerYear := monthlyPayment * 12
			leftAfterOwner := income - ownerSalary
			leftAfterLoan := leftAfterOwner - totalPaidPerYear
			earningsRate := 0.0
			sbaIncome := income - ownerSalary
			if totalPaidPerYear > 0 {
				earningsRate = sbaIncome / totalPaidPerYear
			}

			data = map[string]interface{}{
				"Income":            income,
				"IncomeStr":         fmt.Sprintf("%.0f", income),
				"OwnerSalary":       ownerSalary,
				"OwnerSalaryStr":    fmt.Sprintf("%.0f", ownerSalary),
				"Price":             price,
				"PriceStr":          fmt.Sprintf("%.0f", price),
				"DownPayment":       downPayment,
				"DownPaymentStr":    fmt.Sprintf("%.0f", downPayment),
				"LoanAmount":        loanAmount,
				"Years":             years,
				"Rate":              rate,
				"RateStr":           fmt.Sprintf("%.3f", rate),
				"MonthlyPayment":    monthlyPayment,
				"TotalLoanPayments": totalLoanPayments,
				"TotalPaidPerYear":  totalPaidPerYear,
				"LeftAfterOwner":    leftAfterOwner,
				"LeftAfterLoan":     leftAfterLoan,
				"EarningsRate":      earningsRate,
			}
		}
	}
	if err := indexTmpl.Execute(w, data); err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}
