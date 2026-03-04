import type { FetchProviderType } from '@civet/common';
import { useConfigContext } from '@civet/core';
import { Description, DialogTitle, Field, Label } from '@headlessui/react';
import { useCallback, useState, type SubmitEvent } from 'react';
import { useApplication } from '@/application';
import Button from '@/components/Button';
import Switch from '@/components/Switch';
import Text from '@/components/Text';
import { patchServicesByIDConfig, type Service } from '@/data/dataServices';
import { notifyError } from '@/data/errors';

export default function ServiceAlertsForm({
  service,
  onClose,
}: {
  service: Service;
  onClose: () => void;
}) {
  const { dataProvider } = useConfigContext<FetchProviderType>();
  const { user } = useApplication();

  const [healthAlerts, setHealthAlerts] = useState(service.config.healthAlerts);
  const [versionAlerts, setVersionAlerts] = useState(
    service.config.versionAlerts,
  );

  const submit = useCallback(
    async (e: SubmitEvent<HTMLFormElement>) => {
      e.preventDefault();
      onClose();

      try {
        await patchServicesByIDConfig({
          dataProvider: dataProvider!,
          serviceID: service.id,
          data: { healthAlerts, versionAlerts },
        });
      } catch (error) {
        notifyError(error);
      }
    },
    [onClose, dataProvider, service, healthAlerts, versionAlerts],
  );

  return (
    <form className="service alerts form" onSubmit={submit}>
      <Text
        className="title"
        as={DialogTitle as unknown as 'h2'}
        text="Configure alerts"
      />

      <Field className="health-alerts field" disabled={user?.healthAlerts}>
        <Text
          as={Label as unknown as 'label'}
          className="label"
          text="Health alerts"
        />

        <Text
          as={Description as unknown as 'p'}
          className="description"
          text={
            user?.healthAlerts
              ? 'Health alerts are enabled for all services in your profile.'
              : undefined
          }
        />

        <Switch checked={healthAlerts} onChange={setHealthAlerts} />
      </Field>

      <Field className="version-alerts field" disabled={user?.versionAlerts}>
        <Text
          as={Label as unknown as 'label'}
          className="label"
          text="Version alerts"
        />

        <Text
          as={Description as unknown as 'p'}
          className="description"
          text={
            user?.versionAlerts
              ? 'Version alerts are enabled for all services in your profile.'
              : undefined
          }
        />

        <Switch checked={versionAlerts} onChange={setVersionAlerts} />
      </Field>

      <div className="buttons">
        <Button className="submit" type="submit" text="Save" />

        <Button className="close" onClick={onClose} text="Cancel" />
      </div>
    </form>
  );
}
