package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/jezman/webgorion/models"

	"github.com/astaxie/beego"
)

var (
	employees = []models.Employee{}
	doors     = []models.Door{}
)

type AcsController struct {
	beego.Controller
}

func FirstDate() string {
	timeNow := time.Now().Local()
	now := timeNow.Format("02.01.2006")
	return now
}

func LastDate() string {
	timeNow := time.Now().Local()
	nowAdd := timeNow.AddDate(0, 0, 1).Format("02.01.2006")
	return nowAdd
}

// Validation door input
func ValidationDoor(getDoor string) bool {
	matchDoor, err := regexp.MatchString(`^\d{0,3}$`, getDoor)
	if err != nil {
		fmt.Println("Regexp err:", err)
	}
	if (!matchDoor && len(getDoor) != 0) || len(getDoor) > 3 {
		return false
	}
	return true
}

// Validation employee input
func ValidationEmployees(getEmployees []string) bool {
	for i, _ := range getEmployees {
		matchEmployee, _ := regexp.MatchString(`^\d{0,3}$`, getEmployees[i])
		if (!matchEmployee && len(getEmployees[i]) != 0) || len(getEmployees[i]) > 3 || len(getEmployees) > 4 {
			return false
		}
	}
	return true
}

// Validation date input
func ValidationDate(getDateRange string) bool {
	matchDate, err := regexp.MatchString(`(0[1-9]|[12][0-9]|3[01])[- ..](0[1-9]|1[012])[- ..][201]\d\d\d\s[- +.]\s(0[1-9]|[12][0-9]|3[01])[- ..](0[1-9]|1[012])[- ..][201]\d\d\d`, getDateRange)
	if err != nil {
		fmt.Println("Regexp err:", err)
	}
	if !matchDate {
		return false
	}
	return true
}

func (c *AcsController) SummaryReport() {
	getEmployees := c.GetStrings("employee")
	getDateRange := c.Input().Get("daterange")
	getDoor := c.Input().Get("door")

	if !ValidationDoor(getDoor) || !ValidationDate(getDateRange) || !ValidationEmployees(getEmployees) {
		c.Redirect("/", 302)
	}

	c.Data["Title"] = "СКУД | Общий отчет"
	c.Data["FirstDate"] = FirstDate()
	c.Data["LastDate"] = LastDate()
	c.Data["Events"] = models.GetEvents(getDateRange, getDoor, getEmployees)
	c.Data["Doors"] = models.GetDoors()
	c.Data["DateRange"] = getDateRange
	c.Data["Employees"] = models.GetEmployees()

	c.Layout = "acs/layout.tpl"
	c.TplName = "acs/layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "header.html"
	c.LayoutSections["Menu"] = "menu.html"
	c.LayoutSections["Form"] = "acs/form_summary.html"
	c.LayoutSections["Table"] = "acs/table_summary.html"
	c.LayoutSections["Footer"] = "footer.html"
}

func (c *AcsController) HoursReport() {
	getEmployees := c.GetStrings("employee")
	getDateRange := c.Input().Get("daterange")

	if !ValidationDate(getDateRange) || !ValidationEmployees(getEmployees) {
		c.Redirect("/", 302)
	}

	c.Data["Title"] = "СКУД | Отчет по отработанному времени"
	c.Data["FirstDate"] = FirstDate()
	c.Data["LastDate"] = LastDate()
	c.Data["Events"] = models.GetWorkHours(getDateRange, getEmployees)
	c.Data["DateRange"] = getDateRange
	c.Data["Employees"] = models.GetEmployees()

	c.Layout = "acs/layout.tpl"
	c.TplName = "acs/layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "header.html"
	c.LayoutSections["Menu"] = "menu.html"
	c.LayoutSections["Form"] = "acs/form_hours.html"
	c.LayoutSections["Table"] = "acs/table_hours.html"
	c.LayoutSections["Footer"] = "footer.html"
}

func (c *AcsController) ViewHours() {
	c.Data["Title"] = "СКУД | Отчет по отработанному времени"
	c.Data["FirstDate"] = FirstDate()
	c.Data["LastDate"] = LastDate()
	c.Data["Employees"] = models.GetEmployees()

	c.TplName = "acs/layout.tpl"
	c.Layout = "acs/layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "header.html"
	c.LayoutSections["Menu"] = "menu.html"
	c.LayoutSections["Form"] = "acs/form_hours.html"
	c.LayoutSections["Footer"] = "footer.html"
}

func testHandler(w http.ResponseWriter, r http.Request) {
	fmt.Fprintln(w, "temp")
}

func (c *AcsController) ViewSummary() {
	c.Data["Title"] = "СКУД | Общий отчет"
	c.Data["FirstDate"] = FirstDate()
	c.Data["LastDate"] = LastDate()
	c.Data["Doors"] = models.GetDoors()
	c.Data["Employees"] = models.GetEmployees()

	c.TplName = "acs/layout.tpl"
	c.Layout = "acs/layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "header.html"
	c.LayoutSections["Menu"] = "menu.html"
	c.LayoutSections["Form"] = "acs/form_summary.html"
	c.LayoutSections["Footer"] = "footer.html"
}
