import { useApplication } from '@/application';
import Image from '@/components/Image';
import Markdown from '@/components/Markdown';
import Text from '@/components/Text';
import type { Page } from '@/data/data-layout';

export default function ServicePage({ page }: { page: Page }) {
  const { appName, servicesByID } = useApplication();

  return (
    <div className="service page">
      <title>{`${page.displayName || page.name} | ${appName}`}</title>

      <header className="header">
        <Text className="name" as="h2" text={page.displayName || page.name} />

        <Markdown className="header" markdown={page.header} />
      </header>

      <main className="sections">
        {page.sections?.map((section) =>
          !section.id ? null : (
            <section key={section.id} className="section">
              <header className="header">
                <Image className="icon" image={section.icon} />

                <div className="title">
                  <Text className="name" as="h3" text={section.name} />
                </div>
              </header>

              <div className="services">
                {section.serviceIDs?.map((serviceID) => {
                  const service = servicesByID[serviceID];
                  return !service ? null : (
                    <article key={serviceID} className="service">
                      <header className="header">
                        <a
                          className="header-link"
                          href={service.url || undefined}
                        >
                          <Image className="logo" image={service.logo} />

                          <div className="title">
                            <Text
                              className="name"
                              as="h4"
                              text={service.name}
                            />

                            <Text
                              className="description"
                              text={service.description}
                            />
                          </div>
                        </a>
                      </header>
                    </article>
                  );
                })}
              </div>
            </section>
          ),
        )}
      </main>
    </div>
  );
}
