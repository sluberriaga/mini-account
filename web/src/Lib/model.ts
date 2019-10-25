export interface Balance {
    total_balance: number;
}

export enum TransactionType {
    Credit = "credit",
    Debit = "debit",
}

export interface TransactionBody {
    type: TransactionType;
    amount: number;
}

export interface Transaction {
    id: string;
    type: TransactionType;
    amount: number;
    date: string;
}