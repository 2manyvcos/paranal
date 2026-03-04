import type { FetchProviderType } from '@civet/common';
import { useConfigContext } from '@civet/core';
import {
  MenuHeading,
  MenuItem,
  MenuItems,
  MenuSection,
  MenuSeparator,
} from '@headlessui/react';
import { useMemo } from 'react';
import { useApplication } from '@/application';
import Link from '@/components/Link';
import Text from '@/components/Text';
import type { Service } from '@/data/dataServices';
import { patchServicesByIDConfig } from '@/data/dataServices';
import {
  useServicesByIDActions,
  type Action,
  type ActionGroup,
} from '@/data/dataServicesActions';
import { notifyError, useErrorNotification } from '@/data/errors';
import ServiceMenuAction from './ServiceMenuAction';

type GroupedActions = {
  group: ActionGroup;
  actions: Action[];
};

export default function ServiceMenu({
  service,
  open,
}: {
  service: Service;
  open: boolean;
}) {
  const { dataProvider } = useConfigContext<FetchProviderType>();
  const { admin } = useApplication();

  const actions = useServicesByIDActions({
    disabled: !open,
    serviceID: service.id,
  });
  useErrorNotification(actions);

  const [ungroupedActions, groupedActions] = useMemo<
    [Action[], GroupedActions[]]
  >(() => {
    const ungroupedActions: Action[] = [];
    const groups: { [name: string]: GroupedActions } = {};
    actions.data?.groups.forEach((group) => {
      if (!Object.hasOwn(groups, group.name))
        groups[group.name] = { group, actions: [] };
    });
    actions.data?.actions.forEach((action) => {
      if (!action.group) {
        ungroupedActions.push(action);
        return;
      }
      let group = groups[action.group];
      if (!group) {
        group = { group: { name: action.group, icon: '' }, actions: [] };
        groups[action.group] = group;
      }
      group.actions.push(action);
    });
    return [
      ungroupedActions,
      Object.values(groups).filter((group) => group.actions.length),
    ];
  }, [actions]);

  return (
    <MenuItems
      className="service dropdown"
      anchor="bottom"
      data-loading={actions.isLoading && actions.isInitial ? '' : undefined}
    >
      {!admin ? null : (
        <MenuItem>
          <Link
            className="edit dropdown-item"
            as="a"
            onClick={() => {
              alert('TODO:');
            }}
            text="Edit"
          />
        </MenuItem>
      )}

      <MenuItem>
        <Link
          className="favorite dropdown-item"
          as="a"
          onClick={async () => {
            try {
              await patchServicesByIDConfig({
                dataProvider: dataProvider!,
                serviceID: service.id,
                data: { favorite: !service.config.favorite },
              });
            } catch (error) {
              notifyError(error);
            }
          }}
          text={
            service.config.favorite
              ? 'Remove from favorites'
              : 'Add to favorites'
          }
          data-enabled={service.config.favorite || undefined}
        />
      </MenuItem>

      <MenuItem>
        <Link
          className="hidden dropdown-item"
          as="a"
          onClick={async () => {
            try {
              await patchServicesByIDConfig({
                dataProvider: dataProvider!,
                serviceID: service.id,
                data: { hidden: !service.config.hidden },
              });
            } catch (error) {
              notifyError(error);
            }
          }}
          text={service.config.hidden ? 'Stop hiding' : 'Hide'}
          data-enabled={service.config.hidden || undefined}
        />
      </MenuItem>

      <MenuItem>
        <Link
          className="alerts dropdown-item"
          as="a"
          onClick={() => {
            alert('TODO:');
          }}
          text="Configure alerts"
        />
      </MenuItem>

      {!admin ? null : (
        <MenuItem>
          <Link
            className="scripts dropdown-item"
            as="a"
            onClick={() => {
              alert('TODO:');
            }}
            text="Scripts"
          />
        </MenuItem>
      )}

      <MenuSeparator className="dropdown-separator" />

      {ungroupedActions.map((action) => (
        <ServiceMenuAction
          key={action.name}
          service={service}
          action={action}
        />
      ))}

      {groupedActions.map((bundle) => (
        <MenuSection
          key={bundle.group.name}
          className="action-group dropdown-section"
        >
          <Text className="heading" as={MenuHeading} text={bundle.group.name} />

          {bundle.actions.map((action) => (
            <ServiceMenuAction
              key={action.name}
              service={service}
              action={action}
            />
          ))}
        </MenuSection>
      ))}
    </MenuItems>
  );
}
