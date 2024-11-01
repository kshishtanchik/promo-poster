
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import NewActivityCard from './pages/ActivityCard';
import Invite from "./pages/invite/invite";
import RegistrationStatus from './pages/refistrationStatus';
import Registration from './pages/registration';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="invite/:activityId" element={<Invite />} />
        <Route path="event/create/:activityId" element={<NewActivityCard />} />
        <Route path="registration/:activityId/" element={<Registration />} />
        <Route path="registration/status/:activityId/:userId" element={<RegistrationStatus />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
