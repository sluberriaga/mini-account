import {Balance, Transaction, TransactionBody} from "./model";

const http = <T>(path: string, request: any): Promise<T> => {
    return fetch(`${process.env.REACT_APP_API_URL}${path}`, {
        ...request,
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        }
    }).then(async res => {
        if (!res.ok) {
            throw Error(await res.text());
        }
        return res.json();
    })
};

export const fetchBalance = (): Promise<Balance> => (
    http("/api/account/balance", { method: "GET" })
);

export const fetchTransactions = (offset: number, limit: number): Promise<Transaction[]> => (
    http(`/api/account/transactions?offset=${offset}&limit=${limit}`, { method: "GET" })
);

export const searchTransaction = (id: string): Promise<Transaction> => (
    http(`/api/account/transaction/${id}`, { method: "GET" })
);

export const postTransaction = (t: TransactionBody): Promise<{}> => (
    http("/api/account/transaction", { method: "POST", body: JSON.stringify(t) })
);