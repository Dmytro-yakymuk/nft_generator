import { lazy } from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';

import { Navbar } from '@components/common/Navbar';
import About from '@components/website/About';
import Content from '@components/website/Content';
import Header from '@components/website/Header';
import Team from '@components/website/Team';
import Footer from '@components/website/Footer';
import Roadmap from '@components/website/Roadmap';

/** Root views. */
const AccountSettings = lazy(() => import('@/app/views/AccountSettings'));
const Dashboard = lazy(() => import('@/app/views/Dashboard'));
const Generator = lazy(() => import('@/app/views/Generator'));
const GenerationMatadata = lazy(() => import('@/app/views/GenerationMatadata'));
const CreateContract = lazy(() => import('@/app/views/CreateContract'));
const CollectionEdit = lazy(() => import('@/app/views/CollectionEdit'));
const Collections = lazy(() => import('@/app/views/Collections'));
const Home = lazy(() => import('@/app/views/Home'));
const WebsiteEdit = lazy(() => import('@/app/views/Website'));
const Whitelist = lazy(() => import('@/app/views/Whitelist'));

/**
 * ComponentRoutes holds all needed information to fill up routes config.
 */
export class ComponentRoutes {
    /** Default constructor.
     * @param path - route path.
     * @param component - component mapped to route.
     * @param exact - indicates that path matches the location pathname exactly.
     */
    constructor(
        public path: string,
        public component: React.FC<any>,
        public exact: boolean,
        public children?: ComponentRoutes[]
    ) { }

    /** Creates child sub routes path. */
    public with(
        child: ComponentRoutes,
        parrent: ComponentRoutes
    ): ComponentRoutes {
        child.path = `${parrent.path}/${child.path}`;

        return this;
    };

    /** Adds children components to base parent component. */
    public addChildren(children: ComponentRoutes[]): ComponentRoutes {
        this.children = children.map((child: ComponentRoutes) =>
            child.with(child, this)
        );

        return this;
    };
}

/**
 * RoutesConfig contains configuration of all available routes and subroutes.
 */
export class RoutesConfig {
    public static AccountSettings: ComponentRoutes = new ComponentRoutes(
        '/account-settings',
        AccountSettings,
        true,
    );
    public static Dashboard: ComponentRoutes = new ComponentRoutes(
        '/dashboard/:id',
        Dashboard,
        true,
    );
    public static CreateContract: ComponentRoutes = new ComponentRoutes(
        '/create-contract/:id',
        CreateContract,
        true,
    );
    public static GenerationMatadata: ComponentRoutes = new ComponentRoutes(
        '/generate-metadata/:id',
        GenerationMatadata,
        true,
    );
    public static WebsiteEdit: ComponentRoutes = new ComponentRoutes(
        '/website/:id',
        WebsiteEdit,
        false,
    );
    public static Header: ComponentRoutes = new ComponentRoutes(
        'header',
        Header,
        true,
    );
    public static About: ComponentRoutes = new ComponentRoutes(
        'about',
        About,
        true,
    );
    public static Content: ComponentRoutes = new ComponentRoutes(
        'content',
        Content,
        true,
    );
    public static Team: ComponentRoutes = new ComponentRoutes(
        'team',
        Team,
        true,
    );
    public static Roadmap: ComponentRoutes = new ComponentRoutes(
        'roadmap',
        Roadmap,
        true,
    );
    public static Footer: ComponentRoutes = new ComponentRoutes(
        'footer',
        Footer,
        true,
    );
    public static Generator: ComponentRoutes = new ComponentRoutes(
        '/collection-generator/:id',
        Generator,
        true,
    );
    public static Default: ComponentRoutes = new ComponentRoutes(
        '/',
        Home,
        true
    );
    public static CollectionEdit: ComponentRoutes = new ComponentRoutes(
        '/edit/:id',
        CollectionEdit,
        false
    );
    public static Collections: ComponentRoutes = new ComponentRoutes(
        '/collections',
        Collections,
        true
    );
    public static Whitelist: ComponentRoutes = new ComponentRoutes(
        '/whitelist/:id',
        Whitelist,
        true
    );
    public static Home: ComponentRoutes = new ComponentRoutes(
        '/home',
        Home,
        true
    );

    /** Root routes path list. */
    public static routes: ComponentRoutes[] = [
        RoutesConfig.AccountSettings,
        RoutesConfig.Dashboard,
        RoutesConfig.CreateContract,
        RoutesConfig.Generator,
        RoutesConfig.GenerationMatadata,
        RoutesConfig.Default,
        RoutesConfig.Collections,
        RoutesConfig.Whitelist,
        RoutesConfig.Home,
        RoutesConfig.WebsiteEdit.addChildren([
            RoutesConfig.Header,
            RoutesConfig.About,
            RoutesConfig.Content,
            RoutesConfig.Team,
            RoutesConfig.Roadmap,
            RoutesConfig.Footer,
        ]),
        RoutesConfig.CollectionEdit,
    ];
}

/** Paths of the former email sign-in pages: the platform logs in with the wallet, from the home page. */
const LEGACY_SIGN_IN_PATHS = ['/sign-in', '/sign-up'];

export const Routes = () =>
    <Switch>
        {LEGACY_SIGN_IN_PATHS.map((path: string) =>
            <Redirect key={path} exact from={path} to={RoutesConfig.Default.path} />
        )}
        {RoutesConfig.routes.map(
            (route: ComponentRoutes, index: number) =>
                <Route
                    key={index}
                    path={route.path}
                    component={route.component}
                    exact={route.exact}
                />
        )}
    </Switch>;
