import React from 'react';
import UrlInput from './UrlInput';
import ProgressComponent from './progress';
import LoraInfoComponent from './loraInfo';

const App: React.FC = () => {
  return (
    <div>
      <h1>Request Sender</h1>
      <UrlInput />
      <LoraInfoComponent />
      <ProgressComponent />
    </div>
  );
};

export default App;
