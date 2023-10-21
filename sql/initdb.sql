CREATE TABLE IF NOT EXISTS rules (
  rule_id INT NOT NULL GENERATED ALWAYS AS IDENTITY,
  rule_name varchar(250) NOT NULL,
  
  PRIMARY KEY (rule_id)
);


CREATE TABLE IF NOT EXISTS hits (
  hit_from varchar(250) NOT NULL, -- ip addr
  hit_to varchar(250) NOT NULL, -- routing address
  hit_datetime timestamp without time zone default (now() at time zone 'utc') -- utc time so there is no confusion
);