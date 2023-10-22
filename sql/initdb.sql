CREATE TABLE IF NOT EXISTS rules (
  rule_id INT NOT NULL GENERATED ALWAYS AS IDENTITY,
  rule_type INT NOT NULL, -- current implementations 1: round robin, 2: ip hash -- maybe is better string or enum?
  rule_ipaddr varchar(16) NULL, -- null = no filters
  rule_subnetmask int NULL, -- null = no filter
  rule_servers varchar NOT NULL,
  url_source_match varchar NULL DEFAULT NULL, -- catch the request only if match this regex
  url_dest_rewrite varchar NULL DEFAULT NULL, -- rewrite the request using this 
  PRIMARY KEY (rule_id)
);


CREATE TABLE IF NOT EXISTS hits (
  hit_from varchar(250) NOT NULL, -- requesting ip address
  hit_to varchar(250) NOT NULL, -- routing address
  hit_datetime timestamp without time zone default (now() at time zone 'utc') -- utc time so there is no confusion
);


-- test rules
INSERT INTO rules (rule_type, rule_ipaddr, rule_subnetmask, rule_servers, url_source_filter, url_dest_rewrite) VALUES (1, '0.0.0.0', 0, 'google.it:80,microsoft.it:80')