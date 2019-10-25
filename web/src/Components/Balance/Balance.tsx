import React from 'react';
import {Balance} from "../../Lib/model";

export const BalanceHeader: React.FC<{balance: Balance}> = ({ balance }) => {
    return (
        <header id="header" className="alt">
            <span className="logo"><img src="images/logo.svg" alt=""/></span>
            <h1>$ {new Intl.NumberFormat('en-IN', {
                maximumFractionDigits: 2, minimumFractionDigits: 2
            }).format(balance.total_balance)}</h1>
            <p>This is your current account balance</p>
        </header>
    )
};
