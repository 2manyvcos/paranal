import SwaggerUI from 'swagger-ui';
import 'swagger-ui/dist/swagger-ui.css';

const apiURL = new URL(
  window.paranal.api.replace(/\/*$/g, '/'),
  window.location.href,
);

SwaggerUI({
  url: new URL('schema/openapi.v1.yaml', apiURL).toString(),
  dom_id: '#swagger',
});

/*
Updating the URL can be done with ui.specActions, but there is no predefined UI for that:
https://stackoverflow.com/questions/44816594/swagger-ui-with-multiple-urls/76186038#76186038
*/
// ui.specActions.updateUrl(new URL('schema/openapi.v2.yaml', apiURL).toString());
// ui.specActions.download(new URL('schema/openapi.v2.yaml', apiURL).toString());
