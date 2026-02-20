import { ConfigProvider } from '@civet/core';
import { ConfigProvider as EventConfigProvider } from '@civet/events';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router';
import App from './App.tsx';
import { dataProvider, eventReceiver } from './data';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <ConfigProvider dataProvider={dataProvider}>
        <EventConfigProvider eventReceiver={eventReceiver}>
          <App />
        </EventConfigProvider>
      </ConfigProvider>
    </BrowserRouter>
  </StrictMode>,
);
