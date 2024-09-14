import ReactDOM from 'react-dom/client';
import React from 'react';
import { AppRoutes } from './routes';
import { AppProvider } from './providers/app';

import './index.css';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

async function main() {
  root.render(
    <AppProvider>
      <AppRoutes />
    </AppProvider>,
  );
}

main();
