package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	db        *sql.DB
	rows      *sql.Rows
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

type Config struct {
	Server        string
	Database      string
	User          string
	Password      string
	Smtp2goApiKey string
}

func connect() *sql.DB {
	dsn := "server=" + Conf.Server +
		";user id=" + Conf.User +
		";password=" + Conf.Password +
		";database=" + Conf.Database
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
	connect()
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
	defer db.Close()
	return employees
}

func GetDoors() []Door {
	connect()
	if len(doors) == 0 {
		rows, err := db.Query("SELECT GIndex, Name FROM dbo.AcessPoint ORDER BY Name")
		if err != nil {
			fmt.Println("Query:", err)
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
	defer db.Close()
	return doors
}

func GetEvents(dateRange, door string, employee []string) (events []Events) {
	connect()
	query := []string{`SELECT p.Name, p.FirstName, p.MidName, c.Name, TimeVal, e.Contents, a.Name
		FROM dbo.pLogData l
		JOIN dbo.pList p ON (p.ID = l.HozOrgan)
		JOIN dbo.pCompany c ON (c.ID = p.Company)
		JOIN dbo.Events e ON (e.Event = l.Event)
		JOIN dbo.AcessPoint a ON (a.GIndex = l.DoorIndex)
		WHERE TimeVal BETWEEN '`, dateRange[:10], `' AND '`, dateRange[13:],
		`' AND e.Event BETWEEN 26 AND 29 `,
	}

	// add params to query
	updateRows := func(cmd ...interface{}) {
		rows, err = db.Query(strings.Join(query, ``), cmd...)
		if err != nil {
			fmt.Println("Query:", err)
		}
	}

	pNameOne := ` AND p.Id in (?) `
	pNameTwo := ` AND p.Id in (?, ?) `
	pNameThree := ` AND p.Id in (?, ?, ?) `
	pNameFour := ` AND p.Id in (?, ?, ?, ?) `
	doorIndex := ` AND DoorIndex = ?`

	if door != "" && len(employee) != 0 {
		switch len(employee) {
		case 1:
			query = append(query, pNameOne, doorIndex)
			updateRows(employee[0], door)
		case 2:
			query = append(query, pNameTwo, doorIndex)
			updateRows(employee[0], employee[1], door)
		case 3:
			query = append(query, pNameThree, doorIndex)
			updateRows(employee[0], employee[1], employee[2], door)
		case 4:
			query = append(query, pNameFour, doorIndex)
			updateRows(employee[0], employee[1], employee[2], employee[3], door)
		}

	} else if len(employee) != 0 {
		switch len(employee) {
		case 1:
			query = append(query, pNameOne)
			updateRows(employee[0])
		case 2:
			query = append(query, pNameTwo)
			updateRows(employee[0], employee[1])
		case 3:
			query = append(query, pNameThree)
			updateRows(employee[0], employee[1], employee[2])
		case 4:
			query = append(query, pNameFour)
			updateRows(employee[0], employee[1], employee[2], employee[3])
		}

	} else if door != "" {
		query = append(query, doorIndex)
		updateRows(door)

	} else {
		updateRows()
	}

	event := Events{}

	defer rows.Close()
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
	defer db.Close()
	return events
}

func GetWorkHours(dateRange string, employee []string) (events []Events) {
	connect()
	query := []string{`SELECT p.Name, p.FirstName, p.MidName, c.Name, min(TimeVal), max(TimeVal)
		FROM dbo.pLogData l
		JOIN dbo.pList p ON (p.ID = l.HozOrgan)
		JOIN dbo.pCompany c ON (c.ID = p.Company) `,
		`WHERE TimeVal BETWEEN '`, dateRange[:10], `' AND '`, dateRange[13:], `'`,
		`AND p.Id in(?)`,
		`GROUP BY p.Name, p.FirstName, p.MidName, c.Name, CONVERT(varchar(20), TimeVal, 104)`,
	}

	// add params to query
	updateRows := func(cmd ...interface{}) {
		rows, err = db.Query(strings.Join(query, ``), cmd...)
		if err != nil {
			fmt.Println("Query:", err)
		}
	}

	if len(employee) == 0 {
		query = append(query[:9], query[10])
	}

	switch len(employee) {
	case 1:
		query[6] = ` AND p.Id in (?) `
		updateRows(employee[0])
	case 2:
		query[6] = ` AND p.Id in (?, ?) `
		updateRows(employee[0], employee[1])
	case 3:
		query[6] = ` AND p.Id in (?, ?, ?) `
		updateRows(employee[0], employee[1], employee[2])
	case 4:
		query[6] = ` AND p.Id in(?, ?, ?, ?) `
		updateRows(employee[0], employee[1], employee[2], employee[3])
	}

	event := Events{}

	defer rows.Close()
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
	defer db.Close()
	return events
}
