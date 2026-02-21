import { Route, Routes } from 'react-router';
import LoginPage from './LoginPage';
import NavigateToLoginPage from './NavigateToLoginPage';
import NotFoundPage from './NotFoundPage';
import Admin from './admin/Admin';
import { ApplicationContext, useApplicationState } from './application';
import Dashboard from './dashboard/Dashboard';
import { HTTP_UNAUTHORIZED, HTTPError } from './data';
import { useUser } from './data/user';

function App() {
  const user = useUser();
  const application = useApplicationState(user.data);

  if (
    user.error &&
    (!(user.error instanceof HTTPError) ||
      (user.error as HTTPError).status !== HTTP_UNAUTHORIZED)
  ) {
    return user.error.toString();
  }
  if (user.isLoading && user.isInitial) {
    return null;
  }

  return (
    <ApplicationContext.Provider value={application}>
      <Routes>
        <Route path="/login" element={<LoginPage />} />

        {!application.authorized ? (
          <>
            <Route path="/*" element={<NavigateToLoginPage />} />
          </>
        ) : (
          <>
            {application.admin && (
              <Route path="/admin" element={<Admin />}>
                <Route path="test/*" index element={<div>Test</div>} />

                <Route path="*" element={<NotFoundPage />} />
              </Route>
            )}

            <Route path="/" element={<Dashboard />}>
              <Route path="*" element={<NotFoundPage />} />
            </Route>
          </>
        )}
      </Routes>
    </ApplicationContext.Provider>
  );
}

export default App;
