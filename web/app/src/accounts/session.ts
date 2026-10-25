const IS_LOGGED_IN = 'IS_LOGGED_IN';
const SESSION_CHANGE_EVENT = 'session-change';

/**
 * Session tells whether the user is logged in. The auth cookie is not accessible from js,
 * so the flag is kept in localStorage: it is set after login and cleared after logout
 * or when the server answers 401 (the token has expired).
 */
export const session = {
    /** Reports whether the user is logged in. */
    isActive(): boolean {
        return !!window.localStorage.getItem(IS_LOGGED_IN);
    },

    /** Marks the user as logged in. */
    start() {
        window.localStorage.setItem(IS_LOGGED_IN, JSON.stringify(true));
        window.dispatchEvent(new Event(SESSION_CHANGE_EVENT));
    },

    /** Marks the user as logged out. */
    end() {
        if (!session.isActive()) {
            return;
        }
        window.localStorage.removeItem(IS_LOGGED_IN);
        window.dispatchEvent(new Event(SESSION_CHANGE_EVENT));
    },

    /** Calls the listener on login and logout in this and other tabs; returns the unsubscribe function. */
    subscribe(listener: () => void): () => void {
        const onStorage = (event: StorageEvent) => event.key === IS_LOGGED_IN && listener();
        window.addEventListener(SESSION_CHANGE_EVENT, listener);
        window.addEventListener('storage', onStorage);

        return () => {
            window.removeEventListener(SESSION_CHANGE_EVENT, listener);
            window.removeEventListener('storage', onStorage);
        };
    },
};
