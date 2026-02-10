package sqlite

import (
	"slices"
	"strings"

	"github.com/2manyvcos/paranal/server/data/schema"
)

func ServiceUserAlertChannelQuery(query *schema.ServiceUserAlertChannelQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.ServiceUserAlertChannelQuery)
	}
	var conditions []string
	if query.ErrorAlerts != nil {
		conditions = append(conditions, "useralertchannels.errorAlerts = ?")
		placeholders = append(placeholders, *query.ErrorAlerts)
	}
	if query.UptimeAlerts != nil {
		conditions = append(conditions, "useralertchannels.uptimeAlerts = ?")
		placeholders = append(placeholders, *query.UptimeAlerts)
	}
	if query.VersionAlerts != nil {
		conditions = append(conditions, "useralertchannels.versionAlerts = ?")
		placeholders = append(placeholders, *query.VersionAlerts)
	}
	if query.UserRole != nil {
		conditions = append(conditions, "users.role = ?")
		placeholders = append(placeholders, *query.UserRole)
	}
	if query.UserErrorAlerts != nil {
		conditions = append(conditions, "users.errorAlerts = ?")
		placeholders = append(placeholders, *query.UserErrorAlerts)
	}
	if query.ServiceHidden != nil {
		conditions = append(conditions, "IFNULL(serviceuserconfigs.hidden, FALSE) = ?")
		placeholders = append(placeholders, *query.ServiceHidden)
	}
	uptimeAlertConditions := make([]string, 0, 2)
	uptimeAlertPlaceholders := make([]any, 0, 2)
	if query.UserUptimeAlerts != nil {
		uptimeAlertConditions = append(uptimeAlertConditions, "users.uptimeAlerts = ?")
		uptimeAlertPlaceholders = append(uptimeAlertPlaceholders, *query.UserUptimeAlerts)
	}
	if query.ServiceUptimeAlerts != nil {
		uptimeAlertConditions = append(uptimeAlertConditions, "IFNULL(serviceuserconfigs.uptimeAlerts, FALSE) = ?")
		uptimeAlertPlaceholders = append(uptimeAlertPlaceholders, *query.ServiceUptimeAlerts)
	}
	if len(uptimeAlertConditions) == 0 {
		uptimeAlertConditions = append(uptimeAlertConditions, "1 = 1")
	}
	conditions = append(conditions, "("+strings.Join(uptimeAlertConditions, " OR ")+")")
	placeholders = append(placeholders, uptimeAlertPlaceholders...)
	versionConditions := make([]string, 0, 2)
	versionPlaceholders := make([]any, 0, 2)
	if query.UserVersionAlerts != nil {
		versionConditions = append(versionConditions, "users.versionAlerts = ?")
		versionPlaceholders = append(versionPlaceholders, *query.UserVersionAlerts)
	}
	if query.ServiceVersionAlerts != nil {
		versionConditions = append(versionConditions, "IFNULL(serviceuserconfigs.versionAlerts, FALSE) = ?")
		versionPlaceholders = append(versionPlaceholders, *query.ServiceVersionAlerts)
	}
	if len(versionConditions) == 0 {
		versionConditions = append(versionConditions, "1 = 1")
	}
	conditions = append(conditions, "("+strings.Join(versionConditions, " OR ")+")")
	placeholders = append(placeholders, versionPlaceholders...)
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListServiceUserAlertChannels(serviceID string, query *schema.ServiceUserAlertChannelQuery) ([]schema.ServiceUserAlertChannel, error) {
	where, wherePlaceholders := ServiceUserAlertChannelQuery(query)
	rows, err := i.Query(
		`
      SELECT useralertchannels.id, useralertchannels.userName, useralertchannels.url, useralertchannels.errorAlerts, useralertchannels.uptimeAlerts, useralertchannels.versionAlerts
      FROM useralertchannels
      INNER JOIN users
      ON useralertchannels.userName = users.name
      LEFT JOIN serviceuserconfigs
      ON serviceuserconfigs.userName = users.name AND serviceuserconfigs.serviceID = ?
      WHERE `+where+`
    `,
		slices.Concat(
			[]any{serviceID},
			wherePlaceholders,
		)...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []schema.ServiceUserAlertChannel
	for rows.Next() {
		var record schema.ServiceUserAlertChannel
		if err := rows.Scan(&record.ID, &record.UserName, &record.URL, &record.ErrorAlerts, &record.UptimeAlerts, &record.VersionAlerts); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}
