CREATE TABLE account (
    id serial PRIMARY KEY,
    username varchar(15) UNIQUE NOT NULL,
    password varchar(60) NOT NULL,
    address varchar(54) NOT NULL,
    tkn_enabled boolean NOT NULL,
    tkn_address varchar(54) NOT NULL,
    message_max_char integer NOT NULL,
    min_donation numeric(16,8) NOT NULL,
    show_amount boolean NOT NULL,
    widget_id uuid DEFAULT gen_random_uuid() NOT NULL,
    created timestamp NOT NULL
);

CREATE INDEX account_username_idx ON account (username);
