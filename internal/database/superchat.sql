CREATE TABLE superchat (
      id integer NOT NULL PRIMARY KEY AUTO_INCREMENT,
      tx_id varchar(100) NOT NULL,
      name varchar(100) NOT NULL,
      message varchar(500) NOT NULL,
      amount double NOT NULL,
      hidden tinyint(1) NOT NULL,
      user_id integer NOT NULL,
      created datetime NOT NULL
);
