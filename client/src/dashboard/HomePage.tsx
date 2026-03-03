import { useApplication } from '@/application';
import Text from '@/components/Text';
import ServiceCard from './ServiceCard';

export default function HomePage() {
  const { servicesWithFavorite } = useApplication();

  return (
    <div className="home page">
      <main className="sections">
        <section className="favorites section">
          <header className="header">
            <div className="title">
              <Text className="name" as="h3" text="Favorites" />
            </div>
          </header>

          <div className="services">
            {servicesWithFavorite.map((service) => (
              <ServiceCard key={service.id} service={service} />
            ))}
          </div>
        </section>
      </main>
    </div>
  );
}
