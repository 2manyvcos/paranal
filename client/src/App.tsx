import { useEffect } from 'react';
import { Route, Routes } from 'react-router';
import NavigateToLoginPage from './NavigateToLoginPage';
import AdminScreen from './admin/AdminScreen';
import { ApplicationContext, useApplicationContextState } from './application';
import AuthScreen from './auth/AuthScreen';
import LoginPage from './auth/LoginPage';
import DashboardScreen from './dashboard/DashboardScreen';
import HealthStatusesPage from './dashboard/HealthStatusesPage';
import HomePage from './dashboard/HomePage';
import ServicePage from './dashboard/ServicePage';
import VersionsPage from './dashboard/VersionsPage';
import { useUser } from './data/dataUser';
import { useBooleanDataAttribute, useDataAttribute } from './dataAttributes';
import ErrorScreen from './error/ErrorScreen';
import NotFoundPage from './error/NotFoundPage';
import UnexpectedErrorPage from './error/UnexpectedErrorPage';

const root = document.getElementById('root')!;

function App() {
  const user = useUser();
  const userLoading = user.isLoading && user.isInitial;

  const application = useApplicationContextState(user.data);
  const { appName, tagline, userName, admin, unhide, layout } = application;

  useBooleanDataAttribute(root, 'loading', userLoading);
  useDataAttribute(root, 'app-name', appName);
  useDataAttribute(root, 'tagline', tagline);
  useDataAttribute(root, 'user-name', userName);
  useBooleanDataAttribute(root, 'admin', admin);
  useBooleanDataAttribute(root, 'unhide', unhide);

  useEffect(() => {
    if (user.error) console.error(user.error);
  }, [user.error]);

  if (user.error) {
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
                <Route path="*" element={<NotFoundPage />} />
              </Route>
            )}

            <Route path="/" element={<DashboardScreen />}>
              <Route index element={<HomePage />} />
              <Route path="health-statuses" element={<HealthStatusesPage />} />
              <Route path="versions" element={<VersionsPage />} />

              {layout?.pages?.map((page) =>
                !page.name ? null : (
                  <Route
                    key={page.name}
                    path={encodeURIComponent(page.name)}
                    element={<ServicePage page={page} />}
                  />
                ),
              )}

              <Route path="*" element={<NotFoundPage />} />
            </Route>
          </>
        )}
      </Routes>
    </ApplicationContext.Provider>
  );
}

export default App;
