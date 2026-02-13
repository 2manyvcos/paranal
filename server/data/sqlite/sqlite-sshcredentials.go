package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/2manyvcos/paranal/server/schema"
)

func init() {
	setups = append(setups, func(i *impl) error {
		_, err := i.Exec(`
      CREATE TABLE IF NOT EXISTS sshcredentials (
        name TEXT PRIMARY KEY NOT NULL,
        user TEXT NOT NULL,
        password TEXT NOT NULL,
        privateKey TEXT NOT NULL
      )
    `)
		if err != nil {
			return fmt.Errorf("creating table \"sshcredentials\" failed - %s", err)
		}
		return nil
	})
}

func SSHCredentialQuery(query *schema.SSHCredentialQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.SSHCredentialQuery)
	}
	var conditions []string
	if query.Name != nil {
		conditions = append(conditions, "sshcredentials.name = ?")
		placeholders = append(placeholders, *query.Name)
	}
	if query.User != nil {
		conditions = append(conditions, "sshcredentials.user = ?")
		placeholders = append(placeholders, *query.User)
	}
	if query.HasPassword != nil {
		if *query.HasPassword {
			conditions = append(conditions, "IFNULL(sshcredentials.password, '') <> ''")
		} else {
			conditions = append(conditions, "IFNULL(sshcredentials.password, '') = ''")
		}
	}
	if query.HasPrivateKey != nil {
		if *query.HasPrivateKey {
			conditions = append(conditions, "IFNULL(sshcredentials.privateKey, '') <> ''")
		} else {
			conditions = append(conditions, "IFNULL(sshcredentials.privateKey, '') = ''")
		}
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListSSHCredentials(query *schema.SSHCredentialQuery) ([]schema.SSHCredential, error) {
	where, wherePlaceholders := SSHCredentialQuery(query)
	rows, err := i.Query(
		`
      SELECT name, user, password, privateKey
      FROM sshcredentials
      WHERE `+where+`
      ORDER BY name
    `,
		wherePlaceholders...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.SSHCredential
	for rows.Next() {
		var record schema.SSHCredential
		if err := rows.Scan(&record.Name, &record.User, &record.Password, &record.PrivateKey); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *impl) GetSSHCredential(query schema.SSHCredentialQuery) (schema.SSHCredential, error) {
	where, wherePlaceholders := SSHCredentialQuery(&query)
	var record schema.SSHCredential
	err := i.QueryRow(
		`
      SELECT name, user, password, privateKey
      FROM sshcredentials
      WHERE `+where+`
    `,
		wherePlaceholders...,
	).Scan(&record.Name, &record.User, &record.Password, &record.PrivateKey)
	if errors.Is(err, sql.ErrNoRows) {
		return schema.SSHCredential{}, schema.ErrNotFound
	}
	return record, err
}

func (i *impl) CreateSSHCredential(record schema.SSHCredential) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO sshcredentials (name, user, password, privateKey)
      VALUES (?, ?, ?, ?)
    `,
		&record.Name, &record.User, &record.Password, &record.PrivateKey,
	)
	return requireNoConflict(err)
}

func (i *impl) UpdateSSHCredentials(query schema.SSHCredentialQuery, record schema.SSHCredential) error {
	if err := record.Valid(); err != nil {
		return err
	}
	where, wherePlaceholders := SSHCredentialQuery(&query)
	result, err := i.Exec(
		`
      UPDATE sshcredentials
      SET
        name = ?,
        user = ?,
        password = ?
        privateKey = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{&record.Name, &record.User, &record.Password, &record.PrivateKey},
			wherePlaceholders,
		)...,
	)
	return requireFound(result, err)
}

func (i *impl) DeleteSSHCredentials(query schema.SSHCredentialQuery) error {
	where, wherePlaceholders := SSHCredentialQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM sshcredentials
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
