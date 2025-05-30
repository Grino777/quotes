package models

import "errors"

type Quote struct {
	Id     int32
	Author string
	Quote  string
}

func (q *Quote) Validate() error {
	if q.Author == "" {
		return errors.New("author field cannot be empty")
	}
	if q.Quote == "" {
		return errors.New("quote field cannot be empty")
	}
	return nil
}
