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
	userAlertChannelConditions, userAlertChannelPlaceholders := UserAlertChannelQuery(&query.UserAlertChannelQuery)
	conditions = append(conditions, userAlertChannelConditions)
	placeholders = append(placeholders, userAlertChannelPlaceholders...)
	userConditions, userPlaceholders := UserQuery(&query.UserQuery)
	conditions = append(conditions, userConditions)
	placeholders = append(placeholders, userPlaceholders...)
	serviceConditions, servicePlaceholders := ServiceUserConfigQuery(&query.ServiceUserConfigQuery)
	conditions = append(conditions, serviceConditions)
	placeholders = append(placeholders, servicePlaceholders...)
	if query.ErrorAlerts != nil {
		if *query.ErrorAlerts {
			conditions = append(conditions, "IFNULL(useralertchannels.errorAlerts, FALSE) = TRUE AND IFNULL(users.errorAlerts, FALSE) = TRUE")
		} else {
			conditions = append(conditions, "(IFNULL(useralertchannels.errorAlerts, FALSE) = FALSE OR IFNULL(users.errorAlerts, FALSE) = FALSE)")
		}
	}
	if query.UptimeAlerts != nil {
		if *query.UptimeAlerts {
			conditions = append(conditions, "IFNULL(useralertchannels.uptimeAlerts, FALSE) = TRUE AND (IFNULL(users.uptimeAlerts, FALSE) = TRUE OR IFNULL(serviceuserconfigs.uptimeAlerts, FALSE) = TRUE)")
		} else {
			conditions = append(conditions, "(IFNULL(useralertchannels.uptimeAlerts, FALSE) = FALSE OR (IFNULL(users.uptimeAlerts, FALSE) = FALSE AND IFNULL(serviceuserconfigs.uptimeAlerts, FALSE) = FALSE))")
		}
	}
	if query.VersionAlerts != nil {
		if *query.VersionAlerts {
			conditions = append(conditions, "IFNULL(useralertchannels.versionAlerts, FALSE) = TRUE AND (IFNULL(users.versionAlerts, FALSE) = TRUE OR IFNULL(serviceuserconfigs.versionAlerts, FALSE) = TRUE)")
		} else {
			conditions = append(conditions, "(IFNULL(useralertchannels.versionAlerts, FALSE) = FALSE OR (IFNULL(users.versionAlerts, FALSE) = FALSE AND IFNULL(serviceuserconfigs.versionAlerts, FALSE) = FALSE))")
		}
	}
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
