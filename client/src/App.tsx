import { BrowserRouter, Route, Routes } from 'react-router';
import LoginPage from './LoginPage';
import NavigateToLoginPage from './NavigateToLoginPage';
import NotFoundPage from './NotFoundPage';
import Admin from './admin/Admin';
import { ApplicationContext, useApplicationState } from './application';
import Dashboard from './dashboard/Dashboard';

function App() {
  const application = useApplicationState();

  return (
    <BrowserRouter>
      <ApplicationContext.Provider value={application}>
        <Routes>
          <Route path="/login" element={<LoginPage />} />

          {!application.authorized ? (
            <>
              <Route path="/*" element={<NavigateToLoginPage />} />
            </>
          ) : (
            <>
              {application.user?.admin && (
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
    </BrowserRouter>
  );
}

export default App;
