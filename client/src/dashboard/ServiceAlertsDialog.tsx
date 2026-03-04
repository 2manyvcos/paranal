import { Dialog, DialogBackdrop, DialogPanel } from '@headlessui/react';
import { type Dispatch, type SetStateAction } from 'react';
import { type Service } from '@/data/dataServices';
import ServiceAlertsForm from './ServiceAlertsForm';

export default function ServiceAlertsDialog({
  service,
  open,
  setOpen,
}: {
  service: Service;
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}) {
  return (
    <Dialog
      className="service alerts dialog"
      open={open}
      onClose={() => {
        setOpen(false);
      }}
    >
      <DialogBackdrop className="dialog-backdrop" transition />

      <div className="dialog-frame">
        <DialogPanel className="dialog-panel" transition>
          <ServiceAlertsForm
            service={service}
            onClose={() => {
              setOpen(false);
            }}
          />
        </DialogPanel>
      </div>
    </Dialog>
  );
}
