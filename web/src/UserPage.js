import React from 'react';
import { Card, Button, ListGroup, Spinner } from 'react-bootstrap';
import { Authenticator, AuthenticationError } from './Authentication';

class UserPage extends React.Component {

    constructor(props) {
        super(props);
        this.state = { user: null }
    }

    async componentDidMount() {
        await this.loadUser();
    }

    async loadUser() {
        var user = await Authenticator.getUser()
            .catch(err => {
                if (err instanceof AuthenticationError) {
                    this.handleSignOut();
                } else {
                    this.setState({ user: null });
                }
            });
        this.setState({ user: user });

    }

    handleSignOut() {
        if (!this.props.onSignOut) return;
        this.props.onSignOut();
    }

    handleClick = () => {
        this.handleSignOut();
    }

    render() {
        const user = this.state.user;
        var body = (<Spinner animation="border" />);

        if (user) {
            body = (
                <Card>
                    <Card.Img variant="top" src={user['picture']} />
                    <Card.Body>
                        <Card.Title>{user['name']}</Card.Title>
                    </Card.Body>
                    <ListGroup variant="flush">
                        <ListGroup.Item>ID: {user['id']}</ListGroup.Item>
                        <ListGroup.Item>Name: {user['name']}</ListGroup.Item>
                        <ListGroup.Item>Email: {user['email']}</ListGroup.Item>
                    </ListGroup>
                    <Card.Body className="text-center">
                        <Button variant="primary" onClick={this.handleClick}>Sign out</Button>
                    </Card.Body>
                </Card>
            );
        }

        return (
            <div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
                {body}
            </div>
        );
    }
}

export default UserPage;
