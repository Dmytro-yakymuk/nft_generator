/**
 * HttpClient is a custom wrapper around fetch api.
 * Exposes get, post and delete methods for JSON strings.
 */
 class HttpClient {
    /**
     * Performs POST http request with JSON body.
     * @param path
     * @param body serialized JSON
     */
    async post(path, body) {
        return await this.do('POST', path, body);
    }

    /**
     * Performs PATCH http request with JSON body.
     * @param path
     * @param body serialized JSON
     */
    async patch(path, body) {
        return await this.do('PATCH', path, body);
    }

    /**
     * Performs PUT http request with JSON body.
     * @param path
     * @param body serialized JSON
     * @param _auth indicates if authentication is needed
     */
    async put(path, body, _auth = true) {
        return await this.do('PUT', path, body);
    }

    /**
     * Performs GET http request.
     * @param path
     * @param _auth indicates if authentication is needed
     */
    async get(path, body, _auth = true) {
        return await this.do('GET', path);
    }

    /**
     * Performs DELETE http request.
     * @param path
     * @param _auth indicates if authentication is needed
     */
    /** TODO: DELETE method will be reworked after back-end remarks.
     * Right now needs body here. */
    async delete(path, body, _auth = true) {
        return await this.do('DELETE', path, body);
    };

    /**
     * do sends an HTTP request and returns an HTTP response as configured on the client.
     * @param method holds http method type
     * @param path
     * @param body serialized JSON
     */
    /* eslint-disable-next-line */
    async do(method, path, body) {
        const request = {
            method: method,
            body: body,
        };

        request.headers = {
            'Content-Type': 'application/json',
        };

        return await fetch(path, request);
    };
};
