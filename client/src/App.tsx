import { useEffect } from 'react';
import { Route, Routes } from 'react-router';
import NavigateToLoginPage from './NavigateToLoginPage';
import AdminScreen from './admin/AdminScreen';
import { ApplicationContext, useApplicationState } from './application';
import AuthScreen from './auth/AuthScreen';
import LoginPage from './auth/LoginPage';
import DashboardScreen from './dashboard/DashboardScreen';
import { HTTP_UNAUTHORIZED, HTTPError } from './data/errors';
import { useUser } from './data/user';
import ErrorScreen from './error/ErrorScreen';
import NotFoundPage from './error/NotFoundPage';
import UnexpectedErrorPage from './error/UnexpectedErrorPage';
import Notifications from './notifications/Notifications';

function App() {
  const user = useUser();
  const application = useApplicationState(user.data);

  const userLoading = user.isLoading && user.isInitial;
  useEffect(() => {
    if (!userLoading) return;
    const root = document.getElementById('root');
    root!.setAttribute('data-loading', '');
    return () => {
      root!.removeAttribute('data-loading');
    };
  }, [userLoading]);

  useEffect(() => {
    if (
      user.error &&
      (!(user.error instanceof HTTPError) ||
        (user.error as HTTPError).status !== HTTP_UNAUTHORIZED)
    )
      console.error(user.error);
  }, [user.error]);

  if (
    user.error &&
    (!(user.error instanceof HTTPError) ||
      (user.error as HTTPError).status !== HTTP_UNAUTHORIZED)
  ) {
    return (
      <Routes>
        <Route path="/" element={<ErrorScreen />}>
          <Route index element={<UnexpectedErrorPage />} />
          <Route path="*" element={<UnexpectedErrorPage />} />
        </Route>
      </Routes>
    );
  }

  if (userLoading) {
    return null;
  }

  return (
    <ApplicationContext.Provider value={application}>
      <Routes>
        <Route path="/login" element={<AuthScreen />}>
          <Route index element={<LoginPage />} />
        </Route>

        {!application.authorized ? (
          <Route path="/*" element={<NavigateToLoginPage />} />
        ) : (
          <>
            {application.admin && (
              <Route path="/admin" element={<AdminScreen />}>
                <Route
                  path="test/*"
                  index
                  element={<div>Test TODO: remove</div>}
                />

                <Route path="*" element={<NotFoundPage />} />
              </Route>
            )}

            <Route path="/" element={<DashboardScreen />}>
              <Route path="*" element={<NotFoundPage />} />
            </Route>
          </>
        )}
      </Routes>

      <Notifications />
    </ApplicationContext.Provider>
  );
}

export default App;
