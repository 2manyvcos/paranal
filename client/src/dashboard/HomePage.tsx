import { useMemo } from 'react';
import { useApplication } from '@/application';
import Text from '@/components/Text';
import ServiceCard from './ServiceCard';

export default function HomePage() {
  const { appName, userName, layout, services } = useApplication();

  const assignedServices = useMemo(
    () =>
      Object.fromEntries(
        layout?.pages?.flatMap((page) =>
          !page.name
            ? []
            : (page.sections?.flatMap(
                (section) =>
                  section.serviceIDs?.map((serviceID) => [serviceID, true]) ??
                  [],
              ) ?? []),
        ) ?? [],
      ),
    [layout],
  );

  return (
    <div className="home page">
      <title>{`Home | ${appName}`}</title>

      <header className="header">
        <Text className="name" as="h2" text={`Hi, ${userName}!`} />
      </header>

      <main className="sections">
        <section className="favorites section">
          <header className="header">
            <div className="title">
              <Text className="name" as="h3" text="Favorites" />
            </div>
          </header>

          <div className="services">
            {services.map((service) =>
              !service.config.favorite ? null : (
                <ServiceCard key={service.id} service={service} />
              ),
            )}
          </div>
        </section>

        <section className="unassigned section">
          <header className="header">
            <div className="title">
              <Text className="name" as="h3" text="Unassigned" />
            </div>
          </header>

          <div className="services">
            {services.map((service) =>
              Object.hasOwn(assignedServices, service.id) ? null : (
                <ServiceCard key={service.id} service={service} />
              ),
            )}
          </div>
        </section>
      </main>
    </div>
  );
}
