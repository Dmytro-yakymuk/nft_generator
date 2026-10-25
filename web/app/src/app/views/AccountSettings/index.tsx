import { ChangeEvent, FormEvent, useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';

import { Navbar } from '@components/common/Navbar';
import { NavbarSettings } from '@views/Home';

import pinkLogo from '@static/images/logo.png';

import { RoutesConfig } from '@/app/routes';
import { getUser, logout, updateUser } from '@/app/store/actions/users';
import { RootState } from '@/app/store';
import { UserEditFields } from '@/accounts';
import { notify } from '@/internal/notifications';

import './index.scss';

const AccountSettings: React.FC = () => {
    const dispatch = useDispatch();
    const history = useHistory();
    const navbarSettings: NavbarSettings = new NavbarSettings('border-pink', 'text-black', 'background-white', pinkLogo);
    const user = useSelector((state: RootState) => state.accountsReducer);
    const [isSubscribed, setSubscription] = useState<boolean>(false);
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');

    // the form starts from the saved values: Save sends all three fields, so empty inputs would erase them.
    useEffect(() => {
        dispatch(getUser());
    }, []);
    useEffect(() => {
        setName(user.name);
        setEmail(user.email);
        setSubscription(user.isSubscribed);
    }, [user]);

    /** Updates user fields */
    const updateUserData = async() => {
        try {
            await dispatch(updateUser(new UserEditFields(name, email, isSubscribed)));
            notify('Account settings has been updated!', 'success');
        } catch (error) {
            notify('Account settings could not updated! Try later!');
        }
    };

    /** Logs out and returns to the home page */
    const logoutUser = async() => {
        try {
            await dispatch(logout());
            history.push(RoutesConfig.Default.path);
        } catch (error) {
            notify('Could not log out! Try later!');
        }
    };

    return (
        <>
            <Navbar navbarSettings={navbarSettings} />
            <form className="account-settings" onSubmit={(e: FormEvent<HTMLFormElement>) => e.preventDefault()}>
                <div className="account-settings__wrapper">
                    <h1 className="account-settings__title">Account Settings</h1>
                    <label htmlFor="name">Name</label>
                    <input
                        type="text"
                        id="name"
                        min={2}
                        max={50}
                        value={name}
                        onChange={(e: ChangeEvent<HTMLInputElement>) => setName(e.target.value)}
                    />
                    <label htmlFor="email">Email</label>
                    <input
                        type="email"
                        id="email"
                        min={5}
                        max={100}
                        value={email}
                        onChange={(e: ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
                    />
                    <div className="account-settings__subscription">
                        <input
                            id="subscription"
                            type="checkbox"
                            onChange={() => setSubscription(prev => !prev)}
                            checked={isSubscribed}
                        />
                        <label htmlFor="subscription" className="account-settings__subscription__text">News Subscription</label>
                    </div>
                    <button className="account-settings__btn" onClick={updateUserData}>
                        <span className="account-settings__btn-text">Save</span>
                    </button>
                    <button className="account-settings__btn account-settings__btn--logout" onClick={logoutUser}>
                        <span className="account-settings__btn-text">Log out</span>
                    </button>
                </div>
            </form>
        </>
    );
};

export default AccountSettings;
