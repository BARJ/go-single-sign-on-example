import React from 'react';
import { Authenticator } from './Authentication';
import UserPage from './UserPage';
import SignInPage from './SignInPage';

class App extends React.Component {

  constructor(props) {
    super(props);
    this.state = { isSignedIn: false };
  }

  componentDidMount() {
    this.setState({ isSignedIn: Authenticator.isSignedIn() });
  }

  handleSignOut = () => {
    Authenticator.signOut();
    this.setState({ isSignedIn: false })
  }

  render() {
    const isSignedIn = this.state.isSignedIn;

    if (isSignedIn) {
      return (<UserPage onSignOut={this.handleSignOut} />);
    }
    return (<SignInPage />);
  }
};

export default App;
