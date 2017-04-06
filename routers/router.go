package routers

import (
	"webgorion/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/acs/summary", &controllers.AcsController{}, "get:ViewSummary")
	beego.Router("/acs/hours", &controllers.AcsController{}, "get:ViewHours")
	beego.Router("/report_summary", &controllers.AcsController{}, "get:SummaryReport")
	beego.Router("/report_hours", &controllers.AcsController{}, "get:HoursReport")
}
