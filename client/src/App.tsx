import { ApplicationContext, useApplicationState } from './application';

function App() {
  const application = useApplicationState();

  return (
    <ApplicationContext.Provider value={application}>
      {application.appName}
      <img style={{ height: '1em' }} src={application.logo} />
    </ApplicationContext.Provider>
  );
}

export default App;
