const DEFAULT_VALUE = 0;

/** Exposes user entity */
export class User {
    constructor(
        public createdAt: string = '',
        public email: string = '',
        public id: string = '',
        public lastLogin: string = '',
        public name: string = '',
        public status: number = DEFAULT_VALUE,
        public isSubscribed: boolean = false,
        public wallet: string = '',
    ) { }
}

/** Contains fields for updating user entity: the names match the JSON of the accounts API. */
export class UserEditFields {
    constructor(
        public name: string = '',
        public email: string = '',
        public isSubscribed: boolean = false
    ) { }
}
