import React from "react";
import {TransactionBody, TransactionType} from "../../Lib/model";

export const NewTransaction: React.FC<{
  onChange: (trx: TransactionBody) => void
  onSubmit: (trx: TransactionBody) => void,
  formValue: TransactionBody
}> = (props) => {
  return (
    <div id="form">
      <div className="form-group mb-3">
        <label className="label-form"  htmlFor="type" aria-describedby="amount">Type</label>
        <select
          value={props.formValue.type}
          onChange={e => props.onChange({...props.formValue, type: e.target.value as TransactionType})}
          id="type"
          className="custom-select mb-3"
        >
          <option value="credit" selected>Credit</option>
          <option value="debit">Debit</option>
        </select>
      </div>
      <div className="form-group mb-3">
        <label className="label-form" htmlFor="amount">Amount</label>
        <input
          value={props.formValue.amount}
          onChange={e => props.onChange({...props.formValue, amount: parseInt(e.target.value, 10)})}
          id="amount"
          type="number"
          className="form-control"
          aria-describedby="amount"/>
      </div>
      <a style={{cursor: "pointer"}} onClick={() => props.onSubmit(props.formValue)}>Submit</a>
    </div>
  )
};