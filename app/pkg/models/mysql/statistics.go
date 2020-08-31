// Copyright © Rob Burke inchworks.com, 2020.

// This file is part of PicInch.
//
// PicInch is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// PicInch is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with PicInch.  If not, see <https://www.gnu.org/licenses/>.

package mysql

// SQL operations on statistic table.

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"inchworks.com/gallery/pkg/usage"
)

const (
	statsDelete = `DELETE FROM statistic WHERE id = ?`

	statsInsert = `
		INSERT INTO statistic (event, category, count, start, period) VALUES (:event, :category, :count, :start, :period)`

	statsUpdate = `
		UPDATE statistic
		SET event=:event, category=:category, count=:count, start=:start, period=:period
		WHERE id=:id
	`
)

const (
	statsSelect = `SELECT * FROM statistic`
	statsOrderCategory  = ` ORDER BY category, start`
	statsOrderEvent  = ` ORDER BY event, start`
	statsOrderTime  = ` ORDER BY start DESC, count DESC, event ASC`

	statsWhereBefore = statsSelect + ` WHERE start < ? AND period = ?`
	statsWherePeriod = statsSelect + ` WHERE period = ?`
	statsWhereSingle = statsSelect + ` WHERE event = ? AND start = ? AND period = ?`

	statsBeforeByCategory = statsWhereBefore + statsOrderCategory
	statsBeforeByEvent = statsWhereBefore + statsOrderEvent
	statsBeforeByTime = statsWhereBefore + statsOrderTime

	statsDeleteIf = `DELETE FROM statistic WHERE start < ? AND period = ?`
)

type StatisticStore struct {
	GalleryId int64
	store
}

func NewStatisticStore(db *sqlx.DB, tx **sqlx.Tx, errorLog *log.Logger) *StatisticStore {

	return &StatisticStore{
		store: store{
			DBX:       db,
			ptx:       tx,
			errorLog:  errorLog,
			sqlDelete: statsDelete,
			sqlInsert: statsInsert,
			sqlUpdate: statsUpdate,
		},
	}
}

// Get statistics for specified period, ordered

func (st *StatisticStore) BeforeByTime(before time.Time, period int) []*usage.Statistic {

	var stats []*usage.Statistic

	if err := st.DBX.Select(&stats, statsBeforeByTime, before, period); err != nil {
		st.logError(err)
		return nil
	}
	return stats
}

// Before specified start time, ordered by category and time

func (st *StatisticStore) BeforeByCategory(before time.Time, period int) []*usage.Statistic {

	var stats []*usage.Statistic

	if err := st.DBX.Select(&stats, statsBeforeByCategory, before, period); err != nil {
		st.logError(err)
		return nil
	}
	return stats
}

// Before specified start time, ordered by event and time

func (st *StatisticStore) BeforeByEvent(before time.Time, period int) []*usage.Statistic {

	var stats []*usage.Statistic

	if err := st.DBX.Select(&stats, statsBeforeByEvent, before, period); err != nil {
		st.logError(err)
		return nil
	}
	return stats
}

// Delete old statistics
//
// Note that this is atypical as no other tables have specific functions for updates.

func (st *StatisticStore) DeleteIf(before time.Time, period int) error {

	tx := *st.ptx
	if tx == nil {
		panic("Transaction not begun")
	}

	if _, err := tx.Exec(statsDeleteIf, before, period); err != nil {
		return st.logError(err)
	}

	return nil
}

// Get single statistic, need not exist

func (st *StatisticStore) Get(event string, start time.Time, period int) *usage.Statistic {

	var s usage.Statistic

	if err := st.DBX.Get(&s, statsWhereSingle, event, start, period); err != nil {
		return nil
	}

	return &s
}

// Transaction for updates
// Note that this is atypical as other tables share an application-level transaction.

func (st *StatisticStore) Transaction() func() {

	// start transaction
	*st.ptx = st.DBX.MustBegin()

	return func() {

		// end transaction
		tx := *st.ptx
		tx.Commit()

		*st.ptx = nil
	}
}

// Insert or update statistic

func (st *StatisticStore) Update(s *usage.Statistic) error {

	return st.updateData(&s.Id, s)
}
