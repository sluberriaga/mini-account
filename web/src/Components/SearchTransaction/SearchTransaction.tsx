import React from "react";
import {TransactionBody, TransactionType} from "../../Lib/model";

export const SearchTransaction: React.FC<{
  onChange: (trx: string) => void
  onSubmit: (trx: string) => void,
  formValue: string
}> = (props) => {
  return (
    <div id="form" style={{marginTop: 15}}>
      <div className="form-group mb-3">
        <label className="label-form" htmlFor="amount">ID</label>
        <input
          value={props.formValue}
          onChange={e => props.onChange(e.target.value)}
          id="amount"
          type="string"
          className="form-control"
          aria-describedby="id"/>
      </div>
      <a style={{cursor: "pointer"}} onClick={() => props.onSubmit(props.formValue)}>Search</a>
    </div>
  )
};