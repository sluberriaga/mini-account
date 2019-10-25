import React from 'react';
import './Transactions.css';
import {Transaction, TransactionType} from "../../Lib/model";

export class Transactions extends React.Component<{transactions: Transaction[]}> {
    render () {
        return (
            <div className={'wrapper'}>
                <ul className={'accordion-list'}>
                    {this.props.transactions.map((transaction, key) =>
                        <li className={`accordion-list ${key}`}>
                            <TransactionDetail transaction={transaction} />
                        </li>)
                    }
                </ul>
            </div>
        )
    }
}

const typeColors = {
    [TransactionType.Credit]: "green",
    [TransactionType.Debit]: "red",
};

class TransactionDetail extends React.Component<{transaction: Transaction}, {opened: boolean}> {
    state = { opened: false };

    render () {
        const { props: { transaction }, state: { opened } } = this;

        return (
            <div className={`accordion-item, ${opened && 'accordion-item--opened'}`}>
                <div onClick={() => { this.setState(state => ({ opened: !state.opened })) }}
                     className={'accordion-item__line'}>
                    <h5 className={'accordion-item__title'}>
                        {'   '}
                        <span style={{color: typeColors[transaction.type]}}>
                          $ {new Intl.NumberFormat('en-IN', {
                            maximumFractionDigits: 2, minimumFractionDigits: 2
                          }).format(transaction.amount)}
                        </span>
                    </h5>
                    <span className={'accordion-item__icon'}/>
                </div>
                <div className={'accordion-item__inner'}>
                    <div className={'accordion-item__content'}>
                        <p className={'accordion-item__paragraph'}>
                            ID: {transaction.id}
                            <br />
                            Date: {transaction.date}
                        </p>
                    </div>
                </div>
            </div>
        )
    }
}