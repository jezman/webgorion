package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	db        *sql.DB
	err       error
	doors     = []Door{}
	employees = []Employee{}
	events    = []Events{}
)

type Door struct {
	ID   int64
	Name string
}

type Employee struct {
	ID        int64
	LastName  string
	FirstName string
	MidName   string
	Company   string
}

type Events struct {
	Employee  Employee
	FirstTime time.Time
	LastTime  time.Time
	Events    string
	WorkHours time.Duration
	Door      Door
}

type config struct {
	Server   string
	Database string
	User     string
	Password string
}

func init() {
	connect()
}

func readConfigFile() config {
	confFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Read configuration file error:", err)
	}
	defer confFile.Close()

	decoder := json.NewDecoder(confFile)
	conf := config{}
	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Println("JSON decode error:", err)
	}
	return conf
}
func connect() *sql.DB {
	conf := readConfigFile()
	dsn := "server=" + conf.Server +
		";user id=" + conf.User +
		";password=" + conf.Password +
		";database=" + conf.Database
	db, err = sql.Open("mssql", dsn)
	if err != nil {
		fmt.Println("Connection error", err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Database unreachable", err)
	}
	return db
}

func GetEmployees() []Employee {
	if len(employees) == 0 {
		rows, err := db.Query("SELECT ID, Name, FirstName, MidName FROM dbo.pList ORDER BY Name")
		if err != nil {
			fmt.Println("Query", err)
		}
		defer rows.Close()

		employee := Employee{}
		for rows.Next() {
			err := rows.Scan(&employee.ID,
				&employee.FirstName,
				&employee.LastName,
				&employee.MidName,
			)
			if err != nil {
				fmt.Println("Cols:", err)
			}
			employees = append(employees, employee)
		}
	}
	return employees
}

func GetDoors() []Door {
	if len(doors) == 0 {
		rows, err := db.Query("SELECT GIndex, Name FROM dbo.AcessPoint ORDER BY Name")
		if err != nil {
			fmt.Println("Query", err)
		}
		defer rows.Close()

		for rows.Next() {
			door := Door{}
			err = rows.Scan(&door.ID, &door.Name)
			if err != nil {
				fmt.Println("Cols:", err)
			}
			doors = append(doors, door)
		}
	}
	return doors
}

func GetEvents(daterange, door, employee string) (events []Events) {
	query := []string{"SELECT p.Name, p.FirstName, p.MidName, c.Name, TimeVal, e.Contents, a.Name ",
		"FROM dbo.pLogData l ",
		"JOIN dbo.pList p ON (p.ID = l.HozOrgan) ",
		"JOIN dbo.pCompany c ON (c.ID = p.Company) ",
		"JOIN dbo.Events e ON (e.Event = l.Event) ",
		"JOIN dbo.AcessPoint a ON (a.GIndex = l.DoorIndex) ",
		"WHERE TimeVal BETWEEN '", daterange[:10], "' AND '", daterange[13:], "'",
		" AND e.Event BETWEEN 26 AND 29",
		"ORDER BY TimeVal",
	}
	pName := " AND p.Id = '"
	doorIndex := "' AND DoorIndex = "
	orderBy := "' ORDER BY TimeVal"

	add := func(cmd ...string) {
		query = append(query[:len(query)-1], cmd...)
	}

	if door != "" && employee != "" {
		add(pName, employee, doorIndex, door, orderBy[1:])

	} else if employee != "" {
		add(pName, employee, orderBy)

	} else if door != "" {
		add(doorIndex[1:], door, orderBy[1:])
	}

	rows, err := db.Query(strings.Join(query, ""))
	if err != nil {
		fmt.Println("Query:", err)
	}

	event := Events{}

	for rows.Next() {
		err := rows.Scan(
			&event.Employee.FirstName,
			&event.Employee.LastName,
			&event.Employee.MidName,
			&event.Employee.Company,
			&event.FirstTime,
			&event.Events,
			&event.Door.Name,
		)
		if err != nil {
			fmt.Println("Cols:", err)
		}

		events = append(events, event)
	}
	return events
}

func GetWorkHours(dateRange, employee string) (events []Events) {
	query := []string{"SELECT p.Name, p.FirstName, p.MidName, c.Name, min(TimeVal), max(TimeVal) ",
		"FROM dbo.pLogData l ",
		"JOIN dbo.pList p ON (p.ID = l.HozOrgan) ",
		"JOIN dbo.pCompany c ON (c.ID = p.Company) ",
		"WHERE TimeVal BETWEEN '", dateRange[:10], "' AND '", dateRange[13:], "'",
		" AND p.Name = '", employee, "'",
		" GROUP BY p.Name, p.FirstName, p.MidName, c.Name, CONVERT(varchar(20), TimeVal, 104)",
	}
	if employee == "" {
		query = append(query[:9], query[12])
	}

	rows, err := db.Query(strings.Join(query, ""))
	if err != nil {
		fmt.Println("Query:", err)
	}

	event := Events{}

	for rows.Next() {
		err := rows.Scan(
			&event.Employee.FirstName,
			&event.Employee.LastName,
			&event.Employee.MidName,
			&event.Employee.Company,
			&event.FirstTime,
			&event.LastTime,
		)
		if err != nil {
			fmt.Println("Cols:", err)
		}
		event.WorkHours = event.LastTime.Sub(event.FirstTime)

		events = append(events, event)
	}
	return events
}
