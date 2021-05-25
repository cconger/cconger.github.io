package resume

import (
	"time"
)

type Resume struct {
	Owner    Owner
	Jobs     []Job     `toml:"job"`
	Legacys  []Job     `toml:"legacy"`
	Projects []Project `toml:"project"`
	Degrees  []Degree  `toml:"degree"`
}

type Owner struct {
	Firstname string
	Lastname  string
	Phone     string
	Email     string
}

type Job struct {
	Business     string
	Location     string
	Title        string
	Start        time.Time
	End          time.Time
	Languages    []string
	Technologies []string
	Projects     []string
}

func (j *Job) Tenure() string {
	if j.End.IsZero() {
		return j.Start.Format("Jan 2006" + " - Present")
	}
	return j.Start.Format("Jan 2006") + " - " + j.End.Format("Jan 2006")
}

type Project struct {
	Name         string
	Languages    []string
	Technologies []string
	URL          string
	Code         string
	Description  string
}

type Degree struct {
	Institution string
	Title       string
	Location    string
	Start       time.Time
	End         time.Time
	Accolades   []string
	GPA         GPA
}

func (d *Degree) Tenure() string {
	if d.End.IsZero() {
		return d.Start.Format("Jan 2006" + " - Present")
	}
	return d.Start.Format("Jan 2006") + " - " + d.End.Format("Jan 2006")
}

type GPA struct {
	Earned float32
	Max    float32
}
