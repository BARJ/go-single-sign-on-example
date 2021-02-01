import React from 'react';
import SignInButton from './SignInButton';
import FacebookIcon from './Facebook.png';
import GithubIcon from './Github.png';
import GoogleIcon from './Google.png';

const facebookSignInLink = 'https://localhost/api/v1/single-sign-on/facebook/sign-in';
const googleSignInLink = 'https://localhost/api/v1/single-sign-on/google/sign-in';
const githubSignInLink = 'https://localhost/api/v1/single-sign-on/github/sign-in';

class SignInPage extends React.Component {
    render() {
        return (
            <div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
                <div style={{ width: '400px' }}>
                    <SignInButton
                        identityProviderName='Facebook'
                        identityProviderIcon={FacebookIcon}
                        signInLink={facebookSignInLink}
                    />
                    <SignInButton
                        identityProviderName='Github'
                        identityProviderIcon={GithubIcon}
                        signInLink={githubSignInLink}
                    />
                    <SignInButton
                        identityProviderName='Google'
                        identityProviderIcon={GoogleIcon}
                        signInLink={googleSignInLink}
                    />
                </div>
            </div >
        );
    }
}

export default SignInPage;
