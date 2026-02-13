package sqlite

import (
	"slices"
	"strings"

	"github.com/2manyvcos/paranal/server/schema"
)

func ServiceAlertChannelQuery(query *schema.ServiceAlertChannelQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.ServiceAlertChannelQuery)
	}
	var conditions []string
	userAlertChannelConditions, userAlertChannelPlaceholders := UserAlertChannelQuery(&query.UserAlertChannelQuery)
	conditions = append(conditions, userAlertChannelConditions)
	placeholders = append(placeholders, userAlertChannelPlaceholders...)
	userConditions, userPlaceholders := UserQuery(&query.UserQuery)
	conditions = append(conditions, userConditions)
	placeholders = append(placeholders, userPlaceholders...)
	serviceConfigConditions, serviceConfigPlaceholders := ServiceConfigQuery(&query.ServiceConfigQuery)
	conditions = append(conditions, serviceConfigConditions)
	placeholders = append(placeholders, serviceConfigPlaceholders...)
	if query.ErrorAlerts != nil {
		if *query.ErrorAlerts {
			conditions = append(conditions, "IFNULL(useralertchannels.errorAlerts, FALSE) = TRUE AND IFNULL(users.errorAlerts, FALSE) = TRUE")
		} else {
			conditions = append(conditions, "(IFNULL(useralertchannels.errorAlerts, FALSE) = FALSE OR IFNULL(users.errorAlerts, FALSE) = FALSE)")
		}
	}
	if query.UptimeAlerts != nil {
		if *query.UptimeAlerts {
			conditions = append(conditions, "IFNULL(useralertchannels.uptimeAlerts, FALSE) = TRUE AND (IFNULL(users.uptimeAlerts, FALSE) = TRUE OR IFNULL(serviceconfigs.uptimeAlerts, FALSE) = TRUE)")
		} else {
			conditions = append(conditions, "(IFNULL(useralertchannels.uptimeAlerts, FALSE) = FALSE OR (IFNULL(users.uptimeAlerts, FALSE) = FALSE AND IFNULL(serviceconfigs.uptimeAlerts, FALSE) = FALSE))")
		}
	}
	if query.VersionAlerts != nil {
		if *query.VersionAlerts {
			conditions = append(conditions, "IFNULL(useralertchannels.versionAlerts, FALSE) = TRUE AND (IFNULL(users.versionAlerts, FALSE) = TRUE OR IFNULL(serviceconfigs.versionAlerts, FALSE) = TRUE)")
		} else {
			conditions = append(conditions, "(IFNULL(useralertchannels.versionAlerts, FALSE) = FALSE OR (IFNULL(users.versionAlerts, FALSE) = FALSE AND IFNULL(serviceconfigs.versionAlerts, FALSE) = FALSE))")
		}
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) ListServiceAlertChannels(serviceID string, query *schema.ServiceAlertChannelQuery) ([]schema.ServiceAlertChannel, error) {
	where, wherePlaceholders := ServiceAlertChannelQuery(query)
	rows, err := i.Query(
		`
      SELECT useralertchannels.id, useralertchannels.userName, useralertchannels.url, useralertchannels.errorAlerts, useralertchannels.uptimeAlerts, useralertchannels.versionAlerts
      FROM useralertchannels
      INNER JOIN users
      ON useralertchannels.userName = users.name
      LEFT JOIN serviceconfigs
      ON serviceconfigs.userName = users.name AND serviceconfigs.serviceID = ?
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
	var records []schema.ServiceAlertChannel
	for rows.Next() {
		var record schema.ServiceAlertChannel
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
