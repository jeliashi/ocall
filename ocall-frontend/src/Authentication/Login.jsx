import React from 'react';

const Login = () => {
    return (
        <div className="login-page">
            <h1>OCall</h1>
            <div className="login-form">
                <form>
                    <label>Username</label>
                    <input name="username" type="text"/>
                    <label>Password</label>
                    <input name="password" type="password"/>
                </form>
            </div>
        </div>
    );
}

export default Login;