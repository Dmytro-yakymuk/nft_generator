import { Suspense } from 'react';
import { BrowserRouter } from 'react-router-dom';
import { Provider } from 'react-redux';


import { Notification } from '@components/common/Notification';
import { store } from '@/app/store';
import Loading from './app/views/Loading';
import { MetamaskModal } from './app/components/common/MetamaskModal';
import { useSession } from './app/hooks/useSession';

import { Routes } from './app/routes';

/** Initial App setup. */
function App() {
    const isLoggedIn = useSession();

    return (
        <Provider store={store}>
            <BrowserRouter>
                <Suspense fallback={<Loading />}>
                    <Notification />
                    <MetamaskModal />
                    {/* pages load their data on mount, so they are remounted on login and logout */}
                    <Routes key={String(isLoggedIn)} />
                </Suspense>
            </BrowserRouter >
        </Provider>
    );
}

export default App;
