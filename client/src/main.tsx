import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router';
import App from '@/App.tsx';
import Data from '@/data/Data.tsx';
import Notifications from '@/notifications/Notifications.tsx';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <Data>
        <App />
      </Data>
    </BrowserRouter>

    <Notifications />
  </StrictMode>,
);
