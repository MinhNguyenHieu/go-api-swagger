package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	Time  time.Time `json:"time"`
	Valid bool      `json:"valid"`
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nt.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &nt.Time)
	nt.Valid = err == nil
	return err
}

type NullString struct {
	String string `json:"string"`
	Valid  bool   `json:"valid"`
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ns.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = err == nil
	return err
}

func FromSQLNullTime(t sql.NullTime) NullTime {
	return NullTime{Time: t.Time, Valid: t.Valid}
}

func (nt NullTime) ToSQLNullTime() sql.NullTime {
	return sql.NullTime{Time: nt.Time, Valid: nt.Valid}
}

func FromSQLNullString(s sql.NullString) NullString {
	return NullString{String: s.String, Valid: s.Valid}
}

func (ns NullString) ToSQLNullString() sql.NullString {
	return sql.NullString{String: ns.String, Valid: ns.Valid}
}
