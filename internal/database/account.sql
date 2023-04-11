CREATE TABLE account (
    id serial PRIMARY KEY,
    username varchar(25) NOT NULL,
    password varchar(25) NOT NULL,
    address varchar(54) NOT NULL,
    name_max_char integer NOT NULL,
    message_max_char integer NOT NULL,
    min_donation numeric(16,8) NOT NULL,
    show_amount boolean NOT NULL,
    created timestamp NOT NULL
);
