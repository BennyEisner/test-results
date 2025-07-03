package models

import (
	"database/sql"
	"encoding/xml"
)

type Project struct {
	ID   int64  `json:"id"  xml:"id"`
	Name string `json:"name" xml:"name"`
}
type ProjectXML struct {
	XMLName xml.Name `xml:"project"`
	Project Project  `xml:"project"`
}

// ScanFromRow scans a single row into a Project
func (p *Project) ScanFromRow(row *sql.Row) error {
	return row.Scan(&p.ID, &p.Name)
}

// ScanFromRows scans a row from Rows into a Project
func (p *Project) ScanFromRows(rows *sql.Rows) error {
	return rows.Scan(&p.ID, &p.Name)
}

// Insert inserts a new project into the database
func (p *Project) Insert(db *sql.DB) error {
	return db.QueryRow("INSERT INTO projects(name) VALUES($1) RETURNING id", p.Name).Scan(&p.ID)
}
