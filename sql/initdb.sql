CREATE TABLE IF NOT EXISTS rules (
  rule_id INT NOT NULL GENERATED ALWAYS AS IDENTITY,
  rule_name varchar(250) NOT NULL,
  rule_type varchar(250) NOT NULL,
  rule_subnetmask varchar(250) NULL, -- null = no filters
  PRIMARY KEY (rule_id)
);




CREATE TABLE IF NOT EXISTS hits (
  hit_from varchar(250) NOT NULL, -- requesting ip address
  hit_to varchar(250) NOT NULL, -- routing address
  hit_datetime timestamp without time zone default (now() at time zone 'utc') -- utc time so there is no confusion
);

