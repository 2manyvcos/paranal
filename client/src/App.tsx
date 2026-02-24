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
import ErrorScreen from './error/ErrorScreen';
import NotFoundPage from './error/NotFoundPage';
import UnexpectedErrorPage from './error/UnexpectedErrorPage';
import { useCustomization } from './useCustomization';

function App() {
  const user = useUser();
  const userLoading = user.isLoading && user.isInitial;

  const application = useApplicationContextState(user.data);

  useCustomization(userLoading, application);

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
