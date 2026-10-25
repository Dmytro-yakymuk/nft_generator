/** Describes parameters for sign message */
export class SignedMessage {
    /** Includes message from API, signed message, and wallet address */
    constructor(
        public nonce: string,
        public signature: string,
    ) { }
}
