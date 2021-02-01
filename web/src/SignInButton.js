import React from 'react';

import { Button } from 'react-bootstrap';


class SignInButton extends React.Component {
    render() {
        return (
            <Button variant="outline-dark" size="lg" block
                href={this.props.signInLink}
            >
                <img
                    height={'48px'}
                    style={{ paddingLeft: '8px', marginRight: '24px' }}
                    alt={`Sign in with ${this.props.identityProviderName}`}
                    src={this.props.identityProviderIcon}
                />
                Sign in with {this.props.identityProviderName}
            </Button >
        );
    }
}

export default SignInButton;
