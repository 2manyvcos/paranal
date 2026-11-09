package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (db *sqliteDB) CleanupCollectorStates(collectors []string) (modifiedStates []StateID, err error) {
	db.Lock()
	defer db.Unlock()

	var rows *sql.Rows
	if len(collectors) == 0 {
		rows, err = db.Query(
			`
        DELETE FROM states
        WHERE collector != ''
        RETURNING registry, id
      `,
		)
	} else {
		args := make([]any, len(collectors))
		for i, v := range collectors {
			args[i] = v
		}
		rows, err = db.Query(
			`
        DELETE FROM states
        WHERE
          collector != ''
          AND collector NOT IN (`+strings.TrimSuffix(strings.Repeat("?, ", len(collectors)), ", ")+`)
        RETURNING registry, id
      `,
			args...,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var registry, id string
		if err := rows.Scan(&registry, &id); err != nil {
			return nil, err
		}
		modifiedStates = append(modifiedStates, StateID{Registry: registry, ID: id})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modifiedStates, nil
}

func (db *sqliteDB) PutCollectorStates(collector string, states []State) (modifiedStates []StateID, err error) {
	db.Lock()
	defer db.Unlock()

	if collector == "" {
		return nil, fmt.Errorf("missing collector")
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	revs := make(map[StateID]string)

	rows, err := tx.Query(
		`
      SELECT registry, id, rev
      FROM states
      WHERE collector = ?
    `,
		collector,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var registry, id, rev string
		if err := rows.Scan(&registry, &id, &rev); err != nil {
			return nil, err
		}
		revs[StateID{Registry: registry, ID: id}] = rev
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	modifiedStates = make([]StateID, 0, len(revs)+len(states))
	seen := make(map[StateID]struct{}, len(states))

	stmt, err := tx.Prepare(
		`
      INSERT INTO states (registry, id, rev, collector, data)
      VALUES (?, ?, ?, ?, jsonb(?))
      ON CONFLICT (registry, id) DO UPDATE
        SET
          rev = EXCLUDED.rev,
          collector = EXCLUDED.collector,
          data = EXCLUDED.data
        WHERE
          states.rev IS DISTINCT FROM EXCLUDED.rev
          OR states.collector IS DISTINCT FROM EXCLUDED.collector
          OR states.data IS DISTINCT FROM EXCLUDED.data
    `,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	for _, state := range states {
		id := StateID{Registry: state.Registry, ID: state.ID}
		seen[id] = struct{}{}
		if revs[id] == state.Rev {
			continue
		}
		_, err := stmt.Exec(state.Registry, state.ID, state.Rev, collector, string(state.Data))
		if err != nil {
			return nil, err
		}
		modifiedStates = append(modifiedStates, id)
	}

	stmt, err = tx.Prepare(
		`
      DELETE FROM states
      WHERE
        registry = ?
        AND id = ?
        AND collector = ?
    `,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	for id := range revs {
		if _, ok := seen[id]; ok {
			continue
		}
		_, err = stmt.Exec(id.Registry, id.ID, collector)
		if err != nil {
			return nil, err
		}
		modifiedStates = append(modifiedStates, id)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return modifiedStates, nil
}

// TODO: use it like `SELECT json(data) FROM states …`
