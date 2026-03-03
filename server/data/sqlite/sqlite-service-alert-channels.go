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
			conditions = append(conditions, "IFNULL(user_alert_channels.error_alerts, FALSE) = TRUE AND IFNULL(users.error_alerts, FALSE) = TRUE")
		} else {
			conditions = append(conditions, "(IFNULL(user_alert_channels.error_alerts, FALSE) = FALSE OR IFNULL(users.error_alerts, FALSE) = FALSE)")
		}
	}
	if query.HealthAlerts != nil {
		if *query.HealthAlerts {
			conditions = append(conditions, "IFNULL(user_alert_channels.health_alerts, FALSE) = TRUE AND (IFNULL(users.health_alerts, FALSE) = TRUE OR IFNULL(service_configs.health_alerts, FALSE) = TRUE)")
		} else {
			conditions = append(conditions, "(IFNULL(user_alert_channels.health_alerts, FALSE) = FALSE OR (IFNULL(users.health_alerts, FALSE) = FALSE AND IFNULL(service_configs.health_alerts, FALSE) = FALSE))")
		}
	}
	if query.VersionAlerts != nil {
		if *query.VersionAlerts {
			conditions = append(conditions, "IFNULL(user_alert_channels.version_alerts, FALSE) = TRUE AND (IFNULL(users.version_alerts, FALSE) = TRUE OR IFNULL(service_configs.version_alerts, FALSE) = TRUE)")
		} else {
			conditions = append(conditions, "(IFNULL(user_alert_channels.version_alerts, FALSE) = FALSE OR (IFNULL(users.version_alerts, FALSE) = FALSE AND IFNULL(service_configs.version_alerts, FALSE) = FALSE))")
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
      SELECT user_alert_channels.id, user_alert_channels.user_name, user_alert_channels.url, user_alert_channels.error_alerts, user_alert_channels.health_alerts, user_alert_channels.version_alerts
      FROM user_alert_channels
      INNER JOIN users
      ON user_alert_channels.user_name = users.name
      LEFT JOIN service_configs
      ON service_configs.user_name = users.name AND service_configs.service_id = ?
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
		if err := rows.Scan(&record.ID, &record.UserName, &record.URL, &record.ErrorAlerts, &record.HealthAlerts, &record.VersionAlerts); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}
