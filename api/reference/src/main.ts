import { createApiReference } from '@scalar/api-reference';

import '@scalar/api-reference/style.css';

const apiURL = new URL(
  window.paranal.api.replace(/\/*$/, '/'),
  window.location.href,
);

createApiReference('#root', {
  theme: 'bluePlanet',
  sources: [
    {
      title: 'Paranal API v1 (latest)',
      slug: 'openapi.v1',
      url: new URL('schema/openapi.v1.yaml', apiURL).toString(),
    },
  ],
  darkMode: true,
  hideClientButton: true,
  agent: { disabled: true },
  telemetry: false,
  orderSchemaPropertiesBy: 'preserve',
});
