package controllers

import (
	"time"
	"webgorion/models"

	"github.com/astaxie/beego"
)

var (
	timeNow   = time.Now().Local()
	employees = []models.Employee{}
	doors     = []models.Door{}
)

type AcsController struct {
	beego.Controller
}

func (c *AcsController) SummaryReport() {
	getEmployee := c.Input().Get("employee")
	getDateRange := c.Input().Get("daterange")
	getDoor := c.Input().Get("door")

	c.Data["Title"] = "СКУД | Общий отчет"
	c.Data["FirstDate"] = timeNow.Format("02.01.2006")
	c.Data["LastDate"] = timeNow.AddDate(0, 0, 1).Format("02.01.2006")
	c.Data["Events"] = models.GetEvents(getDateRange, getDoor, getEmployee)
	c.Data["Doors"] = models.GetDoors()
	c.Data["DateRange"] = getDateRange
	c.Data["Employees"] = models.GetEmployees()

	c.Layout = "acs/layout.tpl"
	c.TplName = "acs/layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "acs/header.html"
	c.LayoutSections["Menu"] = "acs/menu.html"
	c.LayoutSections["Form"] = "acs/form_summary.html"
	c.LayoutSections["Table"] = "acs/table_summary.html"
	c.LayoutSections["Footer"] = "acs/footer.html"
}

func (c *AcsController) HoursReport() {
	getEmployee := c.Input().Get("employee")
	getDateRange := c.Input().Get("daterange")

	c.Data["Title"] = "СКУД | Отчет по отработанному времени"
	c.Data["FirstDate"] = timeNow.Format("02.01.2006")
	c.Data["LastDate"] = timeNow.AddDate(0, 0, 1).Format("02.01.2006")
	c.Data["Events"] = models.GetWorkHours(getDateRange, getEmployee)
	c.Data["DateRange"] = getDateRange
	c.Data["Employees"] = models.GetEmployees()

	c.Layout = "acs/layout.tpl"
	c.TplName = "acs/layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "acs/header.html"
	c.LayoutSections["Menu"] = "acs/menu.html"
	c.LayoutSections["Form"] = "acs/form_hours.html"
	c.LayoutSections["Table"] = "acs/table_hours.html"
	c.LayoutSections["Footer"] = "acs/footer.html"
}

func (c *AcsController) ViewHours() {
	c.Data["Title"] = "СКУД | Отчет по отработанному времени"
	c.Data["FirstDate"] = timeNow.Format("02.01.2006")
	c.Data["LastDate"] = timeNow.AddDate(0, 0, 1).Format("02.01.2006")
	c.Data["Employees"] = models.GetEmployees()

	c.TplName = "acs/layout.tpl"
	c.Layout = "acs/layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "acs/header.html"
	c.LayoutSections["Menu"] = "acs/menu.html"
	c.LayoutSections["Form"] = "acs/form_hours.html"
	c.LayoutSections["Footer"] = "acs/footer.html"
}

func (c *AcsController) ViewSummary() {
	c.Data["Title"] = "СКУД | Общий отчет"
	c.Data["FirstDate"] = timeNow.Format("02.01.2006")
	c.Data["LastDate"] = timeNow.AddDate(0, 0, 1).Format("02.01.2006")
	c.Data["Doors"] = models.GetDoors()
	c.Data["Employees"] = models.GetEmployees()

	c.TplName = "acs/layout.tpl"
	c.Layout = "acs/layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "acs/header.html"
	c.LayoutSections["Menu"] = "acs/menu.html"
	c.LayoutSections["Form"] = "acs/form_summary.html"
	c.LayoutSections["Footer"] = "acs/footer.html"
}
