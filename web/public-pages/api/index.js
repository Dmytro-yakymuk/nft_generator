/**
 * ErrorUnauthorized is a custom error type which indicates that the client request has not been
 * completed because it lacks valid authentication credentials for the requested resource.
 */
 class UnauthorizedError extends Error {
    /** Error message while unautorized */
    constructor(message = 'Unauthorized') {
        super(message);
    }
}

/**
 * BadRequestError is a custom error type which indicates that the server cannot or
 * will not process the request due to something that is perceived to be a client error.
 */
class BadRequestError extends Error {
    /** Error message while bad request */
    constructor(message = 'bad request') {
        super(message);
    };
};
class TooManyRequestError extends Error {
    /** Error message while bad request */
    constructor(message = 'To many requests') {
        super(message);
    };
};

/**
 * NotFoundError is a custom error type which indicates that the server can't find the requested resource.
 */
class NotFoundError extends Error {
    /** Error message while not found request */
    constructor(message = 'not found') {
        super(message);
    };
};

/**
 * InternalError is a custom error type which indicates that the server encountered an unexpected condition
 * that prevented it from fulfilling the request.
 */
class InternalError extends Error {
    /** Error message for internal server error */
    constructor(message = 'internal server error') {
        super(message);
    };
};

const BAD_REQUEST_ERROR = 400;
const TOO_MANY_REQUESTS_ERROR = 429;
const UNAUTORISED_ERROR = 401;
const NOT_FOUND_ERROR = 404;
const INTERNAL_ERROR = 500;

/**
 * APIClient is base client that holds http client and error handler.
 */
class APIClient {
    http = new HttpClient();

    /**
     * handles error due to response code.
     * @param response - response from server.
     *
     * @throws {@link NotFoundError}
     * This exception is thrown if the input is not a valid ISBN number.
     *
     * @throws {@link UnauthorizedError}
     * Thrown if the ISBN number is valid, but no such book exists in the catalog.
     *
     * @throws {@link InternalError}
     * Thrown if the ISBN number is valid, but no such book exists in the catalog.
     *
     * @private
     */
    /* eslint-disable */
    async handleError(response){

        switch (response.status) {
            case BAD_REQUEST_ERROR: throw new BadRequestError();
            case TOO_MANY_REQUESTS_ERROR: throw new TooManyRequestError();
            case NOT_FOUND_ERROR: throw new NotFoundError();
            case UNAUTORISED_ERROR: throw new UnauthorizedError();
            case INTERNAL_ERROR:
            default:
                throw new InternalError();
        }
    }
}
