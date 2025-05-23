CREATE TABLE superchat (
      id serial PRIMARY KEY,
      tx_id varchar(64) NOT NULL,
      name varchar(100) NOT NULL,
      message varchar(500) NOT NULL,
      amount numeric(16,8) NOT NULL,
      tkn_category varchar(64) NOT NULL,
      tkn_symbol varchar(32) NOT NULL,
      tkn_logo varchar(200) NOT NULL,
      hidden boolean NOT NULL,
      account_id integer NOT NULL REFERENCES account,
      created timestamp NOT NULL
);

CREATE INDEX superchat_account_idx ON superchat (account_id);
