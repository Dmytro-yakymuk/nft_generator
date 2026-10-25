import { useEffect, useState } from 'react';

import { session } from '@/accounts/session';

/** Hook useSession returns whether the user is logged in and rerenders the component on login and logout. */
export const useSession = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(session.isActive());

    useEffect(() => session.subscribe(() => setIsLoggedIn(session.isActive())), []);

    return isLoggedIn;
};
