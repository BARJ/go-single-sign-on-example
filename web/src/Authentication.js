import Cookies from 'js-cookie';

class AuthenticationError extends Error {
    constructor(message) {
        super(message);
        this.name = 'AuthenticationError';
    }
}

var Authenticator = {

    isSignedIn: function () {
        return !!Cookies.get('token');
    },

    signOut: function () {
        if (Authenticator.isSignedIn()) {
            Cookies.remove('token');
        }
    },

    getUser: function () {
        if (!Authenticator.isSignedIn()) return;

        return fetch('https://localhost/api/v1/me', {
            headers: new Headers({ 'Authorization': 'Bearer ' + Cookies.get('token') })
        })
            .then(rsp => {
                if (rsp.ok) return rsp.json();
                if (rsp.status === 401) throw new AuthenticationError(rsp.statusText);
                throw Error(rsp.status + ': ' + rsp.statusText);
            })
            .then(rsp => rsp['user'])
            .catch(err => {
                console.log(err);
                throw err;
            });
    }
}

export { Authenticator, AuthenticationError }
