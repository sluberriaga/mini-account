import React, {useEffect, useState} from 'react';
import './App.css';
import {Transactions} from "./Components/Transactions/Transactions";
import {BalanceHeader} from "./Components/Balance/Balance";
import {NewTransaction} from "./Components/NewTransaction/NewTransaction";

import {fetchBalance, fetchTransactions, postTransaction, searchTransaction} from "./Lib/api";
import {Transaction, TransactionBody, TransactionType} from "./Lib/model";
import {SearchTransaction} from "./Components/SearchTransaction/SearchTransaction";

function compose<T, T1, T2>(f: (a: T) => T1, g: (b: T1) => T2): (a: T) => T2 {
  return (x) => g(f(x))
}

const Main: React.FC<{
  transactions: Transaction[],
  transactionForm: TransactionBody,
  setTransactionForm: (trx: TransactionBody) => void,
  postNewTransaction: (trx: TransactionBody) => void,
  transactionSearchForm: string,
  setTransactionSearchForm: (trx: string) => void,
  searchTransaction: (trx: string) => void
}> = ({
  transactions,
  transactionForm,
  setTransactionForm,
  postNewTransaction,
  transactionSearchForm,
  setTransactionSearchForm,
  searchTransaction
}) => {
    return (
        <div id="main">
            <section id="intro" className="main">
                <div className="spotlight">
                    <div className="new-transaction">
                        <header className="major">
                            <h3>New Transaction</h3>
                        </header>
                        <NewTransaction formValue={transactionForm}
                                        onChange={setTransactionForm}
                                        onSubmit={postNewTransaction}/>
                        <header className="major" style={{marginTop: 20}}>
                          <h3>Search Transaction</h3>
                        </header>
                        <SearchTransaction formValue={transactionSearchForm}
                                      onChange={setTransactionSearchForm}
                                      onSubmit={searchTransaction}/>
                    </div>
                    <div className="transactions">
                        <header className="major">
                            <h3>Transactions</h3>
                        </header>
                        <Transactions transactions={transactions} />
                    </div>
                </div>
            </section>
        </div>
    )
}

const Pager: React.FC<{setPage: (page: number) => void, page: number}> = ({setPage, page}) => {
    return (
        <nav id="nav">
            <ul>
                <li><a onClick={() => setPage(page - 1)}>Prev</a></li>
                <li><a className="active">{page}</a></li>
                <li><a onClick={() => setPage(page + 1)}>Next</a></li>
            </ul>
        </nav>
    )
}

const App: React.FC = () => {
   const [page, setPage] = useState(0);

    const [balance, setBalance] = useState({total_balance: 0});
    const [transactions, setTransactions] = useState([] as Transaction[]);
    const [transactionForm, setTransactionForm] = useState({type: TransactionType.Credit, amount: 0} as TransactionBody);
    const [transactionSearchForm, setTransactionSearchForm] = useState("");

    async function getBalance() {
        const balance = await fetchBalance();
        setBalance(balance);

        return balance
    }
    async function getTransactions() {
        const TRX_PER_PAGE = 5;
        setTransactions(await fetchTransactions(page * TRX_PER_PAGE, TRX_PER_PAGE));
    }
    async function executeSearchTransaction(id: string) {
        errorHandler(searchTransaction(id).then(compose(JSON.stringify, alert)));
    }

    async function executeTransaction() {
        await errorHandler(postTransaction(transactionForm));
        errorHandler(getBalance());
        errorHandler(getTransactions());
    }

    async function changePage(page: number) {
        setPage(page >= 0 ? page : 0);
    }

    function errorHandler(couldThrow: Promise<any>) {
        couldThrow
          .catch(error => {
            alert(error.message)
          });
    }

    function setTransactionBody(trx: TransactionBody): void {
        setTransactionForm(trx);
    }

    function setTransactionSearch(id: string): void {
      setTransactionSearchForm(id);
    }

    useEffect(() => {
        errorHandler(getTransactions())
    }, [page]);

    useEffect(() => {
        errorHandler(getBalance())
    }, []);

    return (
        <div id="wrapper">
            <BalanceHeader balance={balance} />
            <Main
              postNewTransaction={executeTransaction}
              setTransactionForm={setTransactionBody}
              transactions={transactions}
              transactionForm={transactionForm}
              searchTransaction={executeSearchTransaction}
              setTransactionSearchForm={setTransactionSearch}
              transactionSearchForm={transactionSearchForm}
            />
            <Pager setPage={changePage} page={page} />
        </div>
    );
};

export default App;
