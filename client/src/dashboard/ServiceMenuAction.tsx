import type { FetchProviderType } from '@civet/common';
import { useConfigContext } from '@civet/core';
import { MenuItem } from '@headlessui/react';
import Link from '@/components/Link';
import type { Service } from '@/data/dataServices';
import {
  postServicesByIDActionsByNameRun,
  type Action,
} from '@/data/dataServicesActions';
import { errorMessage, logError } from '@/data/errors';
import { notify } from '@/notifications/notifications';

export default function ServiceMenuAction({
  service,
  action,
}: {
  service: Service;
  action: Action;
}) {
  const { dataProvider } = useConfigContext<FetchProviderType>();

  return (
    <MenuItem>
      <Link
        className="action dropdown-item"
        as="a"
        href={action.url || undefined}
        onClick={() => {
          if (!action.canRun) return;
          notify.promise(
            (async () => {
              try {
                const result = await postServicesByIDActionsByNameRun({
                  dataProvider: dataProvider!,
                  serviceID: service.id,
                  actionName: action.name,
                });
                if (!result.success) throw result.error!;
              } catch (error) {
                logError(error);
                throw errorMessage(error);
              }
            })(),
            {
              loading: `Running action "${action.name}"…`,
              success: `Done.`,
            },
          );
        }}
        text={action.name}
      />
    </MenuItem>
  );
}
