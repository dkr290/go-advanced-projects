package handlers

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

var templates *template.Template

func init() {
	funcMap := template.FuncMap{
		"formatMoney": func(v float64) string {
			return fmt.Sprintf("%.2f", v)
		},
		"formatMoneyInt": func(v float64) string {
			return fmt.Sprintf("%.0f", v)
		},
		"seq": func(start, end int) []int {
			s := make([]int, 0, end-start+1)
			for i := start; i <= end; i++ {
				s = append(s, i)
			}
			return s
		},
		"toString": func(v int) string {
			return strconv.Itoa(v)
		},
		"compoundFreqLabel": func(v string) string {
			labels := map[string]string{
				"1":   "Annually",
				"2":   "Semi-Annually",
				"4":   "Quarterly",
				"12":  "Monthly",
				"365": "Daily",
			}
			if l, ok := labels[v]; ok {
				return l
			}
			return v
		},
	}
	var err error
	templates, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
}

type AmortisationRow struct {
	Year          int
	StartBalance  float64
	PrincipalPaid float64
	InterestPaid  float64
	EndBalance    float64
}

type CompoundInterestRow struct {
	Year           int
	StartBalance   float64
	Contributions  float64
	InterestEarned float64
	EndBalance     float64
}

type CompoundData struct {
	Title                string
	Currency             string
	InitialInvestment    string
	InterestRate         string
	CompoundingFrequency string
	InvestmentTerm       string
	MonthlyContribution  string
	HasResults           bool
	FinalBalance         float64
	TotalContributions   float64
	TotalInterestEarned  float64
	GrowthPercentage     float64
	Schedule             []CompoundInterestRow
	Error                string
}

type MortgageData struct {
	Title          string
	Currency       string
	LoanAmount     string
	DownPayment    string
	InterestRate   string
	LoanTerm       string
	HasResults     bool
	FinalLoanAmt   float64
	MonthlyPayment float64
	TotalRepayment float64
	TotalInterest  float64
	Schedule       []AmortisationRow
	Error          string
}

func Index(c fiber.Ctx) error {
	data := map[string]string{
		"Title":   "Fincalc",
		"Message": "Welcome to Fincalc",
	}
	c.Set("Content-Type", "text/html")
	return templates.ExecuteTemplate(c.Response().BodyWriter(), "index.html", data)
}

func Test(c fiber.Ctx) error {
	data := map[string]string{
		"Title": "Test Page",
	}
	c.Set("Content-Type", "text/html")
	return templates.ExecuteTemplate(c.Response().BodyWriter(), "test.html", data)
}

func Mortgage(c fiber.Ctx) error {
	data := MortgageData{
		Title:    "Mortgage Repayments Calculator",
		Currency: "USD",
	}
	c.Set("Content-Type", "text/html")
	return templates.ExecuteTemplate(c.Response().BodyWriter(), "mortgage.html", data)
}

func MortgageCalculate(c fiber.Ctx) error {
	data := MortgageData{
		Title:        "Mortgage Repayments Calculator",
		Currency:     c.FormValue("currency"),
		LoanAmount:   c.FormValue("loan_amount"),
		DownPayment:  c.FormValue("down_payment"),
		InterestRate: c.FormValue("interest_rate"),
		LoanTerm:     c.FormValue("loan_term"),
	}

	if data.Currency == "" {
		data.Currency = "USD"
	}

	loan, err := strconv.ParseFloat(data.LoanAmount, 64)
	if err != nil || loan <= 0 {
		data.Error = "Loan amount must be a positive number."
		c.Set("Content-Type", "text/html")
		return templates.ExecuteTemplate(c.Response().BodyWriter(), "mortgage.html", data)
	}

	down, err := strconv.ParseFloat(data.DownPayment, 64)
	if err != nil || down < 0 {
		data.Error = "Down payment must be zero or a positive number."
		c.Set("Content-Type", "text/html")
		return templates.ExecuteTemplate(c.Response().BodyWriter(), "mortgage.html", data)
	}

	if down >= loan {
		data.Error = "Down payment must be less than the loan amount."
		c.Set("Content-Type", "text/html")
		return templates.ExecuteTemplate(c.Response().BodyWriter(), "mortgage.html", data)
	}

	rate, err := strconv.ParseFloat(data.InterestRate, 64)
	if err != nil || rate <= 0 || rate > 30 {
		data.Error = "Interest rate must be between 0 and 30%."
		c.Set("Content-Type", "text/html")
		return templates.ExecuteTemplate(c.Response().BodyWriter(), "mortgage.html", data)
	}

	term, err := strconv.Atoi(data.LoanTerm)
	if err != nil || term < 5 || term > 40 {
		data.Error = "Loan term must be between 5 and 40 years."
		c.Set("Content-Type", "text/html")
		return templates.ExecuteTemplate(c.Response().BodyWriter(), "mortgage.html", data)
	}

	finalLoan := loan - down
	monthlyRate := rate / 100.0 / 12.0
	numPayments := term * 12

	var monthlyPayment float64
	if monthlyRate == 0 {
		monthlyPayment = finalLoan / float64(numPayments)
	} else {
		monthlyPayment = finalLoan * (monthlyRate * math.Pow(1+monthlyRate, float64(numPayments))) /
			(math.Pow(1+monthlyRate, float64(numPayments)) - 1)
	}

	totalRepayment := monthlyPayment * float64(numPayments)
	totalInterest := totalRepayment - finalLoan

	data.FinalLoanAmt = finalLoan
	data.MonthlyPayment = monthlyPayment
	data.TotalRepayment = totalRepayment
	data.TotalInterest = totalInterest
	data.HasResults = true

	schedule := make([]AmortisationRow, 0, term)
	balance := finalLoan
	for year := 1; year <= term; year++ {
		startBalance := balance
		yearInterest := 0.0
		yearPrincipal := 0.0

		for month := 0; month < 12; month++ {
			interestPayment := balance * monthlyRate
			principalPayment := monthlyPayment - interestPayment
			if balance < principalPayment {
				principalPayment = balance
			}
			balance -= principalPayment
			yearInterest += interestPayment
			yearPrincipal += principalPayment
		}

		if balance < 0.01 {
			balance = 0
		}

		schedule = append(schedule, AmortisationRow{
			Year:          year,
			StartBalance:  startBalance,
			PrincipalPaid: yearPrincipal,
			InterestPaid:  yearInterest,
			EndBalance:    balance,
		})
	}

	data.Schedule = schedule

	c.Set("Content-Type", "text/html")
	return templates.ExecuteTemplate(c.Response().BodyWriter(), "mortgage.html", data)
}

func Compound(c fiber.Ctx) error {
	data := CompoundData{
		Title:                "Compound Interest Calculator",
		Currency:             "USD",
		CompoundingFrequency: "12",
	}
	c.Set("Content-Type", "text/html")
	return templates.ExecuteTemplate(c.Response().BodyWriter(), "compound.html", data)
}

func CompoundCalculate(c fiber.Ctx) error {
	data := CompoundData{
		Title:                "Compound Interest Calculator",
		Currency:             c.FormValue("currency"),
		InitialInvestment:    c.FormValue("initial_investment"),
		InterestRate:         c.FormValue("interest_rate"),
		CompoundingFrequency: c.FormValue("compounding_frequency"),
		InvestmentTerm:       c.FormValue("investment_term"),
		MonthlyContribution:  c.FormValue("monthly_contribution"),
	}

	if data.Currency == "" {
		data.Currency = "USD"
	}
	if data.CompoundingFrequency == "" {
		data.CompoundingFrequency = "12"
	}

	initial, err := strconv.ParseFloat(data.InitialInvestment, 64)
	if err != nil || initial < 0 {
		data.Error = "Initial investment must be zero or a positive number."
		c.Set("Content-Type", "text/html")
		return templates.ExecuteTemplate(c.Response().BodyWriter(), "compound.html", data)
	}

	rate, err := strconv.ParseFloat(data.InterestRate, 64)
	if err != nil || rate <= 0 || rate > 30 {
		data.Error = "Interest rate must be between 0 and 30%."
		c.Set("Content-Type", "text/html")
		return templates.ExecuteTemplate(c.Response().BodyWriter(), "compound.html", data)
	}

	compoundFreq, err := strconv.Atoi(data.CompoundingFrequency)
	if err != nil || (compoundFreq != 1 && compoundFreq != 2 && compoundFreq != 4 && compoundFreq != 12 && compoundFreq != 365) {
		data.Error = "Compounding frequency must be annually, semi-annually, quarterly, monthly, or daily."
		c.Set("Content-Type", "text/html")
		return templates.ExecuteTemplate(c.Response().BodyWriter(), "compound.html", data)
	}

	term, err := strconv.Atoi(data.InvestmentTerm)
	if err != nil || term < 1 || term > 50 {
		data.Error = "Investment term must be between 1 and 50 years."
		c.Set("Content-Type", "text/html")
		return templates.ExecuteTemplate(c.Response().BodyWriter(), "compound.html", data)
	}

	monthlyContrib := 0.0
	if data.MonthlyContribution != "" {
		monthlyContrib, err = strconv.ParseFloat(data.MonthlyContribution, 64)
		if err != nil || monthlyContrib < 0 {
			data.Error = "Monthly contribution must be zero or a positive number."
			c.Set("Content-Type", "text/html")
			return templates.ExecuteTemplate(c.Response().BodyWriter(), "compound.html", data)
		}
	}

	compoundIntervalMonths := 12 / compoundFreq

	balance := initial
	totalContributions := 0.0
	schedule := make([]CompoundInterestRow, 0, term)

	for year := 1; year <= term; year++ {
		startBalance := balance
		yearContributions := 0.0
		yearInterest := 0.0

		for month := 1; month <= 12; month++ {
			balance += monthlyContrib
			yearContributions += monthlyContrib
			totalContributions += monthlyContrib

			if month%compoundIntervalMonths == 0 {
				interest := balance * (rate / 100.0) / float64(compoundFreq)
				balance += interest
				yearInterest += interest
			}
		}

		schedule = append(schedule, CompoundInterestRow{
			Year:           year,
			StartBalance:   startBalance,
			Contributions:  yearContributions,
			InterestEarned: yearInterest,
			EndBalance:     balance,
		})
	}

	data.FinalBalance = balance
	data.TotalContributions = totalContributions
	data.TotalInterestEarned = balance - initial - totalContributions
	if initial > 0 {
		data.GrowthPercentage = ((balance - initial) / initial) * 100
	}
	data.HasResults = true
	data.Schedule = schedule

	c.Set("Content-Type", "text/html")
	return templates.ExecuteTemplate(c.Response().BodyWriter(), "compound.html", data)
}
