import { createApiReference } from '@scalar/api-reference';

import '@scalar/api-reference/style.css';

const apiURL = new URL(
  window.paranal.api.replace(/\/*$/g, '/'),
  window.location.href,
);

createApiReference('#root', {
  theme: 'default',
  sources: [
    {
      title: 'Paranal API v1 (latest)',
      slug: 'openapi.v1',
      url: new URL('schema/openapi.v1.yaml', apiURL).toString(),
    },
  ],
  hideClientButton: true,
  agent: { disabled: true },
  telemetry: false,
  orderSchemaPropertiesBy: 'preserve',
});
