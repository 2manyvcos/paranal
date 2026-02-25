import { useEffect } from 'react';
import { Route, Routes } from 'react-router';
import NavigateToLoginPage from './NavigateToLoginPage';
import AdminScreen from './admin/AdminScreen';
import { ApplicationContext, useApplicationContextState } from './application';
import AuthScreen from './auth/AuthScreen';
import LoginPage from './auth/LoginPage';
import DashboardScreen from './dashboard/DashboardScreen';
import FavoritesPages from './dashboard/FavoritesPage';
import { useUser } from './data/user';
import { useBooleanDataAttribute, useDataAttribute } from './dataAttributes';
import ErrorScreen from './error/ErrorScreen';
import NotFoundPage from './error/NotFoundPage';
import UnexpectedErrorPage from './error/UnexpectedErrorPage';

const root = document.getElementById('root')!;

function App() {
  const user = useUser();
  const userLoading = user.isLoading && user.isInitial;

  const application = useApplicationContextState(user.data);
  const { appName, tagline, userName, admin, unhide } = application;

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
              <Route path="favorites" element={<FavoritesPages />} />

              <Route path="*" element={<NotFoundPage />} />
            </Route>
          </>
        )}
      </Routes>
    </ApplicationContext.Provider>
  );
}

export default App;
